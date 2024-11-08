package util

import (
	"github.com/google/uuid"
)

func UUID2bytes(uuidStr string) ([]byte, error) {
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}
	return parsedUUID[:], nil
}
