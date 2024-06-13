package main

import (
	"github.com/985492783/sparrow-go/pkg/center"
	"github.com/985492783/sparrow-go/pkg/remote/pb"
	"github.com/985492783/sparrow-go/pkg/remote/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	grpcServer := grpc.NewServer()
	//注册handler
	service := server.NewRequestService()
	service.RegisterHandler(center.NewSwitcherHandler())

	pb.RegisterRequestServer(grpcServer, service)
	listen, err := net.Listen("tcp", ":9854")
	if err != nil {
		log.Fatal("Listen TCP err:", err)
	}
	//最后通过grpcServer.Serve(listen) 在一个监听端口上提供gRPC服务
	log.Fatal(grpcServer.Serve(listen))
}
