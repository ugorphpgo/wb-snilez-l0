package main

import (
	"log"
	"net/http"
	"wb-snilez-l0/internal/app"
	"wb-snilez-l0/internal/cfg"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	myApp, err := app.NewApp(cfg.DSN(), cfg.KafkaBroker)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}
	defer myApp.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", myApp.HomeHandler)

	router.HandleFunc("/order/{order_uid}", myApp.GetById).Methods("GET")

	log.Println("Starting server on port 8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
