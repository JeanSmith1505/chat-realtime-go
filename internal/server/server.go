package server

import (
	"log"
	"net/http"

	"github.com/jeans/chat-realtime-go/internal/chat"
)

// Server estructura principal
type Server struct {
	hub     *chat.Hub
	handler *chat.Handler
}

// NewServer crea una instancia de servidor
func NewServer() *Server {
	hub := chat.NewHub()
	handler := chat.NewHandler(hub)

	return &Server{
		hub:     hub,
		handler: handler,
	}
}

// Run inicia el servidor HTTP
func (s *Server) Run(address string) {
	// Iniciar el Hub en una goroutine
	go s.hub.Run()

	// Rutas
	http.HandleFunc("/ws", s.handler.ServeWS)
	http.Handle("/", http.FileServer(http.Dir("./web/static")))

	log.Printf("Servidor ejecutándose en %s", address)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
