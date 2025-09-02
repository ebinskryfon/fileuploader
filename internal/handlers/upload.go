package handlers

import (
	"net/http"

	"fileuploader/internal/models"
	"fileuploader/internal/services"
	"fileuploader/internal/utils"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService *services.UploadService
	logger        *utils.Logger
}

func NewUploadHandler(uploadService *services.UploadService, logger *utils.Logger) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
		logger:        logger,
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		h.respondWithError(c, models.ErrUnauthorized)
		return
	}

	// Parse multipart form
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		h.logger.Warn("Failed to parse form file", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		h.respondWithError(c, models.NewAppError(http.StatusBadRequest, "No file provided", err))
		return
	}
	defer file.Close()

	// Upload file
	response, appError := h.uploadService.UploadFile(fileHeader, userID)
	if appError != nil {
		h.respondWithError(c, appError)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UploadHandler) respondWithError(c *gin.Context, appError *models.AppError) {
	response := models.ErrorResponse{
		Error:   appError.Message,
		Code:    appError.Code,
		Message: appError.Message,
	}
	c.JSON(appError.Code, response)
}
