package helpers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var defaultAllowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".jfif": true,
}

const defaultMaxSizeMB = 5

type SaveFileOptions struct {
	AllowedExts []string
	MaxSizeMB   int
	CustomName  string
}

func SaveUploadedFileWithOpts(c *gin.Context, fieldName string, uploadDir string, opts *SaveFileOptions) (string, error) {
	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	allowedExts := defaultAllowedExts
	maxSize := int64(defaultMaxSizeMB) * 1024 * 1024

	if opts != nil {
		if len(opts.AllowedExts) > 0 {
			allowedExts = make(map[string]bool)
			for _, ext := range opts.AllowedExts {
				allowedExts[strings.ToLower(ext)] = true
			}
		}
		if opts.MaxSizeMB > 0 {
			maxSize = int64(opts.MaxSizeMB) * 1024 * 1024
		}
	}

	if header.Size > maxSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %dMB", maxSize/1024/1024)
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExts[ext] {
		return "", fmt.Errorf("file type %s is not allowed", ext)
	}

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	var fileName string
	if opts != nil && opts.CustomName != "" {
		fileName = opts.CustomName + ext
	} else {
		randomStr, err := GenerateRandomString(4)
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
		fileName = fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), randomStr, ext)
	}

	filePath := filepath.Join(uploadDir, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	normalizedDir := strings.ReplaceAll(filepath.ToSlash(uploadDir), "\\", "/")
	return fmt.Sprintf("%s/%s", normalizedDir, fileName), nil
}

func SaveUploadedFile(c *gin.Context, fieldName string, uploadDir string) (string, error) {
	return SaveUploadedFileWithOpts(c, fieldName, uploadDir, nil)
}

func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filePath)
}

func GetFileURL(path string) string {
	if path == "" {
		return ""
	}
	baseURL := GetEnv("APP_URL", "http://localhost:8080")
	normalizedPath := strings.ReplaceAll(path, "\\", "/")
	return fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), normalizedPath)
}
