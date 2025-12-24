package udp

import (
	"encoding/json"
	"net"
)

type Notification struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	MangaID  string `json:"manga_id"`
	Chapter  int    `json:"chapter"`
}

type Client struct {
	conn *net.UDPConn
}

func NewClient(addr string) (*Client, error) {
	ra, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, ra)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Send(n Notification) {
	data, _ := json.Marshal(n)
	_, _ = c.conn.Write(data)
}
