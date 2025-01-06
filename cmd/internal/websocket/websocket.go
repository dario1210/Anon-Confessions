package websocket

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

// We do not need a read buffer size because we are not reading from the client.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	slog.Info("Upgrading HTTP connection to WebSocket")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade connection to WebSocket", slog.String("error", err.Error()))
		return
	}
	slog.Info("WebSocket connection upgraded successfully")

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.Register <- client

	slog.Info("Client registered with hub")

	// Start the write pump in a new goroutine.
	go client.writePump()
	slog.Debug("Started writePump for the client")
}
