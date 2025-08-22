package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type App interface {
	Engine() *gin.Engine
	Run() error
}

type Application struct {
	router *gin.Engine
	db     *sql.DB
	config Config
}

type Config struct {
	Port string
}

func NewApplication(db *sql.DB) (App, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	app := &Application{
		router: gin.Default(),
		db:     db,
		config: Config{Port: port},
	}

    return app, nil
}

func (app *Application) Engine() *gin.Engine {
	return app.router
}

func (app *Application) Run() error {
	log.Println("listening on :" + app.config.Port)
	return app.router.Run(":" + app.config.Port)
}
