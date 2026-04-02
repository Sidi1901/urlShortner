package handler

import (
	"net/http"

	"github.com/Sidi1901/urlShortner/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var validate = validator.New()

func (h *Handler) Signup(c *gin.Context) {
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
}

// func (h *Handler)Login(c *gin.Context){

// }
