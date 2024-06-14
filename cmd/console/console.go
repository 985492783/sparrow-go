package console

import (
	"context"
	"github.com/985492783/sparrow-go/pkg/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

type ConsoleServer struct {
	cfg *config.ConsoleConfig
	wg  *sync.WaitGroup
	ctx context.Context
}

func NewConsoleServer(ctx context.Context, wg *sync.WaitGroup, cfg *config.ConsoleConfig) *ConsoleServer {
	return &ConsoleServer{
		cfg: cfg,
		wg:  wg,
		ctx: ctx,
	}
}
func (htt *ConsoleServer) Start() error {
	defer htt.wg.Done()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})

	svc := &http.Server{
		Addr:    htt.cfg.Addr,
		Handler: r,
	}
	go func() {
		<-htt.ctx.Done() // 等待停止信号
		log.Println("http server stopped")
		_ = svc.Shutdown(context.Background())
	}()
	log.Printf("Console server listening on %s\n", htt.cfg.Addr)
	return svc.ListenAndServe()
}
