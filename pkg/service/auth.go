package service

import (
	"crypto/sha1"
	todo "do-app"
	"do-app/pkg/repository"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	salt       = "jvhbfdvoivdfjn2343"
	signingKey = "iubvrdie0werfsdfjk'r9983i"
	tokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {

		return "", fmt.Errorf("generate token: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{time.Now().Add(tokenTTL)},
			IssuedAt:  &jwt.NumericDate{time.Now()},
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(tokenSting string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenSting, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("Parse token: %w", err)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, fmt.Errorf("token claim are not of type *tokenClaims")
	}
	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
