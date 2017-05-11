package main

import (
	"log"

	"github.com/wpferg/house-price-aggregator/structs"
)

func main() {
	log.Println("House Price Aggregator")

	channel := make(chan structs.HouseData, 5000)

	go LoadFile(channel)
	Aggregate(channel)

	log.Println("Aggregation completed successfully.")
}
