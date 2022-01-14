package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/app"
)

func main() {
	var cfg config.Config

	// TODO: read configuration from flag
	err := cleanenv.ReadConfig("./config/local.yml", &cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(&cfg)
}
