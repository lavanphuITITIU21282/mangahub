package udp

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	server *net.UDPAddr
	conn   *net.UDPConn
}

func NewClient(serverAddr string) (*Client, error) {
	if serverAddr == "" {
		serverAddr = "127.0.0.1:9092"
	}

	srv, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("resolve server: %w", err)
	}

	// :0 => random local port
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, fmt.Errorf("listen udp: %w", err)
	}

	return &Client{server: srv, conn: conn}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) SendRaw(cmd string) error {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return fmt.Errorf("empty command")
	}
	_ = c.conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	_, err := c.conn.WriteToUDP([]byte(cmd), c.server)
	return err
}

func (c *Client) Subscribe() error   { return c.SendRaw("SUB") }
func (c *Client) Unsubscribe() error { return c.SendRaw("UNSUB") }
func (c *Client) Ping() error        { return c.SendRaw("PING") }
func (c *Client) List() error        { return c.SendRaw("LIST") }

func (c *Client) Broadcast(msg string) error {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return fmt.Errorf("empty broadcast message")
	}
	return c.SendRaw("BROADCAST " + msg)
}

func (c *Client) Listen(ctx context.Context, onMessage func(text string)) error {
	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		_ = c.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, _, err := c.conn.ReadFromUDP(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				continue
			}
			return err
		}

		text := strings.TrimSpace(string(buf[:n]))
		if onMessage != nil {
			onMessage(text)
		}
	}
}
