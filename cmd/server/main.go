package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/ysomad/go-auth-service/internal/app/server"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/local.toml", "path to config file")
}

func main() {
	// Add -config-path flag from init() to server executable
	flag.Parse()

	// Create config instance for using it in server
	config := server.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	// Create server instance and error handling
	if err := server.Start(config); err != nil {
		log.Fatal(err)
	}
}
