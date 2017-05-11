package services

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/wpferg/house-price-aggregator/structs"
)

type websocketRequest struct {
	Method  string `json:"method"`
	Payload string `json:"payload"`
}

type websocketResponse struct {
	Success bool                          `json:"success"`
	Payload *structs.HouseDataAggregation `json:"payload,omitempty"`
}

func handleRequest(payload string, mapToUse *structs.HouseDataAggregationMap) websocketResponse {
	matchingData, exists := (*mapToUse)[payload]

	if !exists {
		return websocketResponse{
			Success: true,
			Payload: nil,
		}
	}

	return websocketResponse{
		Success: true,
		Payload: &matchingData,
	}
}

func HandleConnection(conn *websocket.Conn, unitData, outcodeData *structs.HouseDataAggregationMap) {
	for {
		request := websocketRequest{}
		err := conn.ReadJSON(&request)
		response := websocketResponse{}

		if err != nil {
			log.Println("Error in request", err.Error())
			response.Success = false
			response.Payload = nil
		} else {
			log.Println("Request information", request)
			switch request.Method {
			case "outcode-search":
				response = handleRequest(request.Payload, outcodeData)
				break
			case "unit-search":
				response = handleRequest(request.Payload, unitData)
				break
			default:
				response.Success = false
				response.Payload = nil
				break
			}
		}

		conn.WriteJSON(response)
	}
}
