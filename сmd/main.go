package main

import (
	"net/http"
	"wb-snilez-l0/internal/app"

	"github.com/gorilla/mux"
)

func main() {
	var myApp app.App
	router := mux.NewRouter()

	router.HandleFunc("/", myApp.HelloHandler)

	http.ListenAndServe(":8080", router)
}
