package repository

//Package which response for getting and saving info into db

import (
	"context"
	"log"
	"sync"
	"time"
	"wb-snilez-l0/models"

	"github.com/jackc/pgx/v5"
)

type OrderRepo struct {
	conn  *pgx.Conn
	cache *sync.Map
}

func (repo *OrderRepo) InitRepo(dburl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var err error
	repo.conn, err = pgx.Connect(ctx, dburl)
	if err != nil {
		log.Printf("Unable to connect to database: %s %v\n", dburl, err)
		return err
	}

	return nil
}

func (repo *OrderRepo) Store(ord models.Order) error {
	repo.saveToDB(ord)
	return nil
}

func (repo *OrderRepo) Find(order_uid string) (order models.Order, found bool, err error) {
	order, found, err = repo.getFromDB(order_uid)
	return
}

func (repo *OrderRepo) GetAllRows() pgx.Rows {
	rows, err := repo.conn.Query(context.Background(), `SELECT order_uid, track_number, entry, locale,
	internal_signature, customer_id,delivery_service,shardkey,smid
    date_created, oof_shard FROM "order"`)
	if err != nil {
		log.Println("Problem with SELECT from order", err)
	}
	return rows
}

func (repo *OrderRepo) saveToDB(order models.Order) error {
	tx, err := repo.conn.Begin(context.Background())
	if err != nil {
		log.Printf("Unable to begin transaction: %v\n", err)
		return err
	}
	defer tx.Rollback(context.Background()) // отменит изменение если не будет подтверждения транзакции
	_, err = tx.Exec(context.Background(), insertIntoOrder,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		log.Printf("Unable to insert into order: %v\n", err)
		return err
	}
	delivery := &order.Delivery
	_, err = tx.Exec(context.Background(), insertIntoDelivery,
		delivery.OrderUID,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	)
	if err != nil {
		log.Printf("Unable to insert into delivery: %v\n", err)
		return err
	}
	payment := &order.Payment
	_, err = tx.Exec(context.Background(), insertIntoPayment,
		payment.OrderUID,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	)
	if err != nil {
		log.Printf("Insert payment failed: %v\n", err)
		return err
	}
	for i := 0; i < len(order.Items); i++ {
		item := &order.Items[i]
		_, err = tx.Exec(context.Background(), insertIntoItem,
			item.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			log.Printf("Insert item failed: %v\n", err)
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Unable to commit transaction: %v\n", err)
		return err
	}
	return nil
}

func (repo *OrderRepo) getFromDB(order_uid string) (order models.Order, found bool, err error) {
	ctx := context.Background()
	found = true
	tx, err := repo.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `SELECT * FROM "order" WHERE order_uid = $1`, order_uid)
	err = row.Scan(&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			found = false
			err = nil
			log.Printf("didn't find order with uid %v\n", order_uid)
			return
		}
		log.Printf("Unable to query row in orders:  %v\n", err)
		return
	}
	row = tx.QueryRow(ctx, "SELECT * FROM delivery WHERE order_uid = $1", order_uid)

	err = row.Scan(
		&order.Delivery.OrderUID,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil && err != pgx.ErrNoRows {
		log.Printf("Unable to query row at delivery: %v\n", err)
		return
	}
	row = tx.QueryRow(ctx, "SELECT * FROM payment WHERE order_uid = $1", order_uid)

	err = row.Scan(
		&order.Payment.OrderUID,
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil && err != pgx.ErrNoRows {
		log.Printf("Unable to query row at payment: %v\n", err)
		return
	}
	rows, err := tx.Query(ctx, `SELECT * FROM item WHERE order_uid = $1`, order_uid)
	if err != nil {
		log.Printf("Unable to query row at item: %v\n", err)
		return
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Item])
	if err != nil {
		log.Printf("Unable to collect items: %v\n", err)
		return
	}
	order.Items = items

	return
}
