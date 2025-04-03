package api

import (
	"log/slog"
	"net/http"

	"github.com/ericbutera/project/internal/repo"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type Server struct {
	router *gin.Engine
}

func New(repo *repo.Repo) (*Server, error) {
	handlers, err := NewHandlers(repo)
	if err != nil {
		return nil, err
	}

	router := NewRouter(handlers)
	if err := router.SetTrustedProxies([]string{}); err != nil {
		return nil, err
	}

	return &Server{
		router: router,
	}, nil
}

func (s *Server) Start() error {
	return s.router.Run()
}

func NewRouter(handlers *Handlers) *gin.Engine {
	router := gin.New()
	routes(router, handlers)
	return router
}

func routes(router *gin.Engine, handlers *Handlers) {
	router.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	v1 := router.Group("/v1").
		Use(sloggin.New(slog.Default())).
		Use(gin.Recovery())
	{
		v1.POST("/devices/:id/readings", handlers.DeviceStoreReadings)
		v1.GET("/devices/:id/readings/latest", handlers.DeviceLatestReading)
		v1.GET("/devices/:id/readings/count", handlers.DeviceReadingCount)
		v1.GET("/devices/:id/readings", handlers.DeviceGetReadings)
	}
}
