package services

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/wpferg/house-price-aggregator-services/structs"
)

type websocketStart struct {
	ID int `json:"id"`
}

type websocketUpdate struct {
	ID      int                          `json:"id"`
	Payload structs.HouseDataAggregation `json:"payload"`
}

type websocketFinish struct {
	ID      int  `json:"id"`
	Success bool `json:"success"`
}

type WebsocketResponse struct {
	ID     int
	Socket *websocket.Conn
}

func (resp *WebsocketResponse) write(data interface{}) {
	resp.Socket.WriteJSON(data)
}

func (resp *WebsocketResponse) Start() {
	resp.write(websocketStart{ID: resp.ID})
}

func (resp *WebsocketResponse) Update(data structs.HouseDataAggregation) {
	resp.write(websocketUpdate{
		ID:      resp.ID,
		Payload: data,
	})
}

func (resp *WebsocketResponse) Finish(success bool) {
	resp.write(websocketFinish{
		ID:      resp.ID,
		Success: success,
	})
}

func (resp *WebsocketResponse) updateLoop(channel chan structs.HouseDataAggregation) {
	update, channelOpen := <-channel
	results := 0
	for channelOpen {
		resp.Update(update)
		update, channelOpen = <-channel
		results++
	}
	resp.Finish(true)

	log.Println("Request update loop completed. Results found:", results)
}

func (resp *WebsocketResponse) Manage(channel chan structs.HouseDataAggregation) {
	go resp.updateLoop(channel)
}
