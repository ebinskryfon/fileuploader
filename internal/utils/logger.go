package utils

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Logger struct {
	*log.Logger
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) logEntry(level, message string, data map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Data:      data,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		l.Printf("Failed to marshal log entry: %v", err)
		return
	}

	l.Println(string(jsonBytes))
}

func (l *Logger) Info(message string, data ...map[string]interface{}) {
	var d map[string]interface{}
	if len(data) > 0 {
		d = data[0]
	}
	l.logEntry("INFO", message, d)
}

func (l *Logger) Error(message string, data ...map[string]interface{}) {
	var d map[string]interface{}
	if len(data) > 0 {
		d = data[0]
	}
	l.logEntry("ERROR", message, d)
}

func (l *Logger) Warn(message string, data ...map[string]interface{}) {
	var d map[string]interface{}
	if len(data) > 0 {
		d = data[0]
	}
	l.logEntry("WARN", message, d)
}

func (l *Logger) Debug(message string, data ...map[string]interface{}) {
	var d map[string]interface{}
	if len(data) > 0 {
		d = data[0]
	}
	l.logEntry("DEBUG", message, d)
}
