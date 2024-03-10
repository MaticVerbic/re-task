package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"retask/api/handler"
	"retask/config"
)

// Server is an object representing a server instance, required dependencies can be added here. s
type Server struct {
	Port   int
	Server *http.Server
}

// New creates a new instance of a Mux.
// handlers could also be a client/consumer interface pattern.
func New(conf *config.Config, handlers *handler.Handler) (*Server, error) {
	mux := &Server{
		Port: conf.ServerPort,
		Server: &http.Server{
			Addr: fmt.Sprintf(":%d", conf.ServerPort),
		},
	}

	// Create a new router
	router := chi.NewRouter()

	// Add middlewares
	router.Use(requestLogging)   // This logs the request's context data.
	router.Use(CORSMiddleware()) // This adds cors

	// This ensures that if there's fatal error during running of an endpoint the server will recover instead of shut down.
	router.Use(middleware.Recoverer)

	// timeout on the request using context.Done()
	router.Use(middleware.Timeout(conf.HttpTimeout))

	// define the endpoints
	// update-package-size is a POST request
	router.Post("/update-package-sizes", handlers.UpdatePackageSizes)
	// calculate-best-packages is a POST request
	router.Post("/calculate-best-packages", handlers.CalculateBestPackages)
	// ping
	router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		if _, err := writer.Write([]byte("pong")); err != nil {
			logrus.WithField("error", err).Error("failed to write ping response")
		}
	})

	// Assign the router to the mux.
	mux.Server.Handler = router
	return mux, nil
}

// ListenAndServe create a new server and runs it in a separate go routine.
// On failure, it logs fatally and shuts down the service.
func (m *Server) ListenAndServe(logger *logrus.Entry) {
	go func() {
		if err := m.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	logger.Infof("server listening on port: %d", m.Port)
}

// Shutdown gracefully shuts the server down.
func (m *Server) Shutdown() error {
	return m.Server.Shutdown(context.Background())
}
