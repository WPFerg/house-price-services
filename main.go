package main

import (
	"log"

	"github.com/wpferg/house-price-aggregator-services/services"
	"github.com/wpferg/house-price-aggregator/structs"
)

func main() {
	log.Println("House Price Aggregator")

	channel := make(chan structs.HouseData, 5000)

	go LoadFile(channel)
	unitAggregate, outcodeAggregate := Aggregate(channel)

	log.Println("Aggregation completed successfully.")

	services.LaunchServices(&unitAggregate, &outcodeAggregate)
}
