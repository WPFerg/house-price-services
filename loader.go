package main

import "log"
import "os"
import "bufio"

import "github.com/wpferg/house-prices/structs"
import "regexp"
import "strconv"

func LoadFile() []structs.HouseData {
	log.Println("Attempting to load the house price file.")

	file, err := os.Open("pp-2017.csv")

	if err != nil {
		log.Println("Error opening file", err.Error())
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var results structs.HouseDataList

	for scanner.Scan() {
		results = append(results, parseLine(scanner.Text()))
	}

	return results
}

func parseLine(data string) structs.HouseData {
	regex, _ := regexp.Compile("\",\"")
	csvData := regex.Split(data[1:len(data)-1], -1)
	cost, err := strconv.Atoi(csvData[1])

	if err != nil {
		cost = -1
	}

	return structs.HouseData{
		Id:                csvData[0],
		Cost:              cost,
		Date:              csvData[2],
		Postcode:          csvData[3],
		FlagA:             csvData[4],
		FlagB:             csvData[5],
		FlagC:             csvData[6],
		HouseNameOrNumber: csvData[7],
		AdditionalNumber:  csvData[8],
		Address:           []string{csvData[9], csvData[10], csvData[11], csvData[12]},
		County:            csvData[13],
	}
}
