package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	inner  http.Server
	logger *zap.Logger
}

func handlerMessage(msg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(msg))
	}
}

func handlerDelayed(timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		<-time.After(timeout)
	}
}

func handlerHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If we can call this at all, we're healthy
		w.Write([]byte("ok"))
	}
}

func middlewareLogRequest(handler http.HandlerFunc, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug(
			"Endpoint called",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)

		handler(w, r)
	}
}

func New(logger *zap.Logger, addr string) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/",
		middlewareLogRequest(
			handlerMessage("Hello from base"),
			logger,
		))

	mux.HandleFunc(
		"/delay",
		middlewareLogRequest(
			handlerDelayed(time.Second*2),
			logger,
		),
	)

	mux.HandleFunc(
		"/healthz",
		middlewareLogRequest(
			handlerHealthz(),
			logger,
		),
	)

	return &Server{
		inner: http.Server{
			Addr:    addr,
			Handler: mux,
		},
		logger: logger.With(zap.String("component", "server")),
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	errServer := make(chan error, 1)
	go func() {
		errServer <- s.inner.ListenAndServe()
	}()

	select {
	case err := <-errServer:
		return fmt.Errorf("failed to listen: %w", err)

	case <-ctx.Done():

		timeout := time.Second * 5

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		s.logger.Info("Shutting down gracefully...", zap.Duration("timeout", timeout))
		err := s.inner.Shutdown(ctx)

		if err != nil {
			return fmt.Errorf("failed to Shutdown(): %w", err)
		}

		s.logger.Info("Graceful shutdown complete")
		return nil
	}
}
