package app

import (
	"log/slog"
	"os"

	"github.com/Negat1v9/link-checker/config"
	"github.com/Negat1v9/link-checker/internal/linkChecker/service"
	"github.com/Negat1v9/link-checker/internal/server"
)

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	server := server.NewServer(a.cfg)

	linkService := service.NewLinkService()

	server.MapHandlers(linkService)

	logger.Info("run servre on")
	if err := server.Run(); err != nil {
		logger.Warn("server is stopped")
	}
}
