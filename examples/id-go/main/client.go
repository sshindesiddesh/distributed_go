package main

import (
	"fmt"
	"log"
	"context"
	"time"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"example.com/idsvc"
)


func GetIdFromServer(conn *grpc.ClientConn, wg *sync.WaitGroup) int64 {
	c := idsvc.NewIdsvcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := c.GetId(ctx, &idsvc.Request{})
	if err != nil {
		log.Fatalf("Error when calling GetId", err)
	}

	defer wg.Done()

	return response.Id
}


func main() {
	conn, err := grpc.Dial("localhost:9000",
						   grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	for ii := 0; ii < 10; ii++ {
		wg.Add(1)
		go fmt.Println("Id: ", GetIdFromServer(conn, &wg))
	}

	wg.Wait()
}