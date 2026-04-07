package routes

import (
	"github.com/Sidi1901/urlShortner/internal/middleware"
	"github.com/gin-gonic/gin"
)

type RouteRegistrar interface {
	RegisterPublicRoutes(r *gin.Engine)
	RegisterProtectedRoutes(rg *gin.RouterGroup)
}

func SetupRoutes(r *gin.Engine, mw *middleware.Middleware, handlers []RouteRegistrar) {

	// 1. Public routes
	for _, h := range handlers {
		h.RegisterPublicRoutes(r)
	}

	// 2. Protected routes
	api := r.Group("/api/v1")
	api.Use(mw.AuthMiddleware())

	for _, h := range handlers {
		h.RegisterProtectedRoutes(api)
	}
}
