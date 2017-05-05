package main

import (
	"log"

	"net/http"

	"github.com/wpferg/house-prices/httpHandlers"
	"github.com/wpferg/house-prices/store"
)

func main() {
	log.Println("House Price Services")

	priceList := LoadFile()

	log.Println("Loaded", len(priceList), "entries")

	store.Set(priceList)

	log.Println("Starting server...")
	startServer()
}

func startServer() {
	http.HandleFunc("/postcode/", httpHandlers.PostcodeSearch)
	err := http.ListenAndServe(":8081", nil)

	if err != nil {
		log.Panicln(err)
	}
}
