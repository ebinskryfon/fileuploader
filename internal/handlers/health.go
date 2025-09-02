package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": c.Request.Header.Get("Date"),
		"service":   "fileuploader",
	}
	c.JSON(http.StatusOK, response)
}

func (h *HealthHandler) Ready(c *gin.Context) {
	// Add actual readiness checks here (database, storage, etc.)
	response := map[string]interface{}{
		"status": "ready",
		"checks": map[string]string{
			"storage": "ok",
		},
	}
	c.JSON(http.StatusOK, response)
}
