package main

import (
	"context"
	"github.com/Ling-Qingran/gRPC-Observability/user"
	"google.golang.org/grpc"
	"log"
	"net"
)

type myUserServer struct {
	user.UnimplementedUserServer
}

func (s myUserServer) Create(ctx context.Context, req *user.CreateRequest) (*user.CreateResponse, error) {
	return &user.CreateResponse{
		Csv: []byte("test"),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	service := &myUserServer{}

	user.RegisterUserServer(server, service)
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
