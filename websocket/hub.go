package websocket

import "encoding/json"

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

var WsHub = NewHub()

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			delete(h.Clients, client)
			close(client.Send)
		case message := <-h.Broadcast:
			for client := range h.Clients {
				client.Send <- message
			}
		}
	}
}

func (h *Hub) BroadcastEvent(eventType string, data interface{}) {
	msg := map[string]interface{}{
		"type": eventType,
		"data": data,
	}
	jsonMsg, _ := json.Marshal(msg)
	h.Broadcast <- jsonMsg
}
