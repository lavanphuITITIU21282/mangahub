package main

import (
	"log"

	"mangahub/internal/config"
	"mangahub/internal/db"
	tcpsync "mangahub/internal/tcp"
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

	addr := ":9090"
	s := tcpsync.NewServer(addr, database, cfg.JWTSecret)

	log.Println("TCP Progress Sync server running on", addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
