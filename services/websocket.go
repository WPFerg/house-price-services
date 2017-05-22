package services

import (
	"log"
	"math"
	"sync"

	"strings"

	"github.com/gorilla/websocket"
	"github.com/wpferg/house-price-aggregator-services/structs"
)

type websocketRequest struct {
	Method  string `json:"method"`
	Payload string `json:"payload"`
}

func processSearch(responseChannel chan structs.HouseDataAggregation, search string, list *[]structs.HouseDataAggregation, waitGroup *sync.WaitGroup) {
	searchLower := strings.ToLower(search)
	for _, value := range *list {
		if strings.Contains(strings.ToLower(value.ID), searchLower) {
			responseChannel <- value
		}
	}

	waitGroup.Done()
}

func handleRequest(response WebsocketResponse, search string, list *[]structs.HouseDataAggregation) {
	log.Println("Starting distributed search for", search)
	NUM_THREADS := float64(2)
	listSize := float64(len(*list))
	countPerThread := math.Ceil(listSize / NUM_THREADS)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(int(NUM_THREADS))
	responseChannel := make(chan structs.HouseDataAggregation, 1024)

	response.Manage(responseChannel)

	for i := float64(0); i < NUM_THREADS; i++ {
		startIndex := int(math.Max(0, i*countPerThread))
		endIndex := int(math.Min(listSize, (i+1)*countPerThread))
		slice := (*list)[startIndex:endIndex]
		go processSearch(responseChannel, search, &slice, &waitGroup)
	}

	waitGroup.Wait()
	close(responseChannel)

	log.Println("Finished distributed search for", search)
}

func HandleConnection(conn *websocket.Conn, unitData, outcodeData *[]structs.HouseDataAggregation) {
	iterationCount := 0
	for {
		request := websocketRequest{}
		err := conn.ReadJSON(&request)
		response := WebsocketResponse{
			ID:     iterationCount,
			Socket: conn,
		}

		if err != nil {
			log.Println("Error in request", err.Error())
			response.Finish(false)
		} else {
			log.Println("Request information", request)
			response.Start()
			switch request.Method {
			case "outcode-search":
				handleRequest(response, request.Payload, outcodeData)
				break
			case "unit-search":
				handleRequest(response, request.Payload, unitData)
				break
			default:
				response.Finish(false)
				break
			}
		}

		iterationCount++
	}
}
