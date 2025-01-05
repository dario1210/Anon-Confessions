package websocket

import "github.com/gin-gonic/gin"

func RegisterWebSocketRoutes(router *gin.RouterGroup, hub *Hub) {
	router.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})
}
