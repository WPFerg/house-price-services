package store

import (
	"strings"

	"github.com/wpferg/house-prices/structs"
)

var store structs.HouseDataList

func Set(data structs.HouseDataList) {
	store = data
}

func SearchPostcode(fragment string) structs.HouseDataList {
	var result structs.HouseDataList
	lowercasePostcodeFragment := strings.ToLower(fragment)

	for _, houseData := range store {
		lowercasePostcode := strings.ToLower(houseData.Postcode)

		if strings.Contains(lowercasePostcode, lowercasePostcodeFragment) {
			result = append(result, houseData)
		}
	}

	return result
}
