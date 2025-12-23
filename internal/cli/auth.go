package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) Login(username, password string) {
	body, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	resp, err := http.Post(c.BaseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	var res map[string]string
	json.NewDecoder(resp.Body).Decode(&res)

	token := res["token"]
	if token == "" {
		fmt.Println("Login failed")
		return
	}

	SaveToken(token)
	fmt.Println("Login successful")
}
