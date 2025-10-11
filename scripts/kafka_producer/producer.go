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
