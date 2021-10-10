### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/geekbim/Go-gRPC.git
   ```
2. Generate Protobuf
   ```sh
   sh run-proto.sh
   ```
3. Running gRPC Server
   ```go
   go run user-management-server/user_management_server.go
   ```
4. Running gRPC Client
   ```go
   go run user-management-client/user_management_client.go
   ```