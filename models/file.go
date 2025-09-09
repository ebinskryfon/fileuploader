package models

import (
	"time"
)

type FileMetadata struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	UploadTime   time.Time `json:"upload_time"`
	URL          string    `json:"url"`
	Checksum     string    `json:"checksum"`
	UserID       string    `json:"user_id"`
}

type UploadResponse struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	UploadTime  time.Time `json:"upload_time"`
	Checksum    string    `json:"checksum"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
