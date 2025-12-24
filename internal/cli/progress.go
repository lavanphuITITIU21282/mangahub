package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (c *Client) SetProgress(mangaID, chapter string) {
	ch, err := strconv.Atoi(chapter)
	if err != nil {
		fmt.Println("Invalid chapter (must be a number):", chapter)
		return
	}
	if ch < 0 {
		fmt.Println("Invalid chapter (must be >= 0):", chapter)
		return
	}

	body, _ := json.Marshal(map[string]any{
		"manga_id": mangaID,
		"chapter":  ch,
	})

	req, _ := http.NewRequest("PUT", c.BaseURL+"/users/progress", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var e map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&e)
		fmt.Println("Error:", e)
		return
	}

	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	fmt.Println("Progress updated:", out)
}

// Deprecated (kept for compatibility with earlier experiments)
func (c *Client) ListLibraryLegacy() {
	req, _ := http.NewRequest("GET", c.BaseURL+"/users/library", nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	var data []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	for _, m := range data {
		fmt.Printf("- %s: chapter %v\n", m["manga_id"], m["current_chapter"])
	}
}
