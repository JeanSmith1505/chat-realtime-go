// File: cmd/server/main.go
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jeans/chat-realtime-go/internal/chat"
)

func main() {
	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()

	// Crear hub y arrancar
	h := chat.NewHub()
	go h.Run()

	// Handlers
	handler := chat.NewHandler(h)

	// Servir archivos estáticos (CSS, JS, imágenes)
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// WebSocket
	http.HandleFunc("/ws", handler.ServeWS)

	// Servir index.html en la raíz
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/static/index.html")
	})

	log.Printf("Servidor escuchando en %s\n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
