package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/lorenas/PaceIt/internal/app"
	"github.com/lorenas/PaceIt/internal/handler"
	"github.com/lorenas/PaceIt/internal/infrastructure"
	"github.com/lorenas/PaceIt/internal/repository"
	"github.com/lorenas/PaceIt/internal/service"
	"github.com/pressly/goose/v3"
)

func main() {
	_ = godotenv.Load()

	db, err := setupDatabase()
    if err != nil {
        log.Fatalf("failed to setup database: %v", err)
    }
    defer db.Close()

	app, err := app.NewApplication(db)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	
	setupUserHandlerRoutes(app, db)

	if err := app.Run(); err != nil {
		log.Fatalf("Application run error: %v", err)
	}
}

func setupDatabase() (*sql.DB, error) {
	db, err := infrastructure.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}
	if err := goose.SetDialect("postgres"); err != nil {
        log.Fatalf("failed to set goose dialect: %v", err)
    }
    if err := goose.Up(db, "internal/migrations"); err != nil {
        log.Fatalf("failed to run migrations: %v", err)
    }
    log.Println("Migrations applied successfully")

	return db, nil
}

func setupUserHandlerRoutes(app app.App, db *sql.DB) {
	userRepo := repository.NewUserRepository(db)
	registerService := service.NewRegisterUserService(userRepo)
	userHandlerInterface := handler.NewUserHandler(registerService)
	userHandlerStruct := userHandlerInterface.(*handler.UserHandler)
	userHandlerStruct.RegisterRoutes(app.Engine())
}
