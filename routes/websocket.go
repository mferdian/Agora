package routes

import (
	"Agora/websocket"
	"github.com/gin-gonic/gin"
)

func WebSocketRoutes(r *gin.Engine, wsHandler *websocket.WsHandler) {
	api := r.Group("/api")
	api.GET("/ws", wsHandler.ServeWs)
}
