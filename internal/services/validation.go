package services

import (
	"mime/multipart"
	"net/http"
	"strings"

	"fileuploader/internal/config"
	"fileuploader/internal/models"
)

type ValidationService struct {
	maxFileSize  int64
	allowedTypes map[string]bool
}

func NewValidationService(cfg *config.Config) *ValidationService {
	allowedTypes := make(map[string]bool)
	for _, t := range cfg.Upload.AllowedTypes {
		allowedTypes[t] = true
	}

	return &ValidationService{
		maxFileSize:  cfg.Upload.MaxFileSize,
		allowedTypes: allowedTypes,
	}
}

func (v *ValidationService) ValidateFile(header *multipart.FileHeader) *models.AppError {
	// Check file size
	if header.Size > v.maxFileSize {
		return models.ErrFileTooLarge
	}

	// Check content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		// Try to detect from filename
		contentType = v.detectContentTypeFromFilename(header.Filename)
	}

	if !v.allowedTypes[contentType] {
		return models.ErrInvalidFileType
	}

	return nil
}

func (v *ValidationService) ValidateFileContent(file multipart.File, contentType string) *models.AppError {
	// Read first 512 bytes to detect actual content type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && n == 0 {
		return models.NewAppError(http.StatusBadRequest, "Cannot read file", err)
	}

	// Reset file position
	file.Seek(0, 0)

	// Detect content type
	detectedType := http.DetectContentType(buffer[:n])

	// Check if detected type matches or is compatible with declared type
	if !v.isCompatibleContentType(contentType, detectedType) {
		return models.ErrInvalidFileType
	}

	return nil
}

func (v *ValidationService) detectContentTypeFromFilename(filename string) string {
	ext := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(ext, ".jpg") || strings.HasSuffix(ext, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(ext, ".png"):
		return "image/png"
	case strings.HasSuffix(ext, ".pdf"):
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

func (v *ValidationService) isCompatibleContentType(declared, detected string) bool {
	// Exact match
	if declared == detected {
		return true
	}

	// Handle common variations
	compatibilityMap := map[string][]string{
		"image/jpeg": {"image/jpg"},
		"image/jpg":  {"image/jpeg"},
	}

	if compatible, exists := compatibilityMap[declared]; exists {
		for _, c := range compatible {
			if c == detected {
				return true
			}
		}
	}

	return false
}
