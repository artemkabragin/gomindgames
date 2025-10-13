package main

import (
	"errors"
	"mindgames/utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthRouteErrorType int

const (
	UserNotFound AuthRouteErrorType = iota
	WrongPassword
	DataBaseIsNotInitialized
	Unknown
)

type User struct {
	ID                   uuid.UUID `json:"id"`
	Username             string    `json:"username"`
	PasswordHash         string    `json:"passwordHash"`
	IsOnboardingComplete bool      `json:"isOnboardingComplete"`
}

func makeErrorResponse(errorType AuthRouteErrorType) (int, map[string]string) {
	switch errorType {
	case UserNotFound:
		return http.StatusBadRequest, map[string]string{"error": "username not found"}
	case WrongPassword:
		return http.StatusBadRequest, map[string]string{"error": "invalid password"}
	case DataBaseIsNotInitialized:
		return http.StatusInternalServerError, map[string]string{"error": "database not initialized"}
	default:
		return http.StatusBadRequest, map[string]string{"error": "something went wrong"}
	}
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshToken struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `json:"userId" gorm:"type:uuid;index"`
	Value  string    `json:"value"`
}

func (r *RefreshToken) BeforeCreate(*gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type UserPublic struct {
	Username             string `json:"username"`
	IsOnboardingComplete bool   `json:"isOnboardingComplete"`
}

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Tokens TokenResponse `json:"tokens"`
	User   UserPublic    `json:"user"`
}

type UserClaims struct {
	UserID   string `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func generateToken(user User) (TokenResponse, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me-dev-secret"
	}

	now := time.Now()
	claims := UserClaims{
		UserID:   user.ID.String(),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return TokenResponse{}, err
	}

	refresh := RefreshToken{
		ID:     uuid.New(),
		UserID: user.ID,
		Value:  uuid.NewString(),
	}

	if err := db.Create(&refresh).Error; err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refresh.Value,
	}, nil
}

func postRegister(c echo.Context) error {
	if db == nil {
		return c.String(http.StatusInternalServerError, "database not initialized")
	}
	var req UserCreateRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	// Check if a user with the same username already exists
	if err := db.First(&User{}, "username = ?", req.Username).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		return errors.New("username already taken")
	}

	passwordHash, err := utils.HashPassword(req.Password)

	if err != nil {
		return err
	}

	user := User{
		ID:                   uuid.New(),
		Username:             req.Username,
		PasswordHash:         passwordHash,
		IsOnboardingComplete: false,
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	tokens, err := generateToken(user)
	if err != nil {
		return err
	}

	resp := AuthResponse{
		Tokens: tokens,
		User: UserPublic{
			Username:             user.Username,
			IsOnboardingComplete: user.IsOnboardingComplete,
		},
	}

	return c.JSON(http.StatusOK, resp)
}

func postLogin(c echo.Context) error {
	if db == nil {
		return c.String(http.StatusInternalServerError, "database not initialized")
	}

	var req UserCreateRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(makeErrorResponse(Unknown))
	}

	var user User

	if err := db.First(&user, "username = ?", req.Username).Error; err != nil {
		return c.JSON(makeErrorResponse(UserNotFound))
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.JSON(makeErrorResponse(WrongPassword))
	}

	tokens, err := generateToken(user)

	if err != nil {
		return c.JSON(makeErrorResponse(Unknown))
	}

	resp := AuthResponse{
		Tokens: tokens,
		User: UserPublic{
			Username:             user.Username,
			IsOnboardingComplete: user.IsOnboardingComplete,
		},
	}

	return c.JSON(http.StatusOK, resp)
}

func postRefresh(c echo.Context) error {
	return c.String(http.StatusOK, "test LOOOOL")
}
