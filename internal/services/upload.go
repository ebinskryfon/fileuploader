package services

import (
	"bytes"
	"io"
	"mime/multipart"
	"time"

	"github.com/ebinskryfon/fileuploader/internal/config"
	"github.com/ebinskryfon/fileuploader/internal/models"
	"github.com/ebinskryfon/fileuploader/internal/storage"
	"github.com/ebinskryfon/fileuploader/internal/utils"
)

type UploadService struct {
	storage    storage.StorageInterface
	validation *ValidationService
	logger     *utils.Logger
}

func NewUploadService(cfg *config.Config, storage storage.StorageInterface, logger *utils.Logger) *UploadService {
	return &UploadService{
		storage:    storage,
		validation: NewValidationService(cfg),
		logger:     logger,
	}
}

func (u *UploadService) UploadFile(fileHeader *multipart.FileHeader, userID string) (*models.UploadResponse, *models.AppError) {
	// Validate file
	if err := u.validation.ValidateFile(fileHeader); err != nil {
		u.logger.Warn("File validation failed", map[string]interface{}{
			"file_name": fileHeader.Filename,
			"file_size": fileHeader.Size,
			"user_id":   userID,
			"error":     err.Message,
		})
		return nil, err
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		u.logger.Error("Failed to open uploaded file", map[string]interface{}{
			"file_name": fileHeader.Filename,
			"user_id":   userID,
			"error":     err.Error(),
		})
		return nil, models.NewAppError(500, "Failed to process file", err)
	}
	defer file.Close()

	// Validate file content
	if validationErr := u.validation.ValidateFileContent(file, fileHeader.Header.Get("Content-Type")); validationErr != nil {
		u.logger.Warn("File content validation failed", map[string]interface{}{
			"file_name": fileHeader.Filename,
			"user_id":   userID,
			"error":     validationErr.Message,
		})
		return nil, validationErr
	}

	// Read file content for checksum calculation
	fileContent, err := io.ReadAll(file)
	if err != nil {
		u.logger.Error("Failed to read file content", map[string]interface{}{
			"file_name": fileHeader.Filename,
			"user_id":   userID,
			"error":     err.Error(),
		})
		return nil, models.NewAppError(500, "Failed to process file", err)
	}

	// Generate file ID and metadata
	fileID := utils.GenerateUUID()
	checksum := utils.CalculateChecksum(fileContent)

	metadata := models.FileMetadata{
		ID:           fileID,
		OriginalName: utils.SanitizeFileName(fileHeader.Filename),
		Size:         fileHeader.Size,
		ContentType:  fileHeader.Header.Get("Content-Type"),
		UploadTime:   time.Now().UTC(),
		URL:          "/files/" + fileID,
		Checksum:     checksum,
		UserID:       userID,
	}

	// Store file
	reader := bytes.NewReader(fileContent)
	if err := u.storage.Store(fileID, reader, metadata); err != nil {
		u.logger.Error("Failed to store file", map[string]interface{}{
			"file_id":   fileID,
			"file_name": fileHeader.Filename,
			"user_id":   userID,
			"error":     err.Error(),
		})
		return nil, models.NewAppError(500, "Failed to store file", err)
	}

	u.logger.Info("File uploaded successfully", map[string]interface{}{
		"file_id":      fileID,
		"file_name":    fileHeader.Filename,
		"file_size":    fileHeader.Size,
		"content_type": metadata.ContentType,
		"user_id":      userID,
		"checksum":     checksum,
	})

	response := &models.UploadResponse{
		ID:          fileID,
		URL:         metadata.URL,
		Size:        metadata.Size,
		ContentType: metadata.ContentType,
		UploadTime:  metadata.UploadTime,
		Checksum:    checksum,
	}

	return response, nil
}

func (u *UploadService) GetFile(fileID, userID string) (io.ReadCloser, models.FileMetadata, *models.AppError) {
	if !u.storage.Exists(fileID) {
		return nil, models.FileMetadata{}, models.ErrFileNotFound
	}

	// Get metadata first to check ownership
	metadata, err := u.storage.GetMetadata(fileID)
	if err != nil {
		u.logger.Error("Failed to get file metadata", map[string]interface{}{
			"file_id": fileID,
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, models.FileMetadata{}, models.ErrFileNotFound
	}

	// Check if user owns the file
	if metadata.UserID != userID {
		u.logger.Warn("Unauthorized file access attempt", map[string]interface{}{
			"file_id":    fileID,
			"user_id":    userID,
			"file_owner": metadata.UserID,
		})
		return nil, models.FileMetadata{}, models.ErrFileNotFound // Don't reveal file exists
	}

	file, _, err2 := u.storage.Retrieve(fileID)
	if err2 != nil {
		u.logger.Error("Failed to retrieve file", map[string]interface{}{
			"file_id": fileID,
			"user_id": userID,
			"error":   err2.Error(),
		})
		return nil, models.FileMetadata{}, models.ErrInternalServer
	}

	u.logger.Info("File retrieved successfully", map[string]interface{}{
		"file_id": fileID,
		"user_id": userID,
	})

	return file, metadata, nil
}
