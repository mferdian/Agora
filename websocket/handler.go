package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader mengubah koneksi HTTP menjadi WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// ⚠️ Production: ganti dengan whitelist origin
		return true
	},
}

// WsHandler menangani endpoint WebSocket
type WsHandler struct {
	Hub *Hub
}

// ServeWs melayani koneksi WebSocket baru
func (h *WsHandler) ServeWs(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "cannot upgrade to websocket"})
		return
	}

	client := &Client{
		hub:  h.Hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	h.Hub.register <- client
	log.Println("🔌 WebSocket client connected:", conn.RemoteAddr())

	go client.WritePump()
	go client.ReadPump()
}
