package main

import (
	"log"

	"mangahub/internal/api"
	"mangahub/internal/config"
	"mangahub/internal/db"

	"github.com/gin-gonic/gin"
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

	// seed data
	_, _ = database.Exec(`
	INSERT OR IGNORE INTO manga (id, title, author, genres, status, total_chapters, description)
	VALUES
	('one-piece','One Piece','Eiichiro Oda','Action,Adventure','Ongoing',1100,'Pirate adventure'),
	('naruto','Naruto','Masashi Kishimoto','Action,Ninja','Completed',700,'Ninja journey')
	`)

	r := gin.Default()

	mangaHandler := &api.MangaHandler{DB: database}
	r.GET("/manga", mangaHandler.Search)
	r.GET("/manga/:id", mangaHandler.Detail)

	log.Println("HTTP server running on", cfg.HTTPAddr)
	r.Run(cfg.HTTPAddr)
}
