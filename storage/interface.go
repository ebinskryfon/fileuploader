package storage

import (
	"github.com/ebinskryfon/fileuploader/models"
	"io"
)

type StorageInterface interface {
	Store(fileID string, reader io.Reader, metadata models.FileMetadata) error
	Retrieve(fileID string) (io.ReadCloser, models.FileMetadata, error)
	Delete(fileID string) error
	Exists(fileID string) bool
	GetMetadata(fileID string) (models.FileMetadata, error)
}
