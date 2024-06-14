package switcher

import (
	"context"
	"github.com/985492783/sparrow-go/pkg/center"
	"github.com/985492783/sparrow-go/pkg/config"
	"github.com/985492783/sparrow-go/pkg/remote/pb"
	"github.com/985492783/sparrow-go/pkg/remote/server"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type SwitcherServer struct {
	cfg *config.SwitcherConfig
	wg  *sync.WaitGroup
	ctx context.Context
}

func NewSwitcherServer(ctx context.Context, wg *sync.WaitGroup, cfg *config.SwitcherConfig) *SwitcherServer {
	return &SwitcherServer{
		ctx: ctx,
		wg:  wg,
		cfg: cfg,
	}
}

func (switcher *SwitcherServer) Start() error {
	defer switcher.wg.Done()

	grpcServer := grpc.NewServer()
	//注册handler
	service := server.NewRequestService()
	service.RegisterHandler(center.NewSwitcherHandler())

	go func() {
		<-switcher.ctx.Done() // 等待停止信号
		grpcServer.GracefulStop()
		log.Println("Switcher server stopped")
	}()

	pb.RegisterRequestServer(grpcServer, service)
	listen, err := net.Listen("tcp", switcher.cfg.Addr)
	if err != nil {
		return err
	}
	log.Printf("Switcher server listening on %s\n", switcher.cfg.Addr)
	//最后通过grpcServer.Serve(listen) 在一个监听端口上提供gRPC服务
	return grpcServer.Serve(listen)
}
