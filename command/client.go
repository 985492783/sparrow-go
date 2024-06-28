package command

import (
	"fmt"
	"github.com/985492783/sparrow-go/pkg/remote/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"syscall"
)

type client struct {
	cn   chan os.Signal
	conn *grpc.ClientConn
	req  pb.RequestClient
}

func newConnect(addr string) (*client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	requestClient := pb.NewRequestClient(conn)
	connect := &client{
		cn:   make(chan os.Signal),
		conn: conn,
		req:  requestClient,
	}
	go func() {
		signal.Notify(connect.cn, syscall.SIGINT, syscall.SIGTERM)
		<-connect.cn
		fmt.Println("grpc连接关闭")
		connect.conn.Close()
		os.Exit(0)
	}()
	fmt.Println("grpc连接开启")

	return connect, nil
}
