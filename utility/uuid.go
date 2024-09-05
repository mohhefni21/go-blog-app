package utility

import (
	"fmt"

	"github.com/google/uuid"
)

func ParseUUID(uuidStr string) (uuid.UUID, error) {
	publicId, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to parse UUID: %w", err)
	}
	return publicId, nil
}
