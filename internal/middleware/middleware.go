package middleware

import (
	"github.com/Sidi1901/urlShortner/internal/config"
)

type Middleware struct {
	Cfg *config.Config
}

func NewMiddleware(cfg *config.Config) *Middleware {
	return &Middleware{
		Cfg: cfg,
	}
}
