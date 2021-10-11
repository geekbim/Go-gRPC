package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "example.com/go-grpc-user-management/_generated/user-management"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port = ":5001"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen %v : ", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("server listening at %v", lis.Addr())

	return s.Serve(lis)
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	readBytes, err := ioutil.ReadFile("users.json")
	var userList *pb.UserList = &pb.UserList{}
	userId := uint32(rand.Intn(1000))
	createdUser := &pb.User{
		Id:   userId,
		Name: in.GetName(),
		Age:  in.GetAge(),
	}

	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("file not found, creating new file")
			userList.Users = append(userList.Users, createdUser)
			jsonBytes, err := protojson.Marshal(userList)
			if err != nil {
				log.Fatalf("json marshaling failed : %v", err)
			}
			if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
				log.Fatalf("failed write to file : %v", err)
			}

			return createdUser, nil
		} else {
			log.Fatalf("error reading file : %v", err)
		}
	}

	if err := protojson.Unmarshal(readBytes, userList); err != nil {
		log.Fatalf("failed to parse user list : %v", err)
	}

	userList.Users = append(userList.Users, createdUser)
	jsonBytes, err := protojson.Marshal(userList)
	if err != nil {
		log.Fatalf("json marshaling failed : %v", err)
	}
	if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
		log.Fatalf("failed write to file : %v", err)
	}

	return createdUser, nil
}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	jsonBytes, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatalf("failed read from file : %v", err)
	}

	var userList *pb.UserList = &pb.UserList{}
	if err := protojson.Unmarshal(jsonBytes, userList); err != nil {
		log.Fatalf("unmarshaling failed : %v", err)
	}

	return userList, nil
}

func main() {
	var userManagementServer *UserManagementServer = NewUserManagementServer()
	if err := userManagementServer.Run(); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
