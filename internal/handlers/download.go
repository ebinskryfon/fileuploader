package handlers

import (
	"net/http"
	"fmt"
	"github.com/ebinskryfon/fileuploader/internal/models"
	"github.com/ebinskryfon/fileuploader/internal/services"
	"github.com/ebinskryfon/fileuploader/internal/utils"
	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	uploadService *services.UploadService
	logger        *utils.Logger
}

func NewDownloadHandler(uploadService *services.UploadService, logger *utils.Logger) *DownloadHandler {
	return &DownloadHandler{
		uploadService: uploadService,
		logger:        logger,
	}
}

func (h *DownloadHandler) GetFile(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		h.respondWithError(c, models.ErrUnauthorized)
		return
	}

	fileID := c.Param("id")
	if fileID == "" {
		h.respondWithError(c, models.NewAppError(http.StatusBadRequest, "File ID required", nil))
		return
	}

	file, metadata, appError := h.uploadService.GetFile(fileID, userID)
	if appError != nil {
		h.respondWithError(c, appError)
		return
	}
	defer file.Close()

	// Check Accept header to determine response format
	acceptHeader := c.GetHeader("Accept")
	if acceptHeader == "application/json" {
		// Return metadata only
		c.JSON(http.StatusOK, metadata)
		return
	}

	// Return file content
	c.Header("Content-Type", metadata.ContentType)
	c.Header("Content-Disposition", "attachment; filename=\""+metadata.OriginalName+"\"")
	c.Header("Content-Length", fmt.Sprintf("%d", metadata.Size))

	// Stream file content
	c.DataFromReader(http.StatusOK, metadata.Size, metadata.ContentType, file, nil)
}

func (h *DownloadHandler) respondWithError(c *gin.Context, appError *models.AppError) {
	response := models.ErrorResponse{
		Error:   appError.Message,
		Code:    appError.Code,
		Message: appError.Message,
	}
	c.JSON(appError.Code, response)
}
