package main

import (
	"log"

	"mangahub/internal/config"
	"mangahub/internal/db"
	tcpsync "mangahub/internal/tcp"
	"mangahub/internal/udp" // ⭐ THÊM
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

	// ⭐ TẠO UDP CLIENT (PHẦN 5)
	udpClient, err := udp.NewClient("127.0.0.1:7070")
	if err != nil {
		log.Fatal(err)
	}

	addr := ":9090"
	s := tcpsync.NewServer(
		addr,
		database,
		cfg.JWTSecret,
		udpClient, // ⭐ TRUYỀN THAM SỐ THỨ 4
	)

	log.Println("TCP Progress Sync server running on", addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
