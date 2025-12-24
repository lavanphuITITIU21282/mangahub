package main

import (
	"log"

	"mangahub/internal/udp"
)

func main() {
	server, err := udp.NewServer(":7070")
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	log.Println("UDP Notification server running on :7070")
	server.Run()
}
