package websocket

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			slog.Warn("Failed to write message to websocket", slog.String("error", err.Error()))
			return
		}
	}

	if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		slog.Warn("Failed to send close message to websocket", slog.String("error", err.Error()))
	}
}
