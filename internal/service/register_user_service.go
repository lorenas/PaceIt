package service

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lorenas/PaceIt/internal/entity"
	"github.com/lorenas/PaceIt/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)

type RegisterUserService interface {
	Register(email, password string) (*entity.User, error)
}

type registerUserService struct {
	repo repository.UserRepository
}

func NewRegisterUserService(repo repository.UserRepository) RegisterUserService {
	return &registerUserService{
		repo: repo,
	}
}

func (service *registerUserService) Register(email, password string) (*entity.User, error) {
	email = strings.TrimSpace(email)

    if !isValidEmail(email) {
        return nil, ErrInvalidEmail
    }
    if !isValidPassword(password) {
        return nil, ErrInvalidPassword
    }

	existing, err := service.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, repository.ErrEmailTaken
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	newUser := &entity.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashBytes),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := service.repo.Create(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func isValidEmail(email string) bool {
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isValidPassword(password string) bool {
    return len(password) >= 8
}