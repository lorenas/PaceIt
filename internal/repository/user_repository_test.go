package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lorenas/PaceIt/internal/entity"
)

func newMockRepo(t *testing.T) (UserRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlMock.New: %v", err)
	}
	return NewUserRepository(db), mock
}

func TestRepository_Create(t *testing.T) {
	repo, sqlMock := newMockRepo(t)
	user := &entity.User{
		ID:           uuid.New(),
		Email:        "homer.simpson@simpsons.com",
		PasswordHash: "hash",
	}
	insertUserSql := regexp.QuoteMeta(`INSERT INTO users (id, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`)
	sqlMock.ExpectExec(insertUserSql).
		WithArgs(user.ID, user.Email, user.PasswordHash).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.Create(user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRepository_Create_WhenSameUserEmailExists(t *testing.T) {
	repo, sqlMock := newMockRepo(t)
	user := &entity.User{
		ID:           uuid.New(),
		Email:        "homer.simpson@simpsons.com",
		PasswordHash: "hash",
	}
	insertUserSql := regexp.QuoteMeta(`INSERT INTO users (id, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`)
	sqlMock.ExpectExec(insertUserSql).
		WithArgs(user.ID, user.Email, user.PasswordHash).
		WillReturnError(&pgconn.PgError{Code: "23505"})

	err := repo.Create(user)
	if err != ErrEmailTaken {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestRepository_GetByEmail_WhenUserExists(t *testing.T) {
	repo, mock := newMockRepo(t)
	wantUserEmail := "bart.simpson@simpsons.com"
	selectQuery := regexp.QuoteMeta(`SELECT id,email,password_hash,created_at,updated_at,deleted_at FROM users WHERE email=$1 AND deleted_at IS NULL`)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "email", "password_hash", "created_at", "updated_at", "deleted_at",
	}).AddRow(uuid.New(), wantUserEmail, "hash", now, now, nil)

	mock.ExpectQuery(selectQuery).
		WithArgs(wantUserEmail).
		WillReturnRows(rows)

	got, err := repo.GetByEmail(wantUserEmail)
	if err != nil {
		t.Fatalf("GetByEmail: %v", err)
	}
	if got == nil {
		t.Fatal("expected user, got nil")
	}
	if got.Email != wantUserEmail {
		t.Fatalf("email mismatch got=%s", got.Email)
	}
}
