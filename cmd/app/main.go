package main

import (
	"VKTest/internal/server"
	"VKTest/internal/subpub"
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerConfig *server.Config `yaml:"server"`
}

func readConfig(filename string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(filename, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	slog.Info("Reading config")
	cfg, err := readConfig("./cfg.yml")
	if err != nil {
		slog.Error("Can't read config, error: " + err.Error())
		return
	}
	slog.Info("Config read successfully")
	slog.Info("Starting SubPub")
	h := subpub.NewSubPub()
	slog.Info(fmt.Sprintf("SubPub started: %v", h))
	slog.Info("Staring server")
	server, err := server.New(*cfg.ServerConfig, h)
	if err != nil {
		slog.Error("Can't create server, error: " + err.Error())
		return
	}
	if err := server.PBServer.Serve(*server.Listener); err != nil {
		slog.Error("Can't run server, error: " + err.Error())
		return
	} else {
		slog.Info("Server is running")
	}

}
