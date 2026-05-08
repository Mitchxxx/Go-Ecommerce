package main

import (
	"github/Mitchxxx/Go-Ecommerce/internal/config"
	"github/Mitchxxx/Go-Ecommerce/internal/database"
	"github/Mitchxxx/Go-Ecommerce/internal/logger"

	"github.com/gin-gonic/gin"
)

func main() {

	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}
	defer mainDB.Close()
	gin.SetMode(cfg.Server.GinMode)

	log.Info().Msg("Starting server")
}
