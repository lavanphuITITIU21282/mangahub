package api

import (
	"database/sql"
	"net/http"

	"mangahub/internal/auth"
	"mangahub/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	DB        *sql.DB
	JWTSecret string
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if !domain.IsValidUsername(req.Username) ||
		!domain.IsValidEmail(req.Email) ||
		!domain.IsValidPassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	hashed, _ := auth.HashPassword(req.Password)
	id := uuid.NewString()

	_, err := h.DB.Exec(
		"INSERT INTO users(id, username, email, password_hash) VALUES(?, ?, ?, ?)",
		id, req.Username, req.Email, hashed,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var id, username, passwordHash string
	err := h.DB.QueryRow(
		"SELECT id, username, password_hash FROM users WHERE username = ?",
		req.Username,
	).Scan(&id, &username, &passwordHash)

	if err != nil || !auth.CheckPassword(passwordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, _ := auth.GenerateToken(h.JWTSecret, id, username)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
