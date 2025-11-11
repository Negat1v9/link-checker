package app

import (
	"context"
	"time"

	"github.com/Negat1v9/link-checker/config"
	"github.com/Negat1v9/link-checker/internal/linkChecker/linkstore"
	"github.com/Negat1v9/link-checker/internal/linkChecker/service"
	"github.com/Negat1v9/link-checker/internal/server"
	"github.com/Negat1v9/link-checker/pkg/logger"
)

const (
	shutDownPeriodSeconds  time.Duration = 5 * time.Second
	serverShutDownDeadLine time.Duration = 10 * time.Second
)

type App struct {
	cfg *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run(shutDown context.Context) {

	logger := logger.NewLogger(a.cfg.Env)

	server := server.NewServer(a.cfg)

	linkStore, err := linkstore.NewLinkStore(&a.cfg.LinkStoreCfg, logger)
	if err != nil {
		logger.Errorf("creating linkStore %v", err)
		panic(1)
	}

	linkService := service.NewLinkService(linkStore)

	server.MapHandlers(linkService)

	logger.Infof("run app")
	go func() {
		if err := server.Run(); err != nil {
			logger.Warnf("http server is stopped: %v", err)
		}
	}()

	// waiting shutdown
	<-shutDown.Done()

	logger.Warnf("start shutDown")
	ctx, cancel := context.WithTimeout(context.Background(), serverShutDownDeadLine)
	// 1 step is stopped server
	if err = server.Stop(ctx); err != nil {
		logger.Errorf("server.Stop: %v", err)
	}
	cancel()
	ctx, cancel = context.WithTimeout(context.Background(), shutDownPeriodSeconds)
	if err = linkStore.Stop(ctx); err != nil {
		logger.Errorf("linkStore.Stop: %v", err)
	}
	cancel()

	logger.Infof("shutdown completed")
}
