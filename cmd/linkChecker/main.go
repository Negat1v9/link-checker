package main

import (
	"github.com/Negat1v9/link-checker/config"
	"github.com/Negat1v9/link-checker/internal/app"
)

func main() {
	loadConfig, err := config.LoadConfig("./config/config")

	if err != nil {
		panic(err.Error())
	}

	cfg, err := config.ParseConfig(loadConfig)
	if err != nil {
		panic(err.Error())
	}

	app := app.NewApp(cfg)
	app.Run()
}
