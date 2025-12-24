package main

import (
	"log"
	"os"

	"mangahub/internal/config"
	"mangahub/internal/db"
	tcpsync "mangahub/internal/tcp"
	"mangahub/internal/udp"
)

func main() {
	cfg := config.Load()

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(database); err != nil {
		log.Fatal(err)
	}

	// UDP notifications (default: 127.0.0.1:9092)
	udpAddr := os.Getenv("MANGAHUB_UDP_ADDR")
	if udpAddr == "" {
		udpAddr = "127.0.0.1:9092"
	}
	udpClient, err := udp.NewClient(udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpClient.Close()

	addr := ":9090"
	s := tcpsync.NewServer(
		addr,
		database,
		cfg.JWTSecret,
		udpClient, 
	)

	log.Println("TCP Progress Sync server running on", addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
