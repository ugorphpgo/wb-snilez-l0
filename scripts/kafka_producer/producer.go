package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ugorphpgo/wb-snilez-l0/pkg/models"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

type TestProducer struct {
	writer *kafka.Writer
}

func MakeTestProducer() *TestProducer {
	return &TestProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP("localhost:9092"),
			Topic:    "test",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (tp *TestProducer) ProduceRandomOrder() error {
	order := models.MakeRandomOrder()

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Printf("Failed to marshall order: %s", err)
		return err
	}
	err = tp.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: orderJSON,
		},
	)
	if err != nil {
		log.Printf("Failed to produce random order: %v", err)
		return err
	}
	return nil
}

func (tp *TestProducer) Close() {
	tp.writer.Close()
}

func main() {
	testKafka()

	fmt.Println("Producer is up")

	producer := MakeTestProducer()
	defer producer.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "producer is ok")
	})

	r.HandleFunc("/producer", func(w http.ResponseWriter, r *http.Request) {
		err := producer.ProduceRandomOrder()
		if err != nil {
			fmt.Fprintf(w, "Failed to produce random order: %v", err)
		} else {
			fmt.Fprintf(w, "Produced random order")
		}
	})
	http.ListenAndServe(":8082", r)
}

func testKafka() {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		log.Printf("Failed to dial kafka: %s", err)
	}
	defer conn.Close()

	topics, _ := conn.ReadPartitions()
	found := false
	for _, p := range topics {
		if p.Topic == "orders" {
			found = true
			break
		}
	}
	if !found {
		err := conn.CreateTopics(kafka.TopicConfig{
			Topic:             "orders",
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
		if err != nil {
			log.Fatalf("Failed to create topic: %s", err)
		}
	}
}
