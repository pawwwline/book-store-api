package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "book-store-api/docs"
	"book-store-api/internal/app"
	"book-store-api/internal/config"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	application, err := app.BuildApp(cfg)
	if err != nil {
		fmt.Println("Error initializing application", err)

		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx, cfg.Cache); err != nil {
		fmt.Println("Error starting application", err)

		return
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {

		return
	}
}
