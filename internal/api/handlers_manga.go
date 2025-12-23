package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MangaHandler struct {
	DB *sql.DB
}

func (h *MangaHandler) Search(c *gin.Context) {
	rows, err := h.DB.Query(
		`SELECT id, title, author, genres, status, total_chapters FROM manga`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()

	var result []gin.H
	for rows.Next() {
		var id, title, author, genres, status string
		var chapters int
		rows.Scan(&id, &title, &author, &genres, &status, &chapters)

		result = append(result, gin.H{
			"id":            id,
			"title":         title,
			"author":        author,
			"genres":        genres,
			"status":        status,
			"totalChapters": chapters,
		})
	}

	c.JSON(http.StatusOK, result)
}

func (h *MangaHandler) Detail(c *gin.Context) {
	id := c.Param("id")

	var title, author, genres, status, description string
	var chapters int

	err := h.DB.QueryRow(
		`SELECT title, author, genres, status, total_chapters, description
		 FROM manga WHERE id = ?`,
		id,
	).Scan(&title, &author, &genres, &status, &chapters, &description)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "manga not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            id,
		"title":         title,
		"author":        author,
		"genres":        genres,
		"status":        status,
		"totalChapters": chapters,
		"description":   description,
	})
}
