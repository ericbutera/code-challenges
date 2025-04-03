package server

// TODO:
// openapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ericbutera/appointments/internal/repo"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type Server struct {
	router *gin.Engine
}

func New(repo *repo.Repo, location *time.Location) (*Server, error) {
	handlers, err := NewHandlers(repo, location) // TODO: Dep Injection
	if err != nil {
		return nil, err
	}
	router := NewRouter(handlers)

	// TODO: configurable middleware
	// logging
	// otel (metrics, traces)
	// error recording
	router.Use(gin.Recovery())
	router.Use(sloggin.New(slog.Default())) // TODO: Dep Injection

	return &Server{
		router: router,
	}, nil
}

func (s *Server) Start() error {
	return s.router.Run()
}

func NewRouter(handlers *Handlers) *gin.Engine {
	//gin.SetMode(gin.TestMode) // TODO: make configurable
	router := gin.New()
	routes(router, handlers)
	return router
}

func routes(router *gin.Engine, handlers *Handlers) {
	router.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	router.GET("/availability", handlers.ListAvailability)                   // TODO: /trainers/:trainer_id/availability
	router.POST("/appointments", handlers.CreateAppointment)                 // TODO: /trainers/:trainer_id/appointments
	router.GET("/appointments", handlers.ListAppointments)                   // TODO: /trainers/:trainer_id/appointments
	router.GET("/clients/:userID/completion_status", handlers.WorkoutStatus) // TODO: /trainers/:trainer_id/appointments
	// /clients/:id/completion_status -> Timeseries[Point]

	// TODO: version api
	// v1 := router.Group("/v1") {}
}
