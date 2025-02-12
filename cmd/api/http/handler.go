package http

import (
	"net/http"

	"github.com/iamsorryprincess/go-project-layout/internal/pkg/http/middleware"
	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
)

func NewHandler(logger log.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", ping)

	return middleware.WithHandler(
		mux,
		middleware.Recovery(logger),
		middleware.CORS,
	)
}
