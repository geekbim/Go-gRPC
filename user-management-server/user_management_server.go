package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "example.com/go-grpc-user-management/_generated/user-management"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	port = ":5001"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

type UserManagementServer struct {
	conn *pgx.Conn
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

	createSql := `
		create table if not exists users(
			id		SERIAL 		PRIMARY KEY,
			name	TEXT,
			age		INT
		);
	`

	_, err := s.conn.Exec(context.Background(), createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "table creation failed %v\n", err)
		os.Exit(1)
	}

	createdUser := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
	}
	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn begin failed : %v", err)
	}

	_, err = tx.Exec(context.Background(), "insert into users(name, age) values ($1, $2)",
		createdUser.Name,
		createdUser.Age,
	)
	if err != nil {
		log.Fatalf("tx.Exec failed :%v", err)
	}

	tx.Commit(context.Background())

	return createdUser, nil
}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var userList *pb.UserList = &pb.UserList{}
	rows, err := s.conn.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		userList.Users = append(userList.Users, &user)
	}

	return userList, nil
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("unable to establish connection : %v", err)
	}
	defer conn.Close(context.Background())

	var userManagementServer *UserManagementServer = NewUserManagementServer()
	userManagementServer.conn = conn
	if err := userManagementServer.Run(); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
