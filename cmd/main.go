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
	config, err := cfg.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config:%v ", err)
	}

	myApp, err := app.NewApp(fmt.Sprintf("postgres://%s:%s@db:5432/%s",
		config.DBUser,
		config.DBPassword,
		config.DBName),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer myApp.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", myApp.HomeHandler)

	router.HandleFunc("/api/add", myApp.HomeHandler)

	router.HandleFunc("/order/{oderuid}", myApp.GetById)

	log.Println("Starting server on port 8080")
	http.ListenAndServe(":8081", router)
}
