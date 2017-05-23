package main

import (
	"log"

	"github.com/wpferg/house-price-services/services"
	"github.com/wpferg/house-price-services/structs"
)

func main() {
	log.Println("House Price Aggregator")

	channel := make(chan structs.HouseData, 5000)

	go LoadFiles(channel, "pp-2016.csv", "pp-2017.csv")
	unitAggregate, outcodeAggregate := Aggregate(channel)

	log.Println("Aggregation completed successfully.")

	services.LaunchServices(&unitAggregate, &outcodeAggregate)
}
