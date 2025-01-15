package main

import (
	"auth/pkg/app"
	"auth/pkg/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()

	logger := setUpLogger(config.Env)
	logger.Info("starting app", slog.Any("config", config))

	application := app.NewApp(logger, config.GRPC.Port, config.StoragePath, config.TokenTTL)
	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	logger.Info("stopping application", slog.String("sign", sign.String()))
	application.GRPCServer.Stop()

	logger.Info("application stopped")
}

func setUpLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
