package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	DB *sql.DB
}

// API 3: PUT /users/progress
func (h *ProgressHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		MangaID string `json:"manga_id"`
		Chapter int    `json:"chapter"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.MangaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	_, err := h.DB.Exec(`
		UPDATE user_progress
		SET current_chapter = ?, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = ? AND manga_id = ?
	`, req.Chapter, userID, req.MangaID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "progress updated"})
}
