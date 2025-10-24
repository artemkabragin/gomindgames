package handler

import (
	"mindgames/internal/domain"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TestHandler interface {
	GetTest(c echo.Context) error
}

type TestHandlerImpl struct {
}

func NewTestHandler() TestHandler {
	return &TestHandlerImpl{}
}

func (h TestHandlerImpl) GetTest(c echo.Context) error {
	tokenIsValid, err := checkToken(c)
	if !tokenIsValid || err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"ok": "true",
	})
}

func checkToken(c echo.Context) (bool, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return false, c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or invalid Authorization header"})
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenString == "" {
		return false, c.JSON(http.StatusUnauthorized, map[string]string{"error": "empty bearer token"})
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me-dev-secret"
	}

	claims := &domain.UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// Ensure HS256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || token == nil || !token.Valid {
		return false, c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
	}

	return true, nil
}
