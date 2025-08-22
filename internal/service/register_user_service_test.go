package service

import (
	"errors"
	"testing"

	"github.com/lorenas/PaceIt/internal/entity"
	"github.com/lorenas/PaceIt/internal/repository"
)

type mockUserRepository struct {
	getByEmailFunc func(email string) (*entity.User, error)
    createFunc     func(user *entity.User) error
}

func (mockRepo *mockUserRepository) GetByEmail(email string) (*entity.User, error) {
	if mockRepo.getByEmailFunc != nil {
        return mockRepo.getByEmailFunc(email)
    }
    return nil, errors.New("getByEmailFunc not implemented")
}

func (mockRepo *mockUserRepository) Create(user *entity.User) error {
	if mockRepo.createFunc != nil {
        return mockRepo.createFunc(user)
    }
    return errors.New("createFunc not implemented")
}

func TestRegisterUserService_WhenSuccess(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByEmailFunc: func(email string) (*entity.User, error) {
			return nil, nil
		},
		createFunc: func(user *entity.User) error {
			return nil
		},
	}

	registerUserService := NewRegisterUserService(mockRepo)

	email := "homer.simpson@simpsons.com"
	password := "dohdohdoh!"
	user, err := registerUserService.Register(email, password)

	if err != nil {
        t.Errorf("Register() error = %v, want err == nil", err)
    }
    if user == nil {
        t.Fatal("expected user to be created, but got nil")
    }
    if user.Email != email {
        t.Errorf("expected email %s, got %s", email, user.Email)
    }
    if user.PasswordHash == "" || user.PasswordHash == password {
        t.Error("password should be hashed")
    }
}

func TestRegisterUserService_WhenPasswordIsShort(t *testing.T) {
	mockRepo := &mockUserRepository{}
	registerUserService := NewRegisterUserService(mockRepo)

	email := "homer.simpson@simpsons.com"
	password := "dohdoh"
	_, err := registerUserService.Register(email, password)

	if !errors.Is(err, ErrInvalidPassword) {
        t.Errorf("Register() error = %v, want err == ErrInvalidPassword", err)
    }
}

func TestRegisterUserService_WhenEmailIsTaken(t *testing.T) {
	mockRepo := &mockUserRepository{
		getByEmailFunc: func(email string) (*entity.User, error) {
			return &entity.User{Email: email}, nil
		},
		createFunc: func(user *entity.User) error {
			return nil
		},
	}

	registerUserService := NewRegisterUserService(mockRepo)

	email := "homer.simpson@simpsons.com"
	password := "dohdohdoh!"
	user, err := registerUserService.Register(email, password)

	if !errors.Is(err, repository.ErrEmailTaken) {
		t.Errorf("Register() error = %v, want err == ErrEmailTaken", err)
	}
	if user != nil {
		t.Fatal("User should not be created with taken email")
	}
}