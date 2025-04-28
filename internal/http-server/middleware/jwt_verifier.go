package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	TokenClaimsKey contextKey = "jwt_claims"
	JWTCookieName  string     = "jwt_token"
)

func JWTVerifier(log *slog.Logger, secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ParseTokenFromCookie(w, r)
			if !ok {
				log.Error("Error in token parsing")
				return
			}

			ctx := context.WithValue(r.Context(), TokenClaimsKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ParseTokenFromCookie(w http.ResponseWriter, r *http.Request) (jwt.MapClaims, bool) {
	cookie, err := r.Cookie(JWTCookieName)
	if err != nil {
		http.Error(w, "Authorization cookie is required", http.StatusUnauthorized)
		return nil, false
	}

	return ParseTokenString(cookie.Value)
}

func ParseTokenString(tokenString string) (jwt.MapClaims, bool) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}

	return claims, true
}

func ParseTokenFromRequest(w http.ResponseWriter, r *http.Request) (jwt.MapClaims, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return nil, false
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return nil, false
	}

	return ParseTokenString(tokenParts[1])
}
