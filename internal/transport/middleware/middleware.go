package middleware

import (
	"log"

	"go.uber.org/zap"
)

type Middleware struct {
	logger *zap.Logger
}

func NewMiddleware(logger *zap.Logger) *Middleware {
	if logger == nil {
		log.Println("логгер не передан")
		logger = zap.NewNop()
	}
	return &Middleware{
		logger: logger,
	}
}
