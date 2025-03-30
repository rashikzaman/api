package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

func DecodeBase64Image(base64String string) ([]byte, error) {
	parts := strings.Split(base64String, ",")
	if len(parts) != 2 {
		return nil, errors.New("invalid base64 image format")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %w", err)
	}

	return decoded, nil
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
