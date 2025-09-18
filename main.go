package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/manuel/make-it-rain/config"
	"github.com/manuel/make-it-rain/db"
	"github.com/manuel/make-it-rain/routes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("No .env file found, using environment variables")
	}

	if err := config.LoadConfig("."); err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	setupLogger()

	log.Info().
		Str("app", config.Cfg.App.Name).
		Str("version", config.Cfg.App.Version).
		Str("env", config.Cfg.Server.Environment).
		Msg("Starting application")

	connStr := config.Cfg.Database.GetConnectionString()
	if err := db.InitDB(connStr); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer db.CloseDB()

	if config.Cfg.Server.Environment != "test" {
		if err := db.RunMigrations(connStr); err != nil {
			log.Fatal().Err(err).Msg("Failed to run migrations")
		}
	}

	if config.Cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	routes.SetupRoutes(router)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  config.Cfg.Server.ReadTimeout,
		WriteTimeout: config.Cfg.Server.WriteTimeout,
	}

	go func() {
		log.Info().Str("port", config.Cfg.Server.Port).Msg("Server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server shutdown completed")
}

func setupLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	if config.Cfg.Server.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	switch config.Cfg.App.LogLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}