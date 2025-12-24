package grpc

import (
	"database/sql"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	mangahubv1 "mangahub/proto/mangahubv1"
)

type Server struct {
	Addr      string
	DB        *sql.DB
	JWTSecret string
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	gs := grpc.NewServer()

	mangahubv1.RegisterMangaHubServiceServer(gs, &MangaHubService{
		DB:        s.DB,
		JWTSecret: s.JWTSecret,
	})

	// Enable reflection for grpcurl / debugging.
	reflection.Register(gs)

	log.Println("gRPC server running on", s.Addr)
	return gs.Serve(lis)
}
