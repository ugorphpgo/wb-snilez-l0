package main

import (
	"fmt"
	"log"
	"net/http"
	"wb-snilez-l0/internal/app"
	"wb-snilez-l0/internal/cfg"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	myApp, err := app.NewApp(fmt.Sprintf("postgres://%s:%s@db:5432/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName),
	)
	if err != nil {
		log.Fatal("Failed to init")
	}
	defer myApp.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", myApp.HomeHandler)
	r.HandleFunc("/order/random/{count}", myApp.GetNOrders).Methods("GET")
	r.HandleFunc("/order/{order_uid}", myApp.GetById).Methods("GET")

	log.Println("Server is up")
	http.ListenAndServe(":8081", r)
}
