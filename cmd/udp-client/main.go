package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mangahub/internal/udp"
)

func main() {
	server := flag.String("server", "127.0.0.1:9092", "udp server address")
	mode := flag.String("mode", "subscribe", "subscribe|broadcast|release|list|ping|unsub")
	msg := flag.String("msg", "", "message for broadcast (raw string)")

	// release notification payload flags
	mangaID := flag.String("manga", "", "manga id for chapter release notification")
	title := flag.String("title", "", "manga title for chapter release notification (optional)")
	chapter := flag.Int("chapter", 0, "chapter number for chapter release notification")
	text := flag.String("text", "", "custom message text for chapter release (optional)")
	flag.Parse()

	c, err := udp.NewClient(*server)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	ctx, cancel := signalContext()
	defer cancel()

	switch *mode {
	case "subscribe":
		if err := c.Subscribe(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("✅ subscribed. Listening... Ctrl+C to stop")
		if err := c.Listen(ctx, func(t string) { fmt.Println(t) }); err != nil {
			log.Fatal(err)
		}

	case "release":
		if *mangaID == "" {
			log.Fatal("missing -manga")
		}
		if *chapter <= 0 {
			log.Fatal("missing/invalid -chapter")
		}
		mTitle := *title
		if mTitle == "" {
			mTitle = *mangaID
		}
		message := *text
		if message == "" {
			message = fmt.Sprintf("New chapter released: %s #%d", mTitle, *chapter)
		}

		payload, err := json.Marshal(map[string]any{
			"type":        "chapter_release",
			"manga_id":    *mangaID,
			"manga_title": mTitle,
			"chapter":     *chapter,
			"message":     message,
			"released_at": time.Now().Format(time.RFC3339),
		})
		if err != nil {
			log.Fatal(err)
		}
		if err := c.Broadcast(string(payload)); err != nil {
			log.Fatal(err)
		}
		fmt.Println("✅ chapter release broadcast sent. Waiting reply...")
		_ = c.Listen(ctx, func(t string) { fmt.Println(t) })

	case "broadcast":
		if *msg == "" {
			log.Fatal("missing -msg")
		}
		if err := c.Broadcast(*msg); err != nil {
			log.Fatal(err)
		}
		// đọc reply 1 chút cho vui
		fmt.Println("✅ broadcast sent. Waiting reply...")
		_ = c.Listen(ctx, func(t string) { fmt.Println(t) })

	case "list":
		if err := c.List(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("✅ list requested. Waiting reply...")
		_ = c.Listen(ctx, func(t string) { fmt.Println(t) })

	case "ping":
		if err := c.Ping(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("✅ ping sent. Waiting reply...")
		_ = c.Listen(ctx, func(t string) { fmt.Println(t) })

	case "unsub":
		if err := c.Unsubscribe(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("✅ unsub sent")

	default:
		log.Fatal("unknown mode")
	}
}

func signalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		cancel()
	}()
	return ctx, cancel
}
