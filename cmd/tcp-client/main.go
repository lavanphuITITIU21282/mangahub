package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	mhcli "mangahub/internal/cli"
)

func send(conn net.Conn, msg any) {
	b, _ := json.Marshal(msg)
	b = append(b, '\n')
	conn.Write(b)
}

func main() {
	// 🔗 Kết nối TCP server
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("✅ Connected to TCP server")

	reader := bufio.NewReader(conn)

	// 1️⃣ Gửi HELLO
	send(conn, map[string]any{
		"type": "hello",
	})

	resp, _ := reader.ReadString('\n')
	fmt.Println("← Server:", resp)

	// 2️⃣ LẤY TOKEN: ưu tiên ENV, fallback file ~/.mangahub_token
	token := strings.TrimSpace(os.Getenv("JWT_TOKEN"))
	if token == "" {
		token = strings.TrimSpace(mhcli.LoadToken())
	}
	if token == "" {
		fmt.Println("❌ Không tìm thấy JWT token")
		fmt.Println("👉 Cách 1: export JWT_TOKEN=\"TOKEN_CUA_BAN\"")
		fmt.Println("👉 Cách 2: go run ./cmd/cli login <username> <password> (sẽ lưu ~/.mangahub_token)")
		return
	}

	// 3️⃣ Gửi PROGRESS UPDATE
	send(conn, map[string]any{
		"type":     "progress_update",
		"token":    token,
		"manga_id": "one-piece",
		"chapter":  10,
	})

	resp, _ = reader.ReadString('\n')
	fmt.Println("← Server:", resp)

	fmt.Println("🎉 DONE")
}
