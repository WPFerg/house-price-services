package main

import "log"
import "os"
import "bufio"

import "github.com/wpferg/house-price-aggregator-services/structs"
import "regexp"
import "strconv"
import "sync"

func loadFile(responseChannel chan structs.HouseData, filepath string, waitGroup *sync.WaitGroup) {
	file, err := os.Open(filepath)

	if err != nil {
		log.Println("Error opening file", err.Error())
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		responseChannel <- parseLine(scanner.Text())
	}

	waitGroup.Done()
	log.Println("Successfully loaded file", filepath)
}

func LoadFiles(responseChannel chan structs.HouseData, filepaths ...string) {
	log.Println("Starting load of prices")
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(filepaths))

	for _, filepath := range filepaths {
		go loadFile(responseChannel, filepath, &waitGroup)
	}

	waitGroup.Wait()
	log.Println("File contents loaded successfully.")
	close(responseChannel)
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
