package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ebinskryfon/fileuploader/models"
	"github.com/ebinskryfon/fileuploader/utils"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
	}
}

func (ls *LocalStorage) Store(fileID string, reader io.Reader, metadata models.FileMetadata) error {
	// Ensure the file path is safe
	filePath := filepath.Join(ls.basePath, fileID)
	if !utils.IsAllowedPath(ls.basePath, filePath) {
		return fmt.Errorf("invalid file path")
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Copy data to file
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	// Store metadata
	metadataPath := filePath + ".meta"
	metadataFile, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %v", err)
	}
	defer metadataFile.Close()

	encoder := json.NewEncoder(metadataFile)
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to write metadata: %v", err)
	}

	return nil
}

func (ls *LocalStorage) Retrieve(fileID string) (io.ReadCloser, models.FileMetadata, error) {
	var metadata models.FileMetadata

	filePath := filepath.Join(ls.basePath, fileID)
	if !utils.IsAllowedPath(ls.basePath, filePath) {
		return nil, metadata, fmt.Errorf("invalid file path")
	}

	// Load metadata
	metadata, err := ls.GetMetadata(fileID)
	if err != nil {
		return nil, metadata, err
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, metadata, fmt.Errorf("failed to open file: %v", err)
	}

	return file, metadata, nil
}

func (ls *LocalStorage) Delete(fileID string) error {
	filePath := filepath.Join(ls.basePath, fileID)
	if !utils.IsAllowedPath(ls.basePath, filePath) {
		return fmt.Errorf("invalid file path")
	}

	// Delete file
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	// Delete metadata
	metadataPath := filePath + ".meta"
	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete metadata: %v", err)
	}

	return nil
}

func (ls *LocalStorage) Exists(fileID string) bool {
	filePath := filepath.Join(ls.basePath, fileID)
	if !utils.IsAllowedPath(ls.basePath, filePath) {
		return false
	}

	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func (ls *LocalStorage) GetMetadata(fileID string) (models.FileMetadata, error) {
	var metadata models.FileMetadata

	filePath := filepath.Join(ls.basePath, fileID)
	if !utils.IsAllowedPath(ls.basePath, filePath) {
		return metadata, fmt.Errorf("invalid file path")
	}

	metadataPath := filePath + ".meta"
	metadataFile, err := os.Open(metadataPath)
	if err != nil {
		return metadata, fmt.Errorf("failed to open metadata file: %v", err)
	}
	defer metadataFile.Close()

	decoder := json.NewDecoder(metadataFile)
	if err := decoder.Decode(&metadata); err != nil {
		return metadata, fmt.Errorf("failed to decode metadata: %v", err)
	}

	return metadata, nil
}
