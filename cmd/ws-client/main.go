package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/ws/chat",
	}

	log.Println("Connecting to", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// Read messages
	go func() {
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Println("← Server:", string(msg))
		}
	}()

	// Send a test message
	testMsg := `{"type":"chat","username":"testuser","content":"Hello WebSocket!"}`
	err = c.WriteMessage(websocket.TextMessage, []byte(testMsg))
	if err != nil {
		log.Println("write error:", err)
		return
	}

	// Wait Ctrl+C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-done:
	case <-interrupt:
		log.Println("interrupt, closing connection")
		_ = c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
	}
}
