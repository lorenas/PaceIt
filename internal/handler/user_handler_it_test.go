package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lorenas/PaceIt/internal/app"
	"github.com/lorenas/PaceIt/internal/repository"
	"github.com/lorenas/PaceIt/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    testDb, cleanup := setupTestDB(t)
    defer cleanup()

    application, err := app.NewApplication(testDb)
    userRepo := repository.NewUserRepository(testDb)
    registerService := service.NewRegisterUserService(userRepo)
    userHandlerInterface := NewUserHandler(registerService)

    userHandlerStruct := userHandlerInterface.(*UserHandler)
    userHandlerStruct.RegisterRoutes(application.Engine())

    if err != nil {
        t.Fatalf("failed to start application: %v", err)
    }

    requestBody := map[string]string{
        "email":    "integration.test@example.com",
        "password": "password123",
    }
    jsonBody, err := json.Marshal(requestBody)
    require.NoError(t, err, "failed to marshal request body")

    req, err := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(jsonBody))
    require.NoError(t, err, "failed to create request")
    req.Header.Set("Content-Type", "application/json")

    rec := httptest.NewRecorder()
    application.Engine().ServeHTTP(rec, req)

    assert.Equal(t, http.StatusCreated, rec.Code, "unexpected status code")

    var responseBody map[string]map[string]interface{}
    err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
    require.NoError(t, err, "failed to unmarshal response body")
    assert.Equal(t, "integration.test@example.com", responseBody["user"]["email"])

    var email string
    err = testDb.QueryRow("SELECT email FROM users WHERE email = $1", "integration.test@example.com").Scan(&email)
    require.NoError(t, err, "user should have been created in the database")
    assert.Equal(t, "integration.test@example.com", email)
}
