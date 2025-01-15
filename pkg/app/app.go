package app

import (
	grpc_app "auth/pkg/app/grpc"
	"auth/pkg/services/auth"
	"auth/pkg/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpc_app.App
}

func NewApp(logger *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.NewStorage(storagePath)
	if err != nil {
		panic("err")
	}

	authService := auth.NewAuth(storage, storage, storage, tokenTTL, logger)
	grpcApp := grpc_app.NewApp(logger, grpcPort, storagePath, tokenTTL, authService)
	return &App{
		GRPCServer: grpcApp,
	}
}
