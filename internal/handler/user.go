package handler

import (
	"mindgames/internal/domain"
	"mindgames/internal/service"
	"mindgames/internal/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthRouteErrorType int

const (
	UserNotFound AuthRouteErrorType = iota
	UserAlreadyExists
	WrongPassword
	DataBaseIsNotInitialized
	Unknown
)

func makeErrorResponse(errorType AuthRouteErrorType) (int, map[string]string) {
	switch errorType {
	case UserNotFound:
		return http.StatusBadRequest, map[string]string{"error": "username not found"}
	case UserAlreadyExists:
		return http.StatusBadRequest, map[string]string{"error": "username already exists"}
	case WrongPassword:
		return http.StatusBadRequest, map[string]string{"error": "invalid password"}
	case DataBaseIsNotInitialized:
		return http.StatusInternalServerError, map[string]string{"error": "database not initialized"}
	default:
		return http.StatusBadRequest, map[string]string{"error": "something went wrong"}
	}
}

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type AuthResponse struct {
	Tokens TokenResponse     `json:"tokens"`
	User   domain.UserPublic `json:"user"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type IUserHandler interface {
	PostRegister(c echo.Context) error
	PostLogin(c echo.Context) error
	PostRefresh(c echo.Context) error
}

type UserHandler struct {
	userService  service.IUserService
	tokenService service.ITokenService
}

func NewUserHandler(userService service.IUserService, tokenService service.ITokenService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (h *UserHandler) PostRegister(c echo.Context) error {
	var req UserCreateRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(makeErrorResponse(Unknown))
	}

	user := domain.User{
		ID:                   uuid.New(),
		Username:             req.Username,
		IsOnboardingComplete: false,
	}

	request := c.Request()
	ctx := request.Context()

	err = h.userService.Create(ctx, &user, req.Password)
	if err != nil {
		return c.JSON(makeErrorResponse(Unknown))
	}

	accessToken, err := h.tokenService.GenerateAccessToken(user)
	if err != nil {
		return c.JSON(makeErrorResponse(Unknown))
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return c.JSON(makeErrorResponse(Unknown))
	}

	resp := AuthResponse{
		Tokens: TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken.Value,
		},
		User: domain.UserPublic{
			Username:             user.Username,
			IsOnboardingComplete: user.IsOnboardingComplete,
		},
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) PostLogin(c echo.Context) error {
	var req UserCreateRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(makeErrorResponse(Unknown))
	}

	user, err := h.userService.GetByUsername(req.Username)
	if err != nil {
		return c.JSON(makeErrorResponse(UserNotFound))
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.JSON(makeErrorResponse(WrongPassword))
	}

	accessToken, err := h.tokenService.GenerateAccessToken(*user)
	if err != nil {
		return c.JSON(makeErrorResponse(Unknown))
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(*user)
	if err != nil {
		return c.JSON(makeErrorResponse(Unknown))
	}

	resp := AuthResponse{
		Tokens: TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken.Value,
		},
		User: domain.UserPublic{
			Username:             user.Username,
			IsOnboardingComplete: user.IsOnboardingComplete,
		},
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) PostRefresh(c echo.Context) error {
	var req RefreshRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	refreshToken, err := h.tokenService.GetRefreshByValue(req.RefreshToken)
	if err != nil {
		return c.JSON(makeErrorResponse(Unknown))
	}

	user := refreshToken.User

	accessToken, err := h.tokenService.GenerateAccessToken(user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate access token"})
	}

	tokenResponse := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Value,
	}

	return c.JSON(http.StatusOK, tokenResponse)
}
