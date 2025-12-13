// File: cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jeans/chat-realtime-go/internal/chat"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Para local
	}

	addr := ":" + port

	// Crear hub y arrancarlo
	h := chat.NewHub()
	go h.Run()

	// Handlers
	handler := chat.NewHandler(h)

	// Archivos estáticos (CSS, JS, imágenes)
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// WebSocket
	http.HandleFunc("/ws", handler.ServeWS)

	// Servir index.html en la raíz
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/static/index.html")
	})

	log.Println("Servidor escuchando en", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
