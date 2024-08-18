package app

import (
	"context"
	"errors"
	"github.com/asb1302/innopolis_go_chat/internal/config"
	"github.com/asb1302/innopolis_go_chat/internal/handler"
	"github.com/asb1302/innopolis_go_chat/internal/repository/cache"
	"github.com/asb1302/innopolis_go_chat/internal/server"
	"github.com/asb1302/innopolis_go_chat/internal/service"
	"github.com/asb1302/innopolis_go_chat/pkg/authclient"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
)

func Run() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup

	chatDB, err := cache.ChatCacheInit(ctx, &wg)
	if err != nil {
		log.Fatalf("ERROR failed to initialize chat database: %v", err)
	}

	config.InitConfig()
	cfg := config.GetConfig()

	// initialize service
	service.Init(chatDB)
	authclient.Init(cfg.AuthServiceHost, cfg.AuthServiceTLS)

	go func() {
		err := server.Run("0.0.0.0", "8000", http.HandlerFunc(handler.HandleHTTPReq))
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("ERROR server run ", err)
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(); err != nil {
		log.Fatal("ERROR server shutdown ", err)
	}
	wg.Wait()
	log.Println("INFO chat service was gracefully shutdown")
}
