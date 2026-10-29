package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"actualiteco/actualite/internal/core/middleware"

	"actualiteco/actualite/internal/core/logger"

	"go.uber.org/zap"
)

type HTTPServer struct {
	mux *http.ServeMux
	cfg Config
	log *logger.Logger
}

func NewHTTPServer(config Config, logger *logger.Logger) *HTTPServer {
	return &HTTPServer{
		mux: http.NewServeMux(),
		cfg: config,
		log: logger,
	}
}

func (s *HTTPServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {

		pattern := fmt.Sprintf("%s %s", route.Method, route.URL)

		middlewares := route.AdditionalMiddleware

		handlerWithMiddleware := middleware.ChainMiddleware(
			route.Handler,
			middlewares...,
		)

		s.mux.Handle(pattern, handlerWithMiddleware)
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	handler := middleware.ChainMiddleware(
		s.mux,
		middleware.PanicMiddleware(s.log),
		middleware.LogMiddleware(s.log),
		middleware.TraceMiddleware(),
		middleware.CORSMiddleware(),
	)

	server := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: handler,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Info("starting http server", zap.String("addr", s.cfg.Addr))
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("error from server: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown htttp server...")

		shutdownCtx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown http server: %w", err)
		}

		s.log.Warn("http server stopped")
	}

	return nil
}
