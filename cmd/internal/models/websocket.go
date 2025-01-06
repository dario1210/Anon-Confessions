package models

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}
