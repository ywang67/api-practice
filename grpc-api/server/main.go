package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"api-project/grpc-api/gen/cablemodems"
	"api-project/grpc-api/methods"
)

func main() {
	// 监听 TCP 端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建 gRPC server 实例
	grpcServer := grpc.NewServer()

	// 注册 CableModemService
	cablemodems.RegisterCableModemServiceServer(grpcServer, &methods.CableModemMethod{})

	log.Println("gRPC server listening on :50051")
	// 启动 server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
