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

type RegisterUserService struct {
	repo *repository.UserRepository
}

func NewRegisterUserService(repo *repository.UserRepository) *RegisterUserService {
	return &RegisterUserService{
		repo: repo,
	}
}

func (service *RegisterUserService) Register(email, password string) (*entity.User, error) {
	email = strings.TrimSpace(email)

	if err := service.isValidEmail(email); err != nil {
		return nil, err
	}

	if err := service.isValidPassword(password); err != nil {
		return nil, err
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

func (service *RegisterUserService) isValidEmail(email string) error {
	isValid := strings.Contains(email, "@") && strings.Contains(email, ".") && email != ""
	if !isValid {
		return ErrInvalidEmail
	}
	existing, err := service.repo.GetByEmail(email)
	if err != nil {
		return err
	}
	if existing != nil {
		return repository.ErrEmailTaken
	}
	return nil
}

func (service *RegisterUserService) isValidPassword(password string) error {
	if len(password) >= 8 {
		return nil
	}
	return ErrInvalidPassword
}
