package server

import "C"
import (
	"github/Mitchxxx/Go-Ecommerce/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
	logger *zerolog.Logger
}

func New(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger) *Server {
	return &Server{
		config: cfg,
		db:     db,
		logger: logger,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Add Middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	// Add Routes
	router.GET("/health", s.healthCheck)

	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Method", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
