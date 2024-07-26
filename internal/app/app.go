package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/TiZir/gRPC_service/internal/app/grpc"
	"github.com/TiZir/gRPC_service/internal/services/auth"
	"github.com/TiZir/gRPC_service/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {

	//Хранилище
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	//Сервисный слой
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService,grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
