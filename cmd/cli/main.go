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

	case "grpc":
		// gRPC commands use MANGAHUB_GRPC_ADDR env var (default: localhost:50051)
		if len(os.Args) < 3 {
			fmt.Println("Usage: mangahub grpc <search|get|progress> ...")
			fmt.Println("  mangahub grpc search <query>")
			fmt.Println("  mangahub grpc get <manga_id>")
			fmt.Println("  mangahub grpc progress <manga_id> <chapter>")
			return
		}
		grpcAddr := os.Getenv("MANGAHUB_GRPC_ADDR")
		g := cli.NewGRPCClient(grpcAddr)
		switch os.Args[2] {
		case "search":
			if len(os.Args) != 4 {
				fmt.Println("Usage: mangahub grpc search <query>")
				return
			}
			g.SearchManga(os.Args[3])
		case "get":
			if len(os.Args) != 4 {
				fmt.Println("Usage: mangahub grpc get <manga_id>")
				return
			}
			g.GetManga(os.Args[3])
		case "progress":
			if len(os.Args) != 5 {
				fmt.Println("Usage: mangahub grpc progress <manga_id> <chapter>")
				return
			}
			g.UpdateProgress(os.Args[3], os.Args[4])
		default:
			fmt.Println("Unknown grpc command")
		}

	default:
		fmt.Println("Unknown command")
	}
}
