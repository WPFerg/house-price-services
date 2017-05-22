package services

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wpferg/house-price-aggregator-services/structs"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var unitData, outcodeData *[]structs.HouseDataAggregation

func handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request, attempting to upgrade")
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Unable to upgrade connection to websocket", err.Error())
		return
	}

	log.Println("Successfully upgraded. Handing over to websocket handler")
	HandleConnection(connection, unitData, outcodeData)
}

func LaunchServices(unit, outcode *[]structs.HouseDataAggregation) {
	log.Println("Attempting to start HTTP Services and Websocket upgrader")

	unitData, outcodeData = unit, outcode
	http.HandleFunc("/", handleHttpRequest)
	http.ListenAndServe(":8080", nil)
}
