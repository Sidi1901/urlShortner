package handler

import (
	"errors"
	"net/http"

	"github.com/Sidi1901/urlShortner/internal/dto"
	errs "github.com/Sidi1901/urlShortner/internal/errors"
	"github.com/Sidi1901/urlShortner/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type userHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *userHandler {
	return &userHandler{service: service}
}

var validate = validator.New()

func (h *userHandler) RegisterPublicRoutes(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/signup", h.Signup)
		user.POST("/login", h.Login)
		user.POST("/refresh", h.RefreshToken)
	}
}

func (h *userHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {}

func (h *userHandler) Signup(c *gin.Context) {
	ctx := c.Request.Context()
	var user dto.CreateUserRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateUser(ctx, user.Email, user.Name, user.Password, user.UserType)

	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, errs.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	return
}

func (h *userHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CreateLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, refreshToken, err := h.service.Login(ctx, req.Email, req.Password)

	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		case errors.Is(err, errs.ErrInvalidInput):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// -------- 3. Success --------
	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func (h *userHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	access, err := h.service.RefreshToken(c, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
	})
}
