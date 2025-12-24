package websocket

type Message struct {
	Type     string `json:"type"`    
	Username string `json:"username"`
	Content  string `json:"content"`
}
