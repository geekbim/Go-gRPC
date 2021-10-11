### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/geekbim/Go-gRPC.git
   ```
2. Generate Protobuf
   ```bash
   sh run-proto.sh
   ```
3. Running PostgreSQL Container
   ```bash
   sh postgre.sh
   ```
4. Running gRPC Server
   ```go
   go run user-management-server/user_management_server.go
   ```
5. Running gRPC Client
   ```go
   go run user-management-client/user_management_client.go
   ```
