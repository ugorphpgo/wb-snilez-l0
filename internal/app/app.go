package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"wb-snilez-l0/internal/repository"
	"wb-snilez-l0/pkg/models"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

type App struct {
	repo         repository.OrderRepo
	kafka_reader *kafka.Reader
	stop_channel chan struct{}
}

func NewApp(dburl string) (*App, error) {
	app := &App{}
	err := app.repo.InitRepo(dburl)
	if err != nil {
		log.Printf("Failed to init repo: %v", err)
		return nil, err
	}
	app.kafka_reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{dburl},
		Topic:   "orders",
		GroupID: "orders-consumer-group",
	})

	go app.runConsumer()
	return app, nil
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("web/templates/index.html")
	if err != nil {
		log.Printf("Error reading index.html: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error reading index.html %v", err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(html)
}
func (a *App) GetNOrders(w http.ResponseWriter, r *http.Request) {
	count := mux.Vars(r)["count"]
	amount, err := strconv.Atoi(count)
	if err != nil {
		log.Printf("Problem with string to int converting: %v", err)
		fmt.Fprintf(w, "Bad Request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orders, err := a.repo.GetOrders(amount)
	if err != nil {
		log.Printf("Error getting orders: %v", err)
		return
	}
	for i := 0; i < len(orders); i++ {
		json_data, err := json.MarshalIndent(orders[i], "", "\t")
		if err != nil {
			log.Printf("Error making json: %v", err)
		}
		fmt.Fprintf(w, "%s\n", json_data)
		fmt.Fprintf(w, "----------------------------------------\n")
	}
}

func (a *App) GetById(w http.ResponseWriter, r *http.Request) {
	order_uid := mux.Vars(r)["order_uid"]
	order, found, err := a.repo.Find(order_uid)
	if !found {
		fmt.Fprintf(w, "order %v not found\n", order_uid)
		return
	} else if err != nil {
		fmt.Fprintf(w, "internal error\n")
		log.Printf("Error while searching by id: %v", err)
		return
	}
	json_data, err := json.MarshalIndent(order, "", "\t")
	if err != nil {
		log.Printf("Error making json: %v", err)
	}
	fmt.Fprintf(w, "%s\n", json_data)

}

func (a *App) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Opened insert page")

	a.repo.Store(models.MakeRandomOrder())

}

func (a *App) runConsumer() {
	log.Println("Consumer started")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-a.stop_channel
		cancel()
	}()

	for {
		m, err := a.kafka_reader.ReadMessage(ctx)
		if err != nil {
			//останавливаем горутинам - отменяем контекст когда закроется a.stop_channel
			if errors.Is(err, context.Canceled) {
				log.Println("Kafka consumer stopped")
				return
			}
			log.Println("read error:", err)
			continue
		}
		new_order := models.Order{}
		err = json.Unmarshal(m.Value, &new_order)
		if err != nil {
			log.Printf("Error unmarshalling json: %v", err)
			continue
		}
		json_text, _ := json.MarshalIndent(new_order, "", "\t")
		log.Printf("New order:\n%s", json_text)

		err = a.repo.Store(&new_order)
		if err != nil {
			log.Printf("Failed to store order: %v", err)
		}
	}

}

func (a *App) Close() {
	close(a.stop_channel)
	a.repo.Close()
	a.kafka_reader.Close()
}
