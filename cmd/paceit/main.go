package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	userapi "github.com/lorenas/PaceIt/internal/api"
	"github.com/lorenas/PaceIt/internal/db"
	"github.com/lorenas/PaceIt/internal/user"
)

func main() {
    _ = godotenv.Load()
    conn, err := db.Open()
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    userRepo := user.NewRepository(conn)
    registerService := user.NewRegisterUserService(userRepo)
    userHandler := userapi.NewUserHandler(registerService)

    router := gin.Default()
    router.POST("/api/v1/users", userHandler.Register)

    log.Println("listening on :" + port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}
