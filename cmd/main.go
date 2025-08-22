package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	userapi "github.com/lorenas/PaceIt/internal/handler"
	"github.com/lorenas/PaceIt/internal/infrastructure"
	"github.com/lorenas/PaceIt/internal/repository"
	"github.com/lorenas/PaceIt/internal/service"
)

func main() {
	_ = godotenv.Load()
	conn, err := infrastructure.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	userRepo := repository.NewRepository(conn)
	registerService := service.NewRegisterUserService(userRepo)
	userHandler := userapi.NewUserHandler(registerService)

	router := gin.Default()
	router.POST("/api/v1/users", userHandler.Register)

	log.Println("listening on :" + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
