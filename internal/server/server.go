package server

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

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{ //nolint:gocritic // I need this for readability
			auth.POST("/register", s.register)
			auth.POST("/login", s.login)
			auth.POST("refresh", s.refreshToken)
			auth.POST("/logout", s.logout)
		}

		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			// User Routes
			users := protected.Group("/users")
			{
				userRoutes := users
				userRoutes.GET("/profile", s.getProfile)
				userRoutes.PUT("/profile", s.updateProfile)
			}
			// category routes
			categories := protected.Group("/categories")
			{
				categoriesRoutes := categories
				categoriesRoutes.POST("/", s.adminMiddleware(), s.createCategory)
				categoriesRoutes.PUT("/:id", s.adminMiddleware(), s.updateCategory)
				categoriesRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)

			}
			// product Routes
			products := protected.Group("/products")
			{
				productRoutes := products
				productRoutes.POST("/", s.adminMiddleware(), s.createProduct)
				productRoutes.PUT("/", s.adminMiddleware(), s.updateProduct)
				productRoutes.DELETE("/", s.adminMiddleware(), s.deleteProduct)
			}
		}

		// public routes
		api.GET("/categories", s.getCategories)
		api.GET("/products", s.getProducts)
		api.GET("/products/:id", s.getProduct)

	}

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
