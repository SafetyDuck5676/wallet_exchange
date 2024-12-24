package config

import "os"

func SetDefaults() {
	_ = os.Setenv("DATABASE_URL", "postgres://admin:securepassword@db_server:5432/my_database?sslmode=disable")
	_ = os.Setenv("GRPC_PORT", "50051")
}
