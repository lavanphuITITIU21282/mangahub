package main

import (
	"log"

	"mangahub/internal/config"
	"mangahub/internal/db"
	"mangahub/internal/grpc"
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

	s := grpc.Server{
		Addr:      cfg.GRPCAddr,
		DB:        database,
		JWTSecret: cfg.JWTSecret,
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
