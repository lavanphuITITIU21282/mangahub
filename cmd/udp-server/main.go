package main

import (
	"log"
	"os"

	"mangahub/internal/udp"
)

func main() {
	addr := os.Getenv("MANGAHUB_UDP_ADDR")
	if addr == "" {
		addr = ":9092"
	}

	s := udp.NewServer(addr)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
