package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Negat1v9/link-checker/config"
	"github.com/Negat1v9/link-checker/internal/app"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	loadConfig, err := config.LoadConfig("./config/config")

	if err != nil {
		panic(err.Error())
	}

	cfg, err := config.ParseConfig(loadConfig)
	if err != nil {
		panic(err.Error())
	}

	app := app.NewApp(cfg)

	app.Run(ctx)

}
