package util

import (
	"strings"
	"sync"

	"github.com/wpferg/house-price-services/structs"
)

func ProcessSearch(search string, list *[]structs.HouseDataAggregation) []structs.HouseDataAggregation {
	searchLower := strings.ToLower(search)
	results := make([]structs.HouseDataAggregation, 1)
	for _, value := range *list {
		if strings.Contains(strings.ToLower(value.ID), searchLower) {
			results = append(results, value)
		}
	}

	return results
}

func ProcessSearchAsync(responseChannel chan structs.HouseDataAggregation, search string, list *[]structs.HouseDataAggregation, waitGroup *sync.WaitGroup) {
	searchLower := strings.ToLower(search)
	for _, value := range *list {
		if strings.Contains(strings.ToLower(value.ID), searchLower) {
			responseChannel <- value
		}
	}

	if waitGroup != nil {
		waitGroup.Done()
	}
}
