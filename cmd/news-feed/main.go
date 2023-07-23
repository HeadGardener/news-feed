package main

import (
	"context"
	"flag"
	"github.com/HeadGardener/news-feed/internal/configs"
	"github.com/HeadGardener/news-feed/internal/fetcher"
	"github.com/HeadGardener/news-feed/internal/handlers"
	"github.com/HeadGardener/news-feed/internal/sender"
	"github.com/HeadGardener/news-feed/internal/server"
	"github.com/HeadGardener/news-feed/internal/services"
	"github.com/HeadGardener/news-feed/internal/storage"
	"github.com/pkg/errors"
	"log"
	"os/signal"
	"syscall"
	"time"
)

var confPath = flag.String("conf-path", "./config/.env", "path to config env")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conf := configs.MustInit(*confPath)

	db, err := storage.NewDB(conf.DBConfig)
	if err != nil {
		log.Fatalf("[FATAL] failed to establish db connection: %e", err)
	}

	var (
		sourceStorage  = storage.NewSourceStorage(db)
		articleStorage = storage.NewArticleStorage(db)
		userStorage    = storage.NewUserStorage(db)
	)

	var (
		sourceService  = services.NewSourceService(sourceStorage)
		articleService = services.NewArticleService(articleStorage)
		userService    = services.NewUserService(userStorage)
		tokenService   = services.NewTokenService(userStorage)
		fetcher        = fetcher.NewFetcher(sourceStorage, articleStorage)
		sender         = sender.NewSender(userStorage)
	)

	handler := handlers.NewHandler(sourceService, articleService, userService, tokenService)

	srv := &server.Server{}
	go func() {
		if err := srv.Run(conf.ServerConfig.ServerPort, handler.InitRoutes()); err != nil {
			log.Printf("[ERROR] failed to run server: %e", err)
		}
	}()

	log.Println("[INFO] server start working")

	go func(ctx context.Context) {
		if err := fetcher.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to start fetcher: %e", err)
			}
		}
		log.Println("[INFO] fetcher stop working")
	}(ctx)

	log.Println("[INFO] fetcher start working")

	go func(ctx context.Context) {
		if err := sender.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to start sender: %e", err)
			}
		}
		log.Println("[INFO] sender stop working")
	}(ctx)

	log.Println("[INFO] sender start working")

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[INFO] server forced to shutdown: %e", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("[INFO] db connection forced to shutdown: %e", err)
	}

	log.Println("[INFO] server exiting")
}
