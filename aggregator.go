package main

import (
	"log"
	"math"
	"strings"

	"sync"

	"github.com/wpferg/house-price-services/structs"
)

func addToMap(key string, data structs.HouseData, mapPtr *structs.HouseDataAggregationMap) {
	mutatedMap := *mapPtr
	value, valueExists := mutatedMap[key]

	if !valueExists {
		value = structs.HouseDataAggregation{
			ID:    key,
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

func aggregateThread(inChannel chan structs.HouseData, outChannel chan structs.HouseDataAggregationMap, waitGroup *sync.WaitGroup) {
	unitAggregate := make(structs.HouseDataAggregationMap)
	outcodeAggregate := make(structs.HouseDataAggregationMap)

	houseData, hasData := <-inChannel

	for hasData {
		houseData, hasData = <-inChannel

		codes := strings.Split(houseData.Postcode, " ")

		addToMap(houseData.Postcode, houseData, &unitAggregate)
		addToMap(codes[0], houseData, &outcodeAggregate)
	}

	outChannel <- unitAggregate
	outChannel <- outcodeAggregate

	(*waitGroup).Done()
}

func parallelisedAggregate(channel chan structs.HouseData) []chan structs.HouseDataAggregationMap {
	log.Println("Starting parallelised aggregation.")
	NUM_THREADS := 1

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(NUM_THREADS)

	resultChannels := make([]chan structs.HouseDataAggregationMap, NUM_THREADS)

	for i := range resultChannels {
		outputChannel := make(chan structs.HouseDataAggregationMap, 2)
		resultChannels[i] = outputChannel
		go aggregateThread(channel, outputChannel, &waitGroup)
	}

	waitGroup.Wait()
	log.Println("Parallelised aggregation complete.")

	return resultChannels
}

func mergeAggregates(a, b structs.HouseDataAggregationMap) structs.HouseDataAggregationMap {
	for key, value := range b {
		aggregateValues, exists := a[key]

		if exists {
			aggregateValues.Count += value.Count
			aggregateValues.Total += value.Total
			aggregateValues.Min = int(math.Min(float64(aggregateValues.Min), float64(value.Min)))
			aggregateValues.Max = int(math.Max(float64(aggregateValues.Max), float64(value.Max)))
			a[key] = aggregateValues
		} else {
			a[key] = value
		}
	}
	return a
}

func aggregateResults(resultChannels []chan structs.HouseDataAggregationMap) (structs.HouseDataAggregationMap, structs.HouseDataAggregationMap) {
	log.Println("Starting marge of parallelised aggregate objects.")

	unitAggregate := <-resultChannels[0]
	outcodeAggregate := <-resultChannels[0]

	for _, channel := range resultChannels[1:] {
		otherUnitAggregate := <-channel
		otherOutcodeAggregate := <-channel

		unitAggregate = mergeAggregates(unitAggregate, otherUnitAggregate)
		outcodeAggregate = mergeAggregates(outcodeAggregate, otherOutcodeAggregate)
	}

	log.Println("Completed merge of parallelised aggregate objects.")
	return unitAggregate, outcodeAggregate
}

func Aggregate(channel chan structs.HouseData) ([]structs.HouseDataAggregation, []structs.HouseDataAggregation) {

	resultChannels := parallelisedAggregate(channel)
	unitAggregate, outcodeAggregate := aggregateResults(resultChannels)

	unitList := make([]structs.HouseDataAggregation, len(unitAggregate))
	outcodeList := make([]structs.HouseDataAggregation, len(outcodeAggregate))

	log.Println("Calculating averages for unit level data.")
	i := 0
	for _, value := range unitAggregate {
		value.Average = float32(value.Total) / float32(value.Count)
		unitList[i] = value

		i++
	}

	log.Println("Calculating averages for outcode level data.")
	i = 0
	for _, value := range outcodeAggregate {
		value.Average = float32(value.Total) / float32(value.Count)
		outcodeList[i] = value
		i++
	}

	return unitList, outcodeList
}
