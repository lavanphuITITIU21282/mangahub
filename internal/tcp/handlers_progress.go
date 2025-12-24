package tcp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) handleProgressUpdate(c *Client, raw []byte) error {
	var msg ProgressUpdateMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return errors.New("invalid_json")
	}
	if msg.Token == "" || msg.MangaID == "" || msg.Chapter <= 0 {
		return errors.New("invalid_payload")
	}

	uid, uname, err := s.verifyToken(msg.Token)
	if err != nil {
		return errors.New("unauthorized")
	}
	c.userID, c.username = uid, uname

	if !s.mangaExists(msg.MangaID) {
		return errors.New("manga_not_found")
	}

	status := "reading"
	now := time.Now().Format(time.RFC3339)

	_, err = s.db.Exec(`
	INSERT INTO user_progress(user_id, manga_id, current_chapter, status, updated_at)
	VALUES(?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(user_id, manga_id)
	DO UPDATE SET current_chapter=excluded.current_chapter, status=excluded.status, updated_at=CURRENT_TIMESTAMP
	`, uid, msg.MangaID, msg.Chapter, status)
	if err != nil {
		return errors.New("db_error")
	}

	s.replyOK(c, map[string]any{
		"user_id":    uid,
		"username":   uname,
		"manga_id":   msg.MangaID,
		"chapter":    msg.Chapter,
		"updated_at": now,
	})

	s.broadcast(ProgressBroadcast{
		Type:      "progress_broadcast",
		UserID:    uid,
		Username:  uname,
		MangaID:   msg.MangaID,
		Chapter:   msg.Chapter,
		Status:    status,
		UpdatedAt: now,
	})

	// ⭐ UDP Notification (broadcast JSON to UDP server)
	if s.udpClient != nil {
		payload, _ := json.Marshal(map[string]any{
			"type":     "progress_update",
			"username": uname,
			"manga_id": msg.MangaID,
			"chapter":  msg.Chapter,
		})
		_ = s.udpClient.Broadcast(string(payload))
	}

	return nil
}

func (s *Server) mangaExists(id string) bool {
	var x string
	err := s.db.QueryRow("SELECT id FROM manga WHERE id = ? LIMIT 1", id).Scan(&x)
	return err == nil
}

func (s *Server) verifyToken(tokenStr string) (string, string, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !t.Valid {
		return "", "", errors.New("invalid")
	}

	claims := t.Claims.(jwt.MapClaims)
	uid, ok1 := claims["user_id"].(string)
	un, ok2 := claims["username"].(string)
	if !ok1 || !ok2 {
		return "", "", errors.New("invalid")
	}
	return uid, un, nil
}

var _ = sql.ErrNoRows
