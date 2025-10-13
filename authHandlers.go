package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserPublic struct {
	Username             string `json:"username"`
	IsOnboardingComplete string `json:"isOnboardingComplete"`
}

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string     `json:"token"`
	User  UserPublic `json:"user"`
}

func postRegister(c echo.Context) error {
	var req UserCreateRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, req.Password)
}

func postLogin(c echo.Context) error {
	return c.String(http.StatusOK, "test LOOOOL")
}

func postRefresh(c echo.Context) error {
	return c.String(http.StatusOK, "test LOOOOL")
}
