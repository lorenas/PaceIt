package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lorenas/PaceIt/internal/repository"
	"github.com/lorenas/PaceIt/internal/service"
)

type UserHandler interface {
	Register(context *gin.Context)
}
type userHandler struct {
	registerService service.RegisterUserService
}

func NewUserHandler(registerService service.RegisterUserService) UserHandler {
	return &userHandler{
		registerService: registerService,
	}
}

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (handler *userHandler) Register(context *gin.Context) {
	var req createUserRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json"})
		return
	}

	userDto, err := handler.registerService.Register(req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidEmail:
			context.JSON(http.StatusBadRequest, gin.H{"error": "invalid_email"})
		case service.ErrInvalidPassword:
			context.JSON(http.StatusBadRequest, gin.H{"error": "invalid_password"})
		case repository.ErrEmailTaken:
			context.JSON(http.StatusConflict, gin.H{"error": "email_taken"})
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		}
		return
	}

	resp := createUserResponse{
		ID:        userDto.ID.String(),
		Email:     userDto.Email,
		CreatedAt: userDto.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: userDto.UpdatedAt.UTC().Format(time.RFC3339),
	}
	context.JSON(http.StatusCreated, gin.H{"user": resp})
}
