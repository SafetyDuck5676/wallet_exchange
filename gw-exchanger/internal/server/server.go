package grpc

import (
	"gw-exchanger/internal/storages"
	"log"
	"net"

	pb "github.com/SafetyDuck5676/grpc_duck/proto-exchange/grpc/pb"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedExchangeServiceServer // Встраиваем необходимую структуру
	storage                               storages.Storage
}

func NewServer(storage storages.Storage) *Server {
	return &Server{storage: storage}
}

func (s *Server) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterExchangeServiceServer(grpcServer, s)
	log.Printf("gRPC server is running on port %s", port)
	return grpcServer.Serve(listener)
}
