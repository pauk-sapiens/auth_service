package grpc_app

import (
	"auth/pkg/grpc/auth"
	auth_grpc "auth/pkg/grpc/auth"
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
)

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(logger *slog.Logger, port int, storagePath string, tokenTTL time.Duration, auth auth.Auth) *App {
	gRPCServer := grpc.NewServer()

	auth_grpc.Register(gRPCServer, auth)
	return &App{
		logger:     logger,
		gRPCServer: gRPCServer,
		port:       port,
	}

}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	const op = "grpcapp.Run"

	log := a.logger.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.logger.With(slog.String("op", op)).Info("stopping grpc server")

	a.gRPCServer.GracefulStop()
}
