// Package rest defines the HTTP server and request handling logic for the API.
package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/faeelol/companies-store/internal/app/database/repositories"
	"github.com/faeelol/companies-store/internal/app/kafka"
	"github.com/faeelol/companies-store/internal/app/logic/companies"
	"github.com/faeelol/companies-store/internal/app/rest/handlers"
	"github.com/faeelol/companies-store/internal/app/rest/jwt"
	"github.com/faeelol/companies-store/internal/app/rest/middlewares"
)

type Server struct {
	cfg        *Config
	httpServer *http.Server
}

func NewServer(cfg *Config, logger *logrus.Logger, db *sqlx.DB, producer *kafka.Producer) *Server {
	jwtParser := jwt.NewJWTParser(*cfg.JWT)

	companiesController := companies.NewCompaniesController(db, repositories.NewCompanyRepository(), producer)

	routes := createRoutingTable(logger, jwtParser, companiesController)

	httpServer := &http.Server{
		Addr:              cfg.Addr,
		WriteTimeout:      cfg.WriteTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		Handler:           routes,
	}

	return &Server{cfg: cfg, httpServer: httpServer}
}

func createRoutingTable(
	logger *logrus.Logger,
	jwtParser *jwt.Parser,
	companiesController *companies.Controller,
) chi.Router {
	r := chi.NewRouter()
	r.Use(middlewares.NewErrorHandlerMiddleware(logger))
	r.Use(middlewares.NewLoggingMiddleware(logger))

	authAdminMiddleware := middlewares.NewJWTMiddleware(jwtParser, []string{"admin"})

	r.Route("/api/companies_repo/v1", func(r chi.Router) {
		r.Method(http.MethodGet, "/companies", handlers.NewGetCompaniesHandler(companiesController))
		r.Group(func(r chi.Router) {
			r.Use(authAdminMiddleware.VerifyToken)
			r.Method(http.MethodPost, "/companies", handlers.NewCreateCompaniesHandler(companiesController))
			r.Method(http.MethodDelete, "/companies", handlers.NewDeleteCompaniesHandler(companiesController))
			r.Method(http.MethodPatch, "/companies", handlers.NewPatchCompaniesHandler(companiesController))
		})
	})
	return r
}

func (s Server) Start(ctx context.Context, logger *logrus.Logger) error {
	errChan := make(chan error, 1)

	go func() {
		logger.WithField("addr", s.cfg.Addr).Info("starting HTTP server")
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("server error: %w", err)
		}
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutting down HTTP service...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("error shutting down server: %w", err)
		}
		return nil

	case err := <-errChan:
		return err
	}
}
