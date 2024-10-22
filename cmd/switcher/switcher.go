package switcher

import (
	"context"
	"github.com/985492783/sparrow-go/pkg/config"
	"github.com/985492783/sparrow-go/pkg/core"
	"github.com/985492783/sparrow-go/pkg/db"
	"github.com/985492783/sparrow-go/pkg/handler"
	"github.com/985492783/sparrow-go/pkg/remote/pb"
	server2 "github.com/985492783/sparrow-go/pkg/server"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type SwitcherServer struct {
	cfg *config.SparrowConfig
	wg  *sync.WaitGroup
	ctx context.Context
}

func NewSwitcherServer(ctx context.Context, wg *sync.WaitGroup, cfg *config.SparrowConfig) *SwitcherServer {
	return &SwitcherServer{
		ctx: ctx,
		wg:  wg,
		cfg: cfg,
	}
}

func (switcher *SwitcherServer) Start(manager *core.SwitcherManager, database *db.Database) error {
	defer switcher.wg.Done()

	grpcServer := grpc.NewServer()
	//注册handler
	service := server2.NewRequestService(switcher.cfg)
	stream := server2.NewRequestStream(manager)

	service.RegisterHandler(handler.NewSwitcherHandler(database, manager, stream))
	service.RegisterHandler(handler.NewSharkHandler())

	go func() {
		<-switcher.ctx.Done() // 等待停止信号
		grpcServer.GracefulStop()
		log.Println("Switcher server stopped")
	}()

	pb.RegisterRequestServer(grpcServer, service)
	pb.RegisterBiRequestStreamServer(grpcServer, stream)

	listen, err := net.Listen("tcp", switcher.cfg.SwitcherConfig.Addr)
	if err != nil {
		return err
	}
	log.Printf("Switcher server listening on %s\n", switcher.cfg.SwitcherConfig.Addr)
	//最后通过grpcServer.Serve(listen) 在一个监听端口上提供gRPC服务
	return grpcServer.Serve(listen)
}
