package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	userapi "github.com/lorenas/PaceIt/internal/handler"
	"github.com/lorenas/PaceIt/internal/infrastructure"
	"github.com/lorenas/PaceIt/internal/repository"
	"github.com/lorenas/PaceIt/internal/service"
	"github.com/pressly/goose/v3"
)

type App interface {
	Run() error
	Shutdown()
}

type Application struct {
	router *gin.Engine
	db     *sql.DB
	config Config
}

type Config struct {
	Port string
}

func NewApplication() (App, error) {
    app := &Application{
        router: gin.Default(),
    }

    type lifecycleStep struct {
        name string
        fn   func() error
    }

    steps := []lifecycleStep{
        {name: "load env variables", fn: app.loadEnvVariables},
        {name: "setup database", fn: app.setupDatabase},
        {name: "migrate database", fn: app.migrateDatabase},
    }

    for _, s := range steps {
        if err := s.fn(); err != nil {
            return nil, fmt.Errorf("failed to %s: %w", s.name, err)
        }
    }

    app.initialiseRoutes()

    return app, nil
}
func (a *Application) Run() error {
	log.Println("listening on :" + a.config.Port)
	return a.router.Run(":" + a.config.Port)
}

func (a *Application) Shutdown() {
	if a.db != nil {
		log.Println("Closing database connection...")
		if err := a.db.Close(); err != nil {
			log.Printf("Error closing database connection: %v\n", err)
		}
	}
	log.Println("Application shut down.")
}

func (a *Application) loadEnvVariables() error {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	a.config.Port = port
	return nil
}

func (a *Application) setupDatabase() error {
	conn, err := infrastructure.Open()
	if err != nil {
		return err
	}
	a.db = conn
	return nil
}

func (a *Application) migrateDatabase() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose dialect error: %w", err)
	}
	if err := goose.Up(a.db, "internal/migrations"); err != nil {
		return fmt.Errorf("migrations error: %w", err)
	}
	log.Println("Migrations applied successfully")
	return nil
}

func (a *Application) initialiseRoutes() {
	userRepo := repository.NewUserRepository(a.db)
	registerService := service.NewRegisterUserService(userRepo)
	userHandler := userapi.NewUserHandler(registerService)

	api := a.router.Group("/api/v1")
	{
		api.POST("/users", userHandler.Register)
	}
}
