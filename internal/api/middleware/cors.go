package middleware

import (
	"github.com/rs/cors"
)

func NewCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "OPTIONS"},
	})
}
