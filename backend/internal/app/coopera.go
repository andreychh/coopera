package app

import (
	"context"
	"fmt"
	"github.com/andreychh/coopera-backend/internal/usecase/task"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/andreychh/coopera-backend/internal/adapter/controller/web_api"
	repomembership "github.com/andreychh/coopera-backend/internal/adapter/repository/membership_repo"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/postgres"
	"github.com/andreychh/coopera-backend/internal/adapter/repository/postgres/dao"
	repotask "github.com/andreychh/coopera-backend/internal/adapter/repository/task_repo"
	repoteams "github.com/andreychh/coopera-backend/internal/adapter/repository/team_repo"
	repouser "github.com/andreychh/coopera-backend/internal/adapter/repository/user_repo"
	"github.com/andreychh/coopera-backend/internal/usecase/memberships"
	"github.com/andreychh/coopera-backend/internal/usecase/team"
	"github.com/andreychh/coopera-backend/internal/usecase/user"
	"github.com/andreychh/coopera-backend/pkg/logger"
	"github.com/andreychh/coopera-backend/pkg/migrator"
	"github.com/go-playground/validator/v10"
)

func Start() error {
	// для локалки "github.com/joho/godotenv"
	//if err := godotenv.Load("config/dev/.env"); err != nil {
	//	return fmt.Errorf("error loading .env: %w", err)
	//}

	logLevel, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if logLevel == 0 {
		logLevel = logger.INFO
	}
	//logService := logger.NewLogger(logLevel)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL not set")
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if err := migrator.Migrate(migrationsPath, dsn, os.Getenv("DB_SCHEMA")); err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	db, err := postgres.NewDB(dsn)
	if err != nil {
		return err
	}

	validate := validator.New()
	web_api.InitValidator(validate)

	userRepo := repouser.NewUserRepository(*dao.NewUserDAO(db))
	teamRepo := repoteams.NewTeamRepository(*dao.NewTeamDAO(db))
	taskRepo := repotask.NewTaskRepository(*dao.NewTaskDAO(db))
	memberRepo := repomembership.NewMembershipRepository(*dao.NewMembershipDAO(db))

	userUC := user.NewUserUsecase(userRepo, db)
	memberUC := memberships.NewMembershipsUsecase(memberRepo, db)
	teamUC := team.NewTeamUsecase(teamRepo, memberUC, db)
	taskUC := task.NewTaskUsecase(taskRepo, memberUC, db)

	router := web_api.NewRouter(userUC, teamUC, taskUC, memberUC).SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:        ":" + port,
		Handler:     router,
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
	log.Println("HTTP server started on port", os.Getenv("PORT"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down...")
	return srv.Shutdown(ctx)
}
