package handler

import (
	"github.com/gin-gonic/gin"
)

type webHandler struct {
}

func NewWebHandler() *webHandler {
	return &webHandler{}
}

func (h *webHandler) RegisterPublicRoutes(r *gin.Engine) {
	r.GET("/", h.Home)
	r.GET("/login", h.ShowLoginPage)
	r.GET("/signup", h.ShowRegisterPage)
	r.GET("/user", h.ShowUserProfile)
}

func (h *webHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	// rg.GET("/user/profile", h.ShowUserProfile)
}

func (h *webHandler) Home(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func (h *webHandler) ShowLoginPage(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func (h *webHandler) ShowRegisterPage(c *gin.Context) {
	c.HTML(200, "signup.html", nil)
}

func (h *webHandler) ShowUserProfile(c *gin.Context) {
	c.HTML(200, "profile.html", nil)
}
