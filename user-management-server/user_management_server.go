package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "example.com/go-grpc-user-management/_generated/user-management"
	"google.golang.org/grpc"
)

const (
	port = ":5001"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	userId := uint32(rand.Intn(1000))

	return &pb.User{
		Id:   userId,
		Name: in.GetName(),
		Age:  in.GetAge(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen %v : ", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, &UserManagementServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
