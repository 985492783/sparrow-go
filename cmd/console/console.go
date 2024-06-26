package console

import (
	"context"
	"github.com/985492783/sparrow-go/pkg/config"
	"github.com/985492783/sparrow-go/pkg/utils"
	"github.com/985492783/sparrow-go/pkg/web/controllers"
	"github.com/985492783/sparrow-go/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

type ConsoleServer struct {
	cfg *config.SparrowConfig
	wg  *sync.WaitGroup
	ctx context.Context
}

func NewConsoleServer(ctx context.Context, wg *sync.WaitGroup, cfg *config.SparrowConfig) *ConsoleServer {
	return &ConsoleServer{
		cfg: cfg,
		wg:  wg,
		ctx: ctx,
	}
}
func (server *ConsoleServer) Start() error {
	defer server.wg.Done()
	r := NewEngine(server.cfg)

	svc := &http.Server{
		Addr:    server.cfg.ConsoleConfig.Addr,
		Handler: r,
	}
	go func() {
		<-server.ctx.Done() // 等待停止信号
		log.Println("http server stopped")
		_ = svc.Shutdown(context.Background())
	}()
	log.Printf("Console server listening on %s\n", server.cfg.ConsoleConfig.Addr)
	return svc.ListenAndServe()
}

func NewEngine(cfg *config.SparrowConfig) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.MaxMultipartMemory = 8 << 20

	sparrow := r.Group("sparrow")
	v1 := sparrow.Group("v1")
	{
		switcherList := v1.Group("switcher")
		switcherList.Use(middleware.Auth(utils.AuthSwitcherList, cfg))
		switcherController := controllers.NewSwitcherController(cfg)
		{
			switcherList.GET("/ns", switcherController.QueryNameSpace)
			switcherList.GET("/class", switcherController.QueryClass)
		}
	}
	return r
}
