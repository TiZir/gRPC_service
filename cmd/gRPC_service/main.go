package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/TiZir/gRPC_service/internal/app"
	"github.com/TiZir/gRPC_service/internal/config"
)

const (
	encLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Конфиг
	cfg := config.MustLoad()
	// Логгер
	log := setupLogger(cfg.Env)
	// Приложение (application)
	log.Info("starting application")
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	// Запуск в горутине для логики остановки и для системных сигналов
	go application.GRPCServer.MustRun()
	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT) // Сигналы ОС
	<-stop                                               // Ожидаем до почтупления в канал
	application.GRPCServer.Stop()
	log.Info("application stopped")

	// Запустить gRPC приложения

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case encLocal:
		// Выводить ошибки в текстовом формате с уровнем инфо и выше
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}
	return log
}
