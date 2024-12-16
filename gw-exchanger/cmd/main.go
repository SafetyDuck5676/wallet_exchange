package main

import (
	"gw-exchanger/internal/config"
	"gw-exchanger/internal/storages/postgres"
	utils "gw-exchanger/pkg"

	grpc "gw-exchanger/internal/server"

	"log"
)

func main() {
	utils.LogInfo("Starting gw-exchanger service...")
	config.SetDefaults()

	cfg, err := config.LoadConfig("config.env")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	storage, err := postgres.NewPostgresStorage(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()

	server := grpc.NewServer(storage)
	if err := server.Start(cfg.GRPCPort); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
