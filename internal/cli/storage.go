package cli

import (
	"os"
	"path/filepath"
)

func tokenPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mangahub_token")
}

func SaveToken(token string) {
	_ = os.WriteFile(tokenPath(), []byte(token), 0600)
}

func LoadToken() string {
	data, err := os.ReadFile(tokenPath())
	if err != nil {
		return ""
	}
	return string(data)
}
