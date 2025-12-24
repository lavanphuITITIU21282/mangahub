package udp

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Addr string
	conn *net.UDPConn

	mu   sync.RWMutex
	subs map[string]*net.UDPAddr
}

func NewServer(addr string) *Server {
	if addr == "" {
		addr = ":9092"
	}
	return &Server{
		Addr: addr,
		subs: make(map[string]*net.UDPAddr),
	}
}

func (s *Server) Start() error {
	udpAddr, err := net.ResolveUDPAddr("udp", s.Addr)
	if err != nil {
		return fmt.Errorf("resolve udp addr: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("listen udp: %w", err)
	}
	s.conn = conn

	log.Printf("[udp] server listening on %s\n", s.Addr)

	buf := make([]byte, 4096)
	for {
		n, from, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("[udp] read error: %v\n", err)
			continue
		}

		msg := strings.TrimSpace(string(buf[:n]))
		s.handle(from, msg)
	}
}

func (s *Server) handle(from *net.UDPAddr, msg string) {
	upper := strings.ToUpper(msg)

	switch {
	case upper == "SUB":
		s.addSubscriber(from)
		_ = s.reply(from, fmt.Sprintf("OK SUBSCRIBED total=%d", s.count()))

	case upper == "UNSUB":
		s.removeSubscriber(from)
		_ = s.reply(from, fmt.Sprintf("OK UNSUBSCRIBED total=%d", s.count()))

	case upper == "PING":
		_ = s.reply(from, "PONG")

	case upper == "LIST":
		_ = s.reply(from, fmt.Sprintf("OK SUBSCRIBERS=%d", s.count()))

	case strings.HasPrefix(upper, "BROADCAST "):
		body := strings.TrimSpace(msg[len("BROADCAST "):])
		sent := s.broadcast(body)
		_ = s.reply(from, fmt.Sprintf("OK BROADCAST sent=%d", sent))

	default:
		_ = s.reply(from, "ERR unknown command (SUB|UNSUB|PING|LIST|BROADCAST <msg>)")
	}
}

func (s *Server) addSubscriber(addr *net.UDPAddr) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subs[addr.String()] = addr
	log.Printf("[udp] subscribed: %s (total=%d)\n", addr.String(), len(s.subs))
}

func (s *Server) removeSubscriber(addr *net.UDPAddr) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.subs, addr.String())
	log.Printf("[udp] unsubscribed: %s (total=%d)\n", addr.String(), len(s.subs))
}

func (s *Server) count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.subs)
}

func (s *Server) broadcast(message string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.conn == nil {
		return 0
	}

	payload := []byte("NOTIFY " + message)
	sent := 0

	for _, addr := range s.subs {
		_ = s.conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
		if _, err := s.conn.WriteToUDP(payload, addr); err == nil {
			sent++
		}
	}

	log.Printf("[udp] broadcast sent=%d msg=%q\n", sent, message)
	return sent
}

func (s *Server) reply(to *net.UDPAddr, text string) error {
	if s.conn == nil {
		return nil
	}
	_ = s.conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	_, err := s.conn.WriteToUDP([]byte(text), to)
	return err
}
