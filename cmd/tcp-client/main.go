package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
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

	// 2️⃣ LẤY TOKEN TỪ ENV (đỡ hardcode)
	token := os.Getenv("JWT_TOKEN")
	if token == "" {
		fmt.Println("❌ Bạn chưa set JWT_TOKEN")
		fmt.Println(`👉 Chạy:  export JWT_TOKEN="TOKEN_CUA_BAN"`)
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
