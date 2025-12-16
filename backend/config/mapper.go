package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	LogLevel                  int
	DatabaseHost              string
	DatabaseUser              string
	DatabasePass              string
	DatabasePort              string
	DatabaseName              string
	BackendPort               string
	MigrationsPath            string
	FrontendURL               string
	DBSchema                  string
	AssignmentsWorkerInterval time.Duration
	TaskMinAge                time.Duration
}

func LoadConfig(envPath string) *Config {
	//if err := godotenv.Load(envPath); err != nil {
	//	log.Fatalf("Ошибка загрузки .env: %v", err)
	//}

	logLevel, err := strconv.Atoi(getEnv("LOG_LEVEL", "0"))
	if err != nil {
		logLevel = 0
	}

	intervalStr := getEnv("ASSIGNMENTS_WORKER_INTERVAL", "1h")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Printf("Не удалось распарсить ASSIGNMENTS_WORKER_INTERVAL: %v, используется default 1h", err)
		interval = time.Hour
	}

	taskMinAgeStr := getEnv("ASSIGNMENTS_TASK_MIN_AGE", "24h")
	taskMinAge, err := time.ParseDuration(taskMinAgeStr)
	if err != nil {
		log.Printf("Не удалось распарсить ASSIGNMENTS_TASK_MIN_AGE: %v, используется default 24h", err)
		taskMinAge = 24 * time.Hour
	}

	return &Config{
		LogLevel:                  logLevel,
		DatabaseHost:              getEnv("DATABASE_HOST", "postgres"),
		DatabaseUser:              getEnv("DATABASE_USER", ""),
		DatabasePass:              getEnv("DATABASE_PASSWORD", ""),
		DatabasePort:              getEnv("DATABASE_PORT", "5432"),
		DatabaseName:              getEnv("DATABASE_NAME", "postgres"),
		BackendPort:               getEnv("BACKEND_PORT", "8080"),
		MigrationsPath:            getEnv("MIGRATIONS_PATH", "./migrations"),
		FrontendURL:               getEnv("FRONTEND_URL", ""),
		DBSchema:                  getEnv("DB_SCHEMA", "coopera"),
		AssignmentsWorkerInterval: interval,
		TaskMinAge:                taskMinAge,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
