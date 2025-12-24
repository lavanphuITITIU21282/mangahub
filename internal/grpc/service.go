package grpc

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mangahubv1 "mangahub/proto/mangahubv1"
)

type MangaHubService struct {
	mangahubv1.UnimplementedMangaHubServiceServer
	DB        *sql.DB
	JWTSecret string
}

func (s *MangaHubService) SearchManga(ctx context.Context, req *mangahubv1.SearchMangaRequest) (*mangahubv1.SearchMangaResponse, error) {
	q := strings.TrimSpace(req.GetQuery())

	var rows *sql.Rows
	var err error

	if q == "" {
		rows, err = s.DB.Query(`SELECT id, title, author, genres, status, total_chapters, description FROM manga`)
	} else {
		like := "%" + q + "%"
		rows, err = s.DB.Query(`
			SELECT id, title, author, genres, status, total_chapters, description
			FROM manga
			WHERE title LIKE ? OR author LIKE ? OR genres LIKE ?
		`, like, like, like)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "db error")
	}
	defer rows.Close()

	resp := &mangahubv1.SearchMangaResponse{}
	for rows.Next() {
		var m mangahubv1.Manga
		var chapters int
		if err := rows.Scan(&m.Id, &m.Title, &m.Author, &m.Genres, &m.Status, &chapters, &m.Description); err != nil {
			return nil, status.Error(codes.Internal, "db error")
		}
		m.TotalChapters = int32(chapters)
		resp.Mangas = append(resp.Mangas, &m)
	}
	return resp, nil
}

func (s *MangaHubService) GetManga(ctx context.Context, req *mangahubv1.GetMangaRequest) (*mangahubv1.GetMangaResponse, error) {
	id := strings.TrimSpace(req.GetId())
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	var m mangahubv1.Manga
	var chapters int
	err := s.DB.QueryRow(
		`SELECT id, title, author, genres, status, total_chapters, description FROM manga WHERE id = ?`,
		id,
	).Scan(&m.Id, &m.Title, &m.Author, &m.Genres, &m.Status, &chapters, &m.Description)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "manga not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "db error")
	}

	m.TotalChapters = int32(chapters)
	return &mangahubv1.GetMangaResponse{Manga: &m}, nil
}

func (s *MangaHubService) UpdateProgress(ctx context.Context, req *mangahubv1.UpdateProgressRequest) (*mangahubv1.UpdateProgressResponse, error) {
	userID, err := s.userIDFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	mangaID := strings.TrimSpace(req.GetMangaId())
	chapter := req.GetChapter()
	if mangaID == "" || chapter <= 0 {
		return nil, status.Error(codes.InvalidArgument, "manga_id and chapter are required")
	}

	// Upsert progress (SQLite ON CONFLICT).
	_, dbErr := s.DB.Exec(`
		INSERT INTO user_progress (user_id, manga_id, current_chapter, status, updated_at)
		VALUES (?, ?, ?, 'reading', CURRENT_TIMESTAMP)
		ON CONFLICT(user_id, manga_id)
		DO UPDATE SET current_chapter = excluded.current_chapter, updated_at = CURRENT_TIMESTAMP
	`, userID, mangaID, int(chapter))

	if dbErr != nil {
		return nil, status.Error(codes.Internal, "db error")
	}

	return &mangahubv1.UpdateProgressResponse{Message: "progress updated"}, nil
}

func (s *MangaHubService) userIDFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	vals := md.Get("authorization")
	if len(vals) == 0 {
		vals = md.Get("Authorization")
	}
	if len(vals) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization")
	}

	h := vals[0]
	if !strings.HasPrefix(h, "Bearer ") {
		return "", status.Error(codes.Unauthenticated, "invalid authorization")
	}

	tokenStr := strings.TrimPrefix(h, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}

	uid, _ := claims["user_id"].(string)
	if strings.TrimSpace(uid) == "" {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}
	return uid, nil
}
