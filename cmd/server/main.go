package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ebinskryfon/fileuploader/config"
    "github.com/ebinskryfon/fileuploader/server"
    "github.com/ebinskryfon/fileuploader/utils"
)

func main() {
	// Initialize logger
	logger := utils.NewLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration: " + err.Error())
		os.Exit(1)
	}

	// Create and start server
	srv := server.New(cfg, logger)
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      srv.Router(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting server on port " + cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed to start: " + err.Error())
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: " + err.Error())
		os.Exit(1)
	}

	logger.Info("Server exited")
}
