package main

import (
	"context"
	"errors"
	"fmt"
	"github/Mitchxxx/Go-Ecommerce/internal/config"
	"github/Mitchxxx/Go-Ecommerce/internal/database"
	"github/Mitchxxx/Go-Ecommerce/internal/events"
	"github/Mitchxxx/Go-Ecommerce/internal/interfaces"
	"github/Mitchxxx/Go-Ecommerce/internal/logger"
	"github/Mitchxxx/Go-Ecommerce/internal/providers"
	"github/Mitchxxx/Go-Ecommerce/internal/server"
	"github/Mitchxxx/Go-Ecommerce/internal/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// @title E-Commerce API
// @version 2.0
// @description A modern e-commerce API built with GO, Gin and GORM
// @termsOfService http://swagger.io/terms/

// @contact.name Mitchel Egboko
// @contact.url http://linkedin.com/in/megboko
// @contact.email megboko@ymail.com

// @license.name Apache 2.0
// @license.url http://ww.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemas http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {

	// Setup Logger
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load database")
	}

	// Connect to Database
	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}
	defer func() {
		if err := mainDB.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}()

	ctx := context.Background()

	eventPublisher, err := events.NewEventPublisher(ctx, &cfg.AWS)
	if err != nil {
		log.Error().Err(err).Msg("failed to create even publisher")
		return
	}
	gin.SetMode(cfg.Server.GinMode)

	authService := services.NewAuthService(db, cfg, eventPublisher)
	productService := services.NewProductService(db)
	userService := services.NewUserService(db)
	cartService := services.NewCartService(db)
	orderService := services.NewOrderService(db)

	var uploadProvider interfaces.UploadProvider
	if cfg.Upload.UploadProvider == "s3" {
		uploadProvider = providers.NewS3Provider(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Upload.Path)
	}
	uploadService := services.NewUploadService(uploadProvider)

	// Launch Server
	srv := server.New(cfg,
		db,
		&log,
		authService,
		productService,
		userService,
		uploadService,
		cartService,
		orderService)

	router := srv.SetupRoutes()
	/// Http Server instance
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Goroutine to start the server

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("starting http server")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start hrrp server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
		return
	}

	log.Info().Msg("shutting down database")
}
