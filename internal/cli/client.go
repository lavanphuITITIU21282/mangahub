package cli

type Client struct {
	BaseURL string
	Token   string
}

func NewClient(url string) *Client {
	return &Client{
		BaseURL: url,
		Token:   LoadToken(),
	}
}
