package main

import (
	"context"
	"flag"
	"github.com/985492783/sparrow-go/cmd/console"
	"github.com/985492783/sparrow-go/cmd/switcher"
	"github.com/985492783/sparrow-go/pkg/config"
	"github.com/985492783/sparrow-go/pkg/core"
	"github.com/985492783/sparrow-go/pkg/db"
	"log"
	"sync"
)

func main() {
	configPath := flag.String("config", "config.toml", "path to sparrow config file")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading sparrow config: %v", err)
	}
	runServer(cfg)
}

func runServer(cfg *config.SparrowConfig) {
	database, err := db.NewDatabase(cfg)
	switcherManager := core.NewSwitcherManager()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	if cfg.SwitcherConfig.Enabled {
		wg.Add(1)
		go func() {
			switcherS := switcher.NewSwitcherServer(ctx, &wg, cfg)
			err := switcherS.Start(switcherManager, database)
			if err != nil {
				log.Printf("Error starting switcher: %v\n", err)
				cancel()
			}
		}()
	}

	if cfg.ConsoleConfig.Enabled {
		wg.Add(1)
		go func() {
			consoleS := console.NewConsoleServer(ctx, &wg, cfg)
			err := consoleS.Start(switcherManager)
			if err != nil {
				log.Printf("Error starting console: %v\n", err)
				cancel()
			}
		}()
	}
	log.Println("server started")
	wg.Wait()
	log.Println("Shutting down...")
}
