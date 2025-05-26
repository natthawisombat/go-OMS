package main

import (
	"context"
	"go-oms/configs"
	pkglogger "go-oms/pkg/logger"
	"go-oms/routers"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	logger := pkglogger.NewLoggerFromEnv()
	ctx = pkglogger.WithLogger(ctx, logger)
	defer stop()

	// app
	app := configs.NewApp(logger)
	if err := app.SetApp(ctx); err != nil {
		logger.Fatal(err)
	}

	routers.SetupRoutes(app)
	errChan := app.RunApp(ctx)

	select {
	case <-ctx.Done():
		logger.Info("shutting down via signal...")
	case err := <-errChan:
		logger.Errorw("server error", "error", err)
	}

	if err := app.StopApp(ctx); err != nil {
		logger.Errorw("shutdown error", "error", err)
	} else {
		logger.Info("shutdown complete")
	}
}
