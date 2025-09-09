// server/server.go
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ebinskryfon/fileuploader/config"
	"github.com/ebinskryfon/fileuploader/handlers"
	"github.com/ebinskryfon/fileuploader/services"
	"github.com/ebinskryfon/fileuploader/storage"
	"github.com/ebinskryfon/fileuploader/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config  *config.Config
	logger  *utils.Logger
	storage storage.StorageInterface
	router  *gin.Engine
}

func New(cfg *config.Config, logger *utils.Logger) *Server {
	// Initialize storage
	localStorage := storage.NewLocalStorage(cfg.Upload.StoragePath)

	// Initialize services
	authService := services.NewAuthService(cfg)
	uploadService := services.NewUploadService(cfg, localStorage, logger)

	// Initialize handlers
	uploadHandler := handlers.NewUploadHandler(uploadService, logger)
	downloadHandler := handlers.NewDownloadHandler(uploadService, logger)
	healthHandler := handlers.NewHealthHandler()

	// Initialize middleware
	middleware := handlers.NewMiddleware(authService, logger, cfg.RateLimit.RequestsPerMinute)

	// Setup Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Health endpoints (no auth required)
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	// API routes with authentication and rate limiting
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	api.Use(middleware.RateLimitMiddleware())
	{
		api.POST("/upload", uploadHandler.Upload)
		api.GET("/files/:id", downloadHandler.GetFile)
	}

	// Direct file access (backward compatibility)
	files := router.Group("/files")
	files.Use(middleware.AuthMiddleware())
	files.Use(middleware.RateLimitMiddleware())
	{
		files.GET("/:id", downloadHandler.GetFile)
	}

	return &Server{
		config:  cfg,
		logger:  logger,
		storage: localStorage,
		router:  router,
	}
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

// Additional idiomatic improvements you might consider:

// 1. Add a Start method to Server for better encapsulation
func (s *Server) Start(addr string) error {
	s.logger.Info("Starting server on " + addr)
	return http.ListenAndServe(addr, s.router)
}

// 2. Add a Shutdown method to Server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")
	// Add any cleanup logic here
	return nil
}

// 3. Consider adding a Run method that handles the full server lifecycle
func (s *Server) Run(addr string) error {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return fmt.Errorf("server failed to start: %w", err)
	case <-quit:
		s.logger.Info("Shutting down server...")
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return httpServer.Shutdown(ctx)
}
