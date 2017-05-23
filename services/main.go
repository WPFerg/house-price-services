package services

import (
	"encoding/json"
	"log"
	"net/http"

	"regexp"

	"github.com/gorilla/websocket"
	"github.com/wpferg/house-price-services/structs"
	"github.com/wpferg/house-price-services/util"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var unitData, outcodeData *[]structs.HouseDataAggregation

func handleWebSocketRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request, attempting to upgrade")
	connection, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Unable to upgrade connection to websocket", err.Error())
		return
	}

	log.Println("Successfully upgraded. Handing over to websocket handler")
	HandleConnection(connection, unitData, outcodeData)
}

func handleUnitRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Received unit request.", r.URL.Path)
	regex := regexp.MustCompile("/unit/(?P<unit>.+)")
	results := regex.FindStringSubmatch(r.URL.Path)
	log.Println(results)

	if len(results) > 1 {
		w.Header().Add("Content-Type", "application/json")

		data, _ := json.Marshal(util.ProcessSearch(results[1], unitData))

		w.Write(data)
	} else {
		w.WriteHeader(200)
		w.Write([]byte("Invalid request"))
	}
}

func handleOutcodeRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Received outcode request.", r.URL.Path)
	regex := regexp.MustCompile("/outcode/(?P<outcode>.+)")
	results := regex.FindStringSubmatch(r.URL.Path)
	log.Println(results)

	if len(results) > 1 {
		w.Header().Add("Content-Type", "application/json")

		data, _ := json.Marshal(util.ProcessSearch(results[1], outcodeData))

		w.Write(data)
	} else {
		w.WriteHeader(200)
		w.Write([]byte("Invalid request"))
	}
}

func LaunchServices(unit, outcode *[]structs.HouseDataAggregation) {
	log.Println("Attempting to start HTTP Services and Websocket upgrader")

	unitData, outcodeData = unit, outcode
	http.HandleFunc("/ws/", handleWebSocketRequest)
	http.HandleFunc("/unit/", handleUnitRequest)
	http.HandleFunc("/outcode/", handleOutcodeRequest)
	http.ListenAndServe(":8080", nil)
}
