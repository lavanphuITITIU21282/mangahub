package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LibraryHandler struct {
	DB *sql.DB
}

// POST /users/library
func (h *LibraryHandler) Add(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		MangaID string `json:"manga_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.MangaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	_, err := h.DB.Exec(`
		INSERT OR IGNORE INTO user_progress (user_id, manga_id, current_chapter, status)
		VALUES (?, ?, 0, 'reading')
	`, userID, req.MangaID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "added to library"})
}

// GET /users/library
func (h *LibraryHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	rows, err := h.DB.Query(`
		SELECT m.id, m.title, up.current_chapter, up.status
		FROM user_progress up
		JOIN manga m ON m.id = up.manga_id
		WHERE up.user_id = ?
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()

	var result []gin.H
	for rows.Next() {
		var id, title, status string
		var chapter int

		rows.Scan(&id, &title, &chapter, &status)
		result = append(result, gin.H{
			"id":              id,
			"title":           title,
			"current_chapter": chapter,
			"status":          status,
		})
	}

	c.JSON(http.StatusOK, result)
}
