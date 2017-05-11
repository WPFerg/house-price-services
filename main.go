package main

import (
	"log"

	"github.com/wpferg/house-price-aggregator/structs"
)

func main() {
	log.Println("House Price Aggregator")

	channel := make(chan structs.HouseData, 5000)

	go LoadFile(channel)
	postcodeData, outcodeData := Aggregate(channel)

	log.Println("Attempting to save postcode-level data")
	SaveMap("postcode-data.json", postcodeData)

	log.Println("Attempting to save outcode-level data")
	SaveMap("outcode-data.json", outcodeData)

	log.Println("Aggregation completed successfully.")
}
