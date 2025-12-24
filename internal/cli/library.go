package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type libraryItem struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	CurrentChapter int   `json:"current_chapter"`
	Status        string `json:"status"`
}

func (c *Client) ListLibrary() {
	req, err := http.NewRequest("GET", c.BaseURL+"/users/library", nil)
	if err != nil {
		fmt.Println("request error:", err)
		return
	}

	// ✅ add auth if token exists
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("http error:", err)
		return
	}
	defer resp.Body.Close()

	// ✅ if not 200, print body for debugging
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("server returned %d: %s\n", resp.StatusCode, string(body))
		return
	}

	var items []libraryItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		fmt.Println("decode error:", err)
		return
	}

	for _, it := range items {
		// ✅ display đẹp
		fmt.Printf("- %s (%s): chapter %d [%s]\n", it.Title, it.ID, it.CurrentChapter, it.Status)
	}
}
