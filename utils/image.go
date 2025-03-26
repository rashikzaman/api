package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveBase64Image(base64String, uploadDir string) (string, error) {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create upload directory: %w", err)
		}
	}

	parts := strings.Split(base64String, ",")
	if len(parts) != 2 {
		return "", errors.New("invalid base64 image format")
	}

	headerParts := strings.Split(parts[0], ";")
	if len(headerParts) != 2 || !strings.HasPrefix(headerParts[0], "data:image/") {
		return "", errors.New("invalid base64 image header")
	}

	imageFormat := strings.TrimPrefix(headerParts[0], "data:image/")
	if imageFormat == "" {
		return "", errors.New("could not determine image format")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %w", err)
	}

	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("image_%d.%s", timestamp, imageFormat)
	filePath := filepath.Join(uploadDir, filename)

	err = os.WriteFile(filePath, decoded, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return filePath, nil
}

func GetMimeTypeFromBase64(base64String string) (string, error) {
	// check if the string has the expected data URI format
	if !strings.HasPrefix(base64String, "data:") {
		return "", errors.New("invalid base64 format: missing data URI prefix")
	}

	// split the base64 string to get the header part
	// format usually is: data:image/jpeg;base64,/9j/4AAQSkZJRgABA...
	parts := strings.Split(base64String, ",")
	if len(parts) != 2 {
		return "", errors.New("invalid base64 format: missing comma separator")
	}

	headerParts := strings.Split(parts[0], ";")
	if len(headerParts) < 1 {
		return "", errors.New("invalid base64 header format")
	}

	mimeType := strings.TrimPrefix(headerParts[0], "data:")

	if mimeType == "" {
		return "", errors.New("could not extract MIME type")
	}

	return mimeType, nil
}
