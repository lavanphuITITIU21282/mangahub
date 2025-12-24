package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) SetProgress(mangaID, chapter string) {
	body, _ := json.Marshal(map[string]string{
		"manga_id": mangaID,
		"chapter":  chapter,
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

	fmt.Println("Progress updated")
}

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
