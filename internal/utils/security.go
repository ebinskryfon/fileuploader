package utils

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"strings"
)

func SanitizeFileName(fileName string) string {
	// Remove path components and dangerous characters
	fileName = filepath.Base(fileName)
	fileName = strings.ReplaceAll(fileName, "..", "")
	fileName = strings.ReplaceAll(fileName, "/", "")
	fileName = strings.ReplaceAll(fileName, "\\", "")
	return fileName
}

func IsAllowedPath(basePath, filePath string) bool {
	// Ensure the file path is within the base path
	cleanBase := filepath.Clean(basePath)
	cleanFile := filepath.Clean(filePath)

	rel, err := filepath.Rel(cleanBase, cleanFile)
	if err != nil {
		return false
	}

	return !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".."
}

func CalculateChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
