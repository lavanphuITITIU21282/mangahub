package udp

import (
	"encoding/json"
	"log"
	"net"
)

type Server struct {
	conn *net.UDPConn
}

func NewServer(addr string) (*Server, error) {
	ua, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", ua)
	if err != nil {
		return nil, err
	}

	return &Server{conn: conn}, nil
}

func (s *Server) Close() error {
	return s.conn.Close()
}

func (s *Server) Run() {
	buf := make([]byte, 1024)

	for {
		n, _, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		var msg map[string]any
		_ = json.Unmarshal(buf[:n], &msg)
		log.Println("UDP notification:", msg)
	}
}
