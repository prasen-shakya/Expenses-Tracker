package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/prasen-shakya/todo/internal/users"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUsername    = errors.New("username is required")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters")
)

type Service struct {
	users *users.Repository
}

func NewService(usersRepo *users.Repository) *Service {
	return &Service{users: usersRepo}
}

func (s *Service) Register(ctx context.Context, username, password string) (users.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return users.User{}, ErrInvalidUsername
	}

	if len(password) < 8 {
		return users.User{}, ErrInvalidPassword
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return users.User{}, err
	}

	return s.users.Create(ctx, username, string(passwordHash))
}

func (s *Service) Login(ctx context.Context, username, password string) (users.User, error) {
	user, err := s.users.GetByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return users.User{}, ErrInvalidCredentials
		}
		return users.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return users.User{}, ErrInvalidCredentials
	}

	return user, nil
}
