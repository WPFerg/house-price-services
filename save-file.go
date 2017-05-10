package main

import (
	"encoding/json"

	"io/ioutil"

	"github.com/wpferg/house-price-aggregator/structs"
)

func SaveMap(filename string, data structs.HouseDataAggregationMap) {
	marshalled, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filename, marshalled, 0644)

	if err != nil {
		panic(err)
	}
}
