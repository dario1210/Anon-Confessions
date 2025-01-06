package models

// WebSocketMessage represents a message to be sent via WebSocket.
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}
