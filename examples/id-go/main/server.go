/*
- Steps to init a go module
- go mod init example.com/main
- go mod edit -replace example.com/idsvc=../proto
- import example.com/idsvc
- go mod tidy
- proto command: protoc --proto_path=. --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative idsvc.proto
*/

package main

import (
	"fmt"
	"net"
	"log"
	"context"
	"sync"

	"google.golang.org/grpc"
	"example.com/idsvc"
)

// Dummy declaration
type server struct {
	idsvc.UnimplementedIdsvcServer
}

var global_cnt int64 = 0
var m sync.Mutex

func (s *server) GetId(ctx context.Context, in *idsvc.Request) (*idsvc.Response, error) {
	m.Lock()
	global_cnt++
	m.Unlock()
	return &idsvc.Response{Id:global_cnt}, nil
}

func main() {
	fmt.Println("Starting server...")
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("Failed to listen: ", err, lis)
	}

	grpc_serv := grpc.NewServer()
	idsvc.RegisterIdsvcServer(grpc_serv, &server{})

	// Start listening for connections.
	if err := grpc_serv.Serve(lis); err != nil {
		log.Fatal("Failed to serve; ", err)
	}
}