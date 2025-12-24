package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	mangahubv1 "mangahub/proto/mangahubv1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// GRPCClient is a minimal gRPC client used by the CLI.
type GRPCClient struct {
	Addr  string
	Token string
}

func NewGRPCClient(addr string) *GRPCClient {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		addr = "localhost:50051"
	}
	return &GRPCClient{
		Addr:  addr,
		Token: LoadToken(),
	}
}

func (c *GRPCClient) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx,
		c.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

func (c *GRPCClient) SearchManga(query string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := c.dial(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	client := mangahubv1.NewMangaHubServiceClient(conn)
	resp, err := client.SearchManga(ctx, &mangahubv1.SearchMangaRequest{Query: query})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(resp.GetMangas()) == 0 {
		fmt.Println("No manga found")
		return
	}

	for _, m := range resp.GetMangas() {
		fmt.Printf("- %s (%s)\n", m.GetTitle(), m.GetId())
	}
}

func (c *GRPCClient) GetManga(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := c.dial(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	client := mangahubv1.NewMangaHubServiceClient(conn)
	resp, err := client.GetManga(ctx, &mangahubv1.GetMangaRequest{Id: id})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	m := resp.GetManga()
	if m == nil {
		fmt.Println("Manga not found")
		return
	}
	fmt.Printf("Title: %s\nAuthor: %s\n", m.GetTitle(), m.GetAuthor())
}

func (c *GRPCClient) UpdateProgress(mangaID, chapterStr string) {
	ch, err := strconv.Atoi(strings.TrimSpace(chapterStr))
	if err != nil || ch <= 0 {
		fmt.Println("Invalid chapter:", chapterStr)
		return
	}

	if strings.TrimSpace(c.Token) == "" {
		fmt.Println("Missing token. Please login first (mangahub login ...)")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := c.dial(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	client := mangahubv1.NewMangaHubServiceClient(conn)

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+strings.TrimSpace(c.Token))

	resp, err := client.UpdateProgress(ctx, &mangahubv1.UpdateProgressRequest{
		MangaId: mangaID,
		Chapter: int32(ch),
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	msg := strings.TrimSpace(resp.GetMessage())
	if msg == "" {
		msg = "Progress updated"
	}
	fmt.Println(msg)
}
