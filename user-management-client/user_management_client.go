package main

import (
	"context"
	"log"
	"time"

	pb "example.com/go-grpc-user-management/_generated"
	"google.golang.org/grpc"
)

const (
	address = "localhost:5001"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect : %v", err)
	}
	defer conn.Close()
	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var newUsers = make(map[string]uint32)
	newUsers["Abim"] = 19
	newUsers["Dhanu"] = 15

	for name, age := range newUsers {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{
			Name: name,
			Age:  age,
		})
		if err != nil {
			log.Fatalf("could not create user : %v", err)
		}
		log.Printf(`
			User Details: 
			Id: %d,
			Name: %s,
			Age: %d
		`, r.GetId(), r.GetName(), r.GetAge())
	}
}
