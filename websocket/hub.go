package websocket

import (
	"Agora/dto"
	"Agora/service"
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

// Hub menyimpan semua client aktif dan mengatur broadcast
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	commentSvc service.ICommentService
}

// NewHub membuat instance baru dari Hub
func NewHub(commentSvc service.ICommentService) *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		commentSvc: commentSvc,
	}
}

// Run menangani event register, unregister, dan broadcast
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("✅ Client connected:", client.conn.RemoteAddr())

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("❌ Client disconnected:", client.conn.RemoteAddr())
			}

		case message := <-h.broadcast:
			h.handleMessage(message)
		}
	}
}

// handleMessage memproses pesan dari client
func (h *Hub) handleMessage(message *Message) {
	log.Printf("📨 Received message: %+v\n", message)

	if message.Type == "comment" {
		h.handleCommentMessage(message)
	}

	data, _ := json.Marshal(map[string]interface{}{
		"type":        message.Type,
		"user_id":     message.UserID,
		"proposal_id": message.ProposalID,
		"content":     message.Content,
	})

	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

// handleCommentMessage memproses pesan tipe "comment" dan simpan ke DB
func (h *Hub) handleCommentMessage(msg *Message) {
	ctx := context.Background()

	req := dto.CreateCommentRequest{
		ProposalID: uuid.MustParse(msg.ProposalID),
		Content:    msg.Content,
	}

	_, err := h.commentSvc.CreateComment(ctx, req)
	if err != nil {
		log.Println("❗ DB Error:", err)
		h.sendErrorToAll("failed to save comment to database")
	}
}

// sendErrorToAll mengirim error ke semua client
func (h *Hub) sendErrorToAll(msg string) {
	data, _ := json.Marshal(map[string]string{
		"type":    "error",
		"message": msg,
	})
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
