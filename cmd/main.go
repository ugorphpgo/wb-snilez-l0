package main

import (
	"fmt"
	"log"
	"net/http"
	"wb-snilez-l0/internal/app"

	"github.com/gorilla/mux"
)

func main() {
	dbuser := "postgres"
	dbpassword := "1234"
	dbname := "order_db"

	myApp, err := app.NewApp(fmt.Sprintf("postgres://%s:%s@db:5432/%s", dbuser, dbpassword, dbname))
	if err != nil {
		log.Fatal(err)
	}
	defer myApp.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", myApp.HomeHandler)

	router.HandleFunc("/api/add", myApp.HomeHandler)

	router.HandleFunc("/order/{oderuid}", myApp.GetById)

	log.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", router)
}
