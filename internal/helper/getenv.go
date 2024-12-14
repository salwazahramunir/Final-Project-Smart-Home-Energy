package helper

import (
	"fmt"
	"os"
)

func GetENV(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", fmt.Errorf("variabel %s is not set in the .env file", key)
	}

	return value, nil
}
