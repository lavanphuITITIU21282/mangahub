package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ListManga() {
	resp, _ := http.Get(c.BaseURL + "/manga")
	defer resp.Body.Close()

	var data []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	for _, m := range data {
		fmt.Printf("- %s (%s)\n", m["title"], m["id"])
	}
}

func (c *Client) ViewManga(id string) {
	resp, _ := http.Get(c.BaseURL + "/manga/" + id)
	defer resp.Body.Close()

	var m map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&m)

	fmt.Printf("Title: %s\nAuthor: %s\n", m["title"], m["author"])
}
