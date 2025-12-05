// File: internal/chat/client.go
package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Tiempo máximo para escribir un mensaje desde el servidor al cliente
	writeWait = 10 * time.Second
	// Tiempo máximo para leer el siguiente pong del cliente
	pongWait = 60 * time.Second
	// Frecuencia de ping al cliente
	pingPeriod = (pongWait * 9) / 10
	// Tamaño máximo de mensaje recibido
	maxMessageSize = 512
)

// Client representa una conexión WebSocket
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	// Buffered channel de envío
	send chan []byte
	// Username (opcional)
	Name string
}

func NewClient(h *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

// readPump lee mensajes del websocket y los reenvía al hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.hub.broadcast <- message
	}
}

// writePump envía mensajes desde el canal send al websocket
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub cerró el canal
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// Primero escribe el mensaje recibido
			w.Write(message)

			// Enviar mensajes restantes en buffer
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write([]byte(<-c.send))
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
