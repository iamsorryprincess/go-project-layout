package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/iamsorryprincess/go-project-layout/internal/pkg/log"
)

type Server struct {
	logger log.Logger
	config Config
	server *http.Server
}

func NewServer(logger log.Logger, config Config, handler http.Handler) *Server {
	return &Server{
		logger: logger,
		config: config,
		server: &http.Server{
			Handler: handler,
			Addr:    fmt.Sprintf(":%d", config.Port),

			ReadTimeout:       config.ReadTimeout,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
			WriteTimeout:      config.WriteTimeout,
			IdleTimeout:       config.IdleTimeout,

			MaxHeaderBytes: config.MaxHeaderBytes,
		},
	}
}

func (s *Server) Start() {
	go func() {
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error().Err(err).Msg("http server listen error")
		}
	}()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("http server shutdown error")
	}
}
