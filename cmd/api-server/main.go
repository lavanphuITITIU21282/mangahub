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

	// seed manga
	_, _ = database.Exec(`
	INSERT OR IGNORE INTO manga (id, title, author, genres, status, total_chapters, description)
	VALUES
	('one-piece','One Piece','Eiichiro Oda','Action,Adventure','Ongoing',1100,'Pirate adventure'),
	('naruto','Naruto','Masashi Kishimoto','Action,Ninja','Completed',700,'Ninja journey')
	`)

	r := gin.Default()

	// AUTH
	authHandler := &api.AuthHandler{
		DB:        database,
		JWTSecret: cfg.JWTSecret,
	}
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	// MANGA
	mangaHandler := &api.MangaHandler{DB: database}
	r.GET("/manga", mangaHandler.Search)
	r.GET("/manga/:id", mangaHandler.Detail)

	// USER APIs
	authMW := api.AuthMiddleware(cfg.JWTSecret)
	user := r.Group("/users")
	user.Use(authMW)

	libraryHandler := &api.LibraryHandler{DB: database}
	progressHandler := &api.ProgressHandler{DB: database}

	user.POST("/library", libraryHandler.Add)
	user.GET("/library", libraryHandler.List)
	user.PUT("/progress", progressHandler.Update)

	log.Println("HTTP server running on", cfg.HTTPAddr)
	r.Run(cfg.HTTPAddr)
}
