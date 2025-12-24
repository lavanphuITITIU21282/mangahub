package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	DB *sql.DB
}

// PUT /users/progress
// Body: {"manga_id":"one-piece","chapter":10}
// Auth: Bearer JWT
func (h *ProgressHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		MangaID string `json:"manga_id"`
		Chapter int    `json:"chapter"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.MangaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Chapter < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chapter must be >= 0"})
		return
	}

	// Validate manga exists and get total chapters for status
	var totalChapters int
	err := h.DB.QueryRow(
		"SELECT total_chapters FROM manga WHERE id = ? LIMIT 1",
		req.MangaID,
	).Scan(&totalChapters)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "manga not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	status := "reading"
	if totalChapters > 0 && req.Chapter >= totalChapters {
		status = "completed"
	}

	_, err = h.DB.Exec(`
		INSERT INTO user_progress (user_id, manga_id, current_chapter, status, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id, manga_id)
		DO UPDATE SET
			current_chapter = excluded.current_chapter,
			status = excluded.status,
			updated_at = CURRENT_TIMESTAMP
	`, userID, req.MangaID, req.Chapter, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "progress updated",
		"manga_id":    req.MangaID,
		"chapter":     req.Chapter,
		"status":      status,
		"updated_at":   time.Now().Format(time.RFC3339),
		"total_chapters": totalChapters,
	})
}
