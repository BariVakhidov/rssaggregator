package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
)

type Server struct {
	log    *slog.Logger
	server *http.Server
}

func New(log *slog.Logger, port int, router http.Handler) *Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return &Server{
		log:    log,
		server: server,
	}
}

func (s *Server) MustRun() {
	if err := s.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (s *Server) Run() error {
	const op = "delivery.httpserver.Run"
	log := s.log.With(slog.String("op", op), slog.String("port", s.server.Addr))

	log.Info("starting http server")

	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) {
	const op = "delivery.httpserver.Stop"
	log := s.log.With(slog.String("op", op), slog.String("port", s.server.Addr))
	log.Info("stopping http server")

	if err := s.server.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown http server", sl.Err(err))
	}
}
