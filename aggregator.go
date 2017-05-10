package main

import (
	"log"
	"strings"

	"github.com/wpferg/house-price-aggregator/structs"
)

func Aggregate(channel chan structs.HouseData) (structs.HouseDataAggregationMap, structs.HouseDataAggregationMap) {
	unitAggregate := make(structs.HouseDataAggregationMap)
	outcodeAggregate := make(structs.HouseDataAggregationMap)

	houseData, hasData := <-channel
	iteration := 0

	for hasData {
		houseData, hasData = <-channel

		codes := strings.Split(houseData.Postcode, " ")

		addToMap(houseData.Postcode, houseData, &unitAggregate)
		addToMap(codes[0], houseData, &outcodeAggregate)

		iteration++

		if iteration%100000 == 0 {
			log.Println("Aggregated", iteration, "entries")
		}
	}

	log.Println("Aggregation complete. Processed", iteration, "entries in total.")

	log.Println("Calculating averages for unit level data.")
	for key, value := range unitAggregate {
		value.Average = float32(value.Total) / float32(value.Count)
		unitAggregate[key] = value
	}

	log.Println("Calculating averages for outcode level data.")
	for key, value := range outcodeAggregate {
		value.Average = float32(value.Total) / float32(value.Count)
		outcodeAggregate[key] = value
	}

	return unitAggregate, outcodeAggregate
}

func addToMap(key string, data structs.HouseData, mapPtr *structs.HouseDataAggregationMap) {
	mutatedMap := *mapPtr
	value, valueExists := mutatedMap[key]

	if !valueExists {
		value = structs.HouseDataAggregation{
			Min:   data.Cost,
			Max:   data.Cost,
			Total: data.Cost,
			Count: 1,
		}
	} else {
		if data.Cost < value.Min {
			value.Min = data.Cost
		}

		if data.Cost > value.Max {
			value.Max = data.Cost
		}

		value.Total += data.Cost
		value.Count++
	}

	mutatedMap[key] = value
}
