package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prasen-shakya/todo/internal/users"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUsername    = errors.New("username is required")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters")
)

type Service struct {
	users     *users.Repository
	secretKey []byte
}

func NewService(usersRepo *users.Repository, jwtSecretKey []byte) *Service {
	return &Service{users: usersRepo, secretKey: jwtSecretKey}
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
	if strings.TrimSpace(username) == "" || password == "" {
		return users.User{}, ErrInvalidCredentials
	}

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

func (s *Service) CreateJwtToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) VerifyJwtToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secretKey, nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}

	userIdFloat, ok := claims["userId"].(float64)
	if !ok {
		return 0, fmt.Errorf("missing userId")
	}

	return int(userIdFloat), nil
}
