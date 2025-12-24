package main

import (
	"fmt"
	"os"

	"mangahub/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// Allow overriding API base URL from env
	baseURL := os.Getenv("MANGAHUB_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	client := cli.NewClient(baseURL)

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
		switch os.Args[2] {
		case "list":
			client.ListManga()
		case "view":
			if len(os.Args) != 4 {
				fmt.Println("Usage: mangahub manga view <manga_id>")
				return
			}
			client.ViewManga(os.Args[3])
		default:
			fmt.Println("Usage: mangahub manga <list|view>")
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
		if grpcAddr == "" {
			grpcAddr = "localhost:50051"
		}

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
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: mangahub <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  login <username> <password>")
	fmt.Println("  manga list")
	fmt.Println("  manga view <manga_id>")
	fmt.Println("  progress set <manga_id> <chapter>")
	fmt.Println("  library list")
	fmt.Println("  grpc search <query>")
	fmt.Println("  grpc get <manga_id>")
	fmt.Println("  grpc progress <manga_id> <chapter>")
	fmt.Println("")
	fmt.Println("Env:")
	fmt.Println("  MANGAHUB_BASE_URL   (default: http://localhost:8080)")
	fmt.Println("  MANGAHUB_GRPC_ADDR  (default: localhost:50051)")
}
