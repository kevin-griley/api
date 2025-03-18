package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/kevin-griley/api/internal/handlers"
	"github.com/kevin-griley/api/internal/types"
)

func ExtractBearerToken(authHeader string) (string, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("authorization header must be in Bearer format")
	}
	return strings.TrimSpace(strings.TrimPrefix(authHeader, prefix)), nil
}

func JwtAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenStr, err := ExtractBearerToken(authHeader)
		if err != nil {
			slog.Error("JwtAuthMiddleware", "ExtractBearerToken", err)
			PermissionDenied(w)
			return
		}

		token, err := ValidateJWT(tokenStr)
		if err != nil || !token.Valid {
			slog.Error("JwtAuthMiddleware", "ValidateJWT", err)
			PermissionDenied(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.Error("JwtAuthMiddleware", "token.Claims", err)
			PermissionDenied(w)
			return
		}

		subject, err := claims.GetSubject()
		if err != nil {
			slog.Error("JwtAuthMiddleware", "claims.GetSubject", err)
			PermissionDenied(w)
			return
		}

		userID, err := uuid.Parse(subject)
		if err != nil {
			PermissionDenied(w)
			return
		}

		ctx := context.WithValue(r.Context(), types.ContextKeyUserID, userID)
		ctx = context.WithValue(ctx, types.ContextKeyClaims, claims)
		next(w, r.WithContext(ctx))

	}
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func PermissionDenied(w http.ResponseWriter) {
	handlers.WriteJSON(w, http.StatusUnauthorized, handlers.ApiError{
		Status:  http.StatusUnauthorized,
		Message: "Permission Denied",
	})
}
