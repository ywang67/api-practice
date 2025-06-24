package main

import (
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	cmClient := NewCableModemsClient(conn)
	cmClient.ByMac([]string{"00:11:22:33:44:55"})
}
