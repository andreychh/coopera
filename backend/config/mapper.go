package config

import (
	"os"
	"strconv"
)

type Config struct {
	LogLevel       int
	DatabaseHost   string
	DatabaseUser   string
	DatabasePass   string
	DatabasePort   string
	DatabaseName   string
	BackendPort    string
	MigrationsPath string
	FrontendURL    string
	DBSchema       string
}

func LoadConfig(envPath string) *Config {
	_ = envPath

	//if err := godotenv.Load(envPath); err != nil {
	//	log.Fatalf("Ошибка загрузки .env: %v", err)
	//}

	logLevel, err := strconv.Atoi(getEnv("LOG_LEVEL", "0"))
	if err != nil {
		logLevel = 0
	}

	return &Config{
		LogLevel:       logLevel,
		DatabaseHost:   getEnv("DATABASE_HOST", "postgres"),
		DatabaseUser:   getEnv("DATABASE_USER", ""),
		DatabasePass:   getEnv("DATABASE_PASSWORD", ""),
		DatabasePort:   getEnv("DATABASE_PORT", "5432"),
		DatabaseName:   getEnv("DATABASE_NAME", "postgres"),
		BackendPort:    getEnv("BACKEND_PORT", "8080"),
		MigrationsPath: getEnv("MIGRATIONS_PATH", "./migrations"),
		FrontendURL:    getEnv("FRONTEND_URL", ""),
		DBSchema:       getEnv("DB_SCHEMA", "coopera"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
