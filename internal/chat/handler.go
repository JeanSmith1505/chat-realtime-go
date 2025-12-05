package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Permitir conexiones desde cualquier origen (para demo). En producción ajustar.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handler encapsula el hub
type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{hub: h}
}

// ServeWS maneja la conexión websocket
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}
	client := NewClient(h.hub, conn)
	client.hub.register <- client

	// Lanzar goroutines de lectura y escritura
	go client.writePump()
	go client.readPump()
}
