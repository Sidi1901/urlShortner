package middleware

import (
	"github.com/Sidi1901/urlShortner/internal/config"
)

var (
	ContextUserIDKey = "user_id"
)

type Middleware struct {
	Cfg *config.Config
}

func NewMiddleware(cfg *config.Config) *Middleware {
	return &Middleware{
		Cfg: cfg,
	}
}
