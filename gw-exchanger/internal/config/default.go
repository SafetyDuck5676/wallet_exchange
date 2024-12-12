package config

import "os"

func SetDefaults() {
	_ = os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/exchanger")
	_ = os.Setenv("GRPC_PORT", "50051")
}
