package handlers

import (
	"strings"
	"time"

	"fileuploader/internal/models"
	"fileuploader/internal/services"
	"fileuploader/internal/utils"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	authService *services.AuthService
	logger      *utils.Logger
	rateLimiter map[string][]time.Time // Simple in-memory rate limiter
	rateLimit   int
}

func NewMiddleware(authService *services.AuthService, logger *utils.Logger, rateLimit int) *Middleware {
	return &Middleware{
		authService: authService,
		logger:      logger,
		rateLimiter: make(map[string][]time.Time),
		rateLimit:   rateLimit,
	}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.respondWithError(c, models.ErrUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			m.respondWithError(c, models.ErrUnauthorized)
			return
		}

		token := tokenParts[1]
		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			m.logger.Warn("Invalid token", map[string]interface{}{
				"error": err.Error(),
				"ip":    c.ClientIP(),
			})
			m.respondWithError(c, models.ErrUnauthorized)
			return
		}

		// Set user ID in context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func (m *Middleware) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		if userID == "" {
			userID = c.ClientIP() // Use IP if no user ID
		}

		now := time.Now()
		cutoff := now.Add(-time.Minute)

		// Clean old entries
		var validRequests []time.Time
		for _, requestTime := range m.rateLimiter[userID] {
			if requestTime.After(cutoff) {
				validRequests = append(validRequests, requestTime)
			}
		}

		// Check rate limit
		if len(validRequests) >= m.rateLimit {
			m.logger.Warn("Rate limit exceeded", map[string]interface{}{
				"user_id": userID,
				"ip":      c.ClientIP(),
			})
			m.respondWithError(c, models.ErrRateLimitExceeded)
			return
		}

		// Add current request
		validRequests = append(validRequests, now)
		m.rateLimiter[userID] = validRequests

		c.Next()
	}
}

func (m *Middleware) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		m.logger.Info("Request processed", map[string]interface{}{
			"method":      method,
			"path":        path,
			"status_code": statusCode,
			"duration_ms": duration.Milliseconds(),
			"ip":          c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
			"user_id":     c.GetString("userID"),
		})
	}
}

func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (m *Middleware) respondWithError(c *gin.Context, appError *models.AppError) {
	response := models.ErrorResponse{
		Error:   appError.Message,
		Code:    appError.Code,
		Message: appError.Message,
	}
	c.JSON(appError.Code, response)
	c.Abort()
}
