package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhilashdk2016/golang-ecommerce/internal/config"
	"github.com/abhilashdk2016/golang-ecommerce/internal/database"
	"github.com/abhilashdk2016/golang-ecommerce/internal/logger"
	"github.com/abhilashdk2016/golang-ecommerce/internal/providers"
	"github.com/abhilashdk2016/golang-ecommerce/internal/server"
	"github.com/abhilashdk2016/golang-ecommerce/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}
	defer func() {
		err := mainDB.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown db")
		}
	}()

	gin.SetMode(cfg.Server.GinMode)

	authService := services.NewAuthService(db, cfg)
	productService := services.NewProductService(db)
	userService := services.NewUserService(db)
	uploadService := services.NewUploadService(providers.NewLocalUploadProvider(cfg.Upload.Path))

	srv := server.New(cfg, db, &log, authService, productService, userService, uploadService)
	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("starting http server...")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
	}

	log.Info().Msg("shutting down database...")
}
