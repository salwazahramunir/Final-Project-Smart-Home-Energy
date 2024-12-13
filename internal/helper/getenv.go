package helper

import (
	"fmt"
	"os"
)

func GetToken(key string) (string, error) {
	token := os.Getenv(key)

	if token == "" {
		return "", fmt.Errorf("variabel %s is not set in the .env file", key)
	}

	return token, nil
}
