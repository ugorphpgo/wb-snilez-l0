package main

import (
	"log"
	"wb-snilez-l0/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("init app: %v", err)
	}
	if err := a.Run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
