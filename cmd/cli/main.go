package main

import (
	"fmt"
	"os"

	"mangahub/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mangahub <command>")
		return
	}

	client := cli.NewClient("http://localhost:8080")

	switch os.Args[1] {

	case "login":
		if len(os.Args) != 4 {
			fmt.Println("Usage: mangahub login <username> <password>")
			return
		}
		client.Login(os.Args[2], os.Args[3])

	case "manga":
		if len(os.Args) < 3 {
			fmt.Println("Usage: mangahub manga <list|view>")
			return
		}
		if os.Args[2] == "list" {
			client.ListManga()
		} else if os.Args[2] == "view" && len(os.Args) == 4 {
			client.ViewManga(os.Args[3])
		}

	case "progress":
		if len(os.Args) != 5 || os.Args[2] != "set" {
			fmt.Println("Usage: mangahub progress set <manga_id> <chapter>")
			return
		}
		client.SetProgress(os.Args[3], os.Args[4])

	case "library":
		if len(os.Args) != 3 || os.Args[2] != "list" {
			fmt.Println("Usage: mangahub library list")
			return
		}
		client.ListLibrary()

	default:
		fmt.Println("Unknown command")
	}
}
