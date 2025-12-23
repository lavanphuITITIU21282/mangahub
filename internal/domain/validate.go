package domain

import "strings"

func IsValidUsername(username string) bool {
	return len(username) >= 3 && !strings.Contains(username, " ")
}

func IsValidPassword(password string) bool {
	return len(password) >= 8
}

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@")
}
