package tcp

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	addr      string
	db        *sql.DB
	jwtSecret string

	mu      sync.Mutex
	clients map[*Client]struct{}
}

type Client struct {
	conn     net.Conn
	send     chan []byte
	userID   string
	username string
}

func NewServer(addr string, db *sql.DB, jwtSecret string) *Server {
	return &Server{
		addr:      addr,
		db:        db,
		jwtSecret: jwtSecret,
		clients:   make(map[*Client]struct{}),
	}
}

func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		c := &Client{
			conn: conn,
			send: make(chan []byte, 32),
		}
		s.addClient(c)
		go s.writeLoop(c)
		go s.readLoop(c)
	}
}

func (s *Server) addClient(c *Client) {
	s.mu.Lock()
	s.clients[c] = struct{}{}
	s.mu.Unlock()
	log.Println("client connected:", c.conn.RemoteAddr())
}

func (s *Server) removeClient(c *Client) {
	s.mu.Lock()
	delete(s.clients, c)
	s.mu.Unlock()

	_ = c.conn.Close()
	close(c.send)
	log.Println("client disconnected:", c.conn.RemoteAddr())
}

func (s *Server) broadcast(msg any) {
	b, _ := json.Marshal(msg)
	b = append(b, '\n')

	s.mu.Lock()
	defer s.mu.Unlock()
	for c := range s.clients {
		select {
		case c.send <- b:
		default:
			// client bị nghẽn -> drop để tránh treo server
		}
	}
}

func (s *Server) writeLoop(c *Client) {
	for data := range c.send {
		_, err := c.conn.Write(data)
		if err != nil {
			s.removeClient(c)
			return
		}
	}
}

func (s *Server) readLoop(c *Client) {
	defer s.removeClient(c)

	sc := bufio.NewScanner(c.conn)
	// tăng buffer nếu payload dài
	sc.Buffer(make([]byte, 0, 4096), 1024*1024)

	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}

		var base BaseMessage
		if err := json.Unmarshal(line, &base); err != nil {
			s.replyError(c, "invalid_json")
			continue
		}

		switch base.Type {
		case "hello":
			s.handleHello(c)
		case "progress_update":
			if err := s.handleProgressUpdate(c, line); err != nil {
				s.replyError(c, err.Error())
			}
		default:
			s.replyError(c, "unknown_type")
		}
	}

	if err := sc.Err(); err != nil {
		log.Println("read error:", err)
	}
}

func (s *Server) replyOK(c *Client, payload any) {
	resp := map[string]any{"type": "ok", "data": payload}
	b, _ := json.Marshal(resp)
	b = append(b, '\n')
	select {
	case c.send <- b:
	default:
	}
}

func (s *Server) replyError(c *Client, code string) {
	resp := map[string]any{"type": "error", "code": code}
	b, _ := json.Marshal(resp)
	b = append(b, '\n')
	select {
	case c.send <- b:
	default:
	}
}

func (s *Server) handleHello(c *Client) {
	s.replyOK(c, map[string]any{
		"server_time": time.Now().Format(time.RFC3339),
		"msg":         "welcome to mangahub tcp progress sync",
	})
}
