package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/ysomad/go-auth-service/config"
	"github.com/ysomad/go-auth-service/internal/app"
)

func main() {
	var cfg config.Config

	err := cleanenv.ReadConfig("./config/config.yml", &cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(&cfg)
}
