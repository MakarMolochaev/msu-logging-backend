package jwtservice

import (
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtService struct {
	log *slog.Logger
}

func New(
	log *slog.Logger,
) *JwtService {
	return &JwtService{
		log: log,
	}
}

func (a *JwtService) GenerateToken(taskId int32, tokenTTL time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(tokenTTL).Unix(),
		"iat":    time.Now().Unix(),
		"taskId": taskId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}
