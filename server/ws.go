package server

// Implementaion of the web server struct
type WebSocket struct {
	test string
}

func NewWebSocket(test string) *WebSocket {
	return &WebSocket{
		test: test,
	}
}
