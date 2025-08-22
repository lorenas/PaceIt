package handler_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lorenas/PaceIt/internal/handler"
	"github.com/lorenas/PaceIt/internal/service"
	"github.com/lorenas/PaceIt/internal/repository"
	"github.com/pressly/goose/v3"
)


func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "postgres://paceit_user:paceit_password@localhost:5432/paceit_db?sslmode=disable"
    }
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        t.Fatalf("nepavyko prisijungti prie DB: %v", err)
    }
    if err := db.Ping(); err != nil {
        t.Fatalf("DB ping klaida: %v", err)
    }
    if err := goose.Up(db, "../../migrations"); err != nil {
        t.Fatalf("migracijų klaida: %v", err)
    }
    return db
}

func setupRouter(db *sql.DB) *gin.Engine {
    // Išvalome lentelę prieš kiekvieną testą
    db.Exec("DELETE FROM users")

    userRepo := repository.NewUserRepository(db)
    registerService := service.NewRegisterUserService(userRepo)
    userHandler := handler.NewUserHandler(registerService)

    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.POST("/api/v1/users", userHandler.Register)

    return router
}