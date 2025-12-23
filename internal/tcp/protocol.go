package tcp

type BaseMessage struct {
	Type string `json:"type"`
}

type ProgressUpdateMessage struct {
	Type    string `json:"type"`
	Token   string `json:"token"`
	MangaID string `json:"manga_id"`
	Chapter int    `json:"chapter"`
}

type ProgressBroadcast struct {
	Type      string `json:"type"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	MangaID   string `json:"manga_id"`
	Chapter   int    `json:"chapter"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}
