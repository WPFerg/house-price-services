package services

import (
	"log"
	"math"
	"sync"

	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/wpferg/house-price-services/structs"
	"github.com/wpferg/house-price-services/util"
)

type websocketRequest struct {
	Method  string `json:"method"`
	Payload string `json:"payload"`
}

type websocketError struct {
	isWebsocketError bool
	error
}

type requestError struct {
	error
	isRequestError bool
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
		go util.ProcessSearchAsync(responseChannel, search, &slice, &waitGroup)
	}

	waitGroup.Wait()
	close(responseChannel)

	log.Println("Finished distributed search for", search)
}

func handleIncomingFrame(conn *websocket.Conn, unitData, outcodeData *[]structs.HouseDataAggregation) (*websocketRequest, error) {
	_, bytes, err := conn.ReadMessage()

	if err != nil {
		return nil, websocketError{
			isWebsocketError: true,
			error:            err,
		}
	}

	request := websocketRequest{}
	err = json.Unmarshal(bytes, &request)

	if err != nil {
		return nil, requestError{
			isRequestError: true,
			error:          err,
		}
	}

	return &request, nil
}

func HandleConnection(conn *websocket.Conn, unitData, outcodeData *[]structs.HouseDataAggregation) {
	iterationCount := 0
	defer conn.Close()
	for {
		request, err := handleIncomingFrame(conn, unitData, outcodeData)

		if err != nil {
		}

		response := WebsocketResponse{
			ID:     iterationCount,
			Socket: conn,
		}
		if err != nil {
			if requestErrorDetails, ok := err.(requestError); ok {
				log.Println("Error in request\n", requestErrorDetails.Error())
				response.Finish(false)
			} else if websocketError, ok := err.(websocketError); ok {
				log.Println("Error reading data from socket (Potentially closed). Halting. Error details:\n", websocketError.Error())
				break
			} else {
				log.Println("Unhandled error:\n", err.Error())
				response.Finish(false)
			}
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
