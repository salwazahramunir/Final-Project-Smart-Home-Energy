package helper

import (
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

func InitCloudinary() (*cloudinary.Cloudinary, error) {
	CLOUDINARY_NAME := os.Getenv("CLOUDINARY_NAME")
	CLOUDINARY_API_KEY := os.Getenv("CLOUDINARY_API_KEY")
	CLOUDINARY_API_SECRET := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(CLOUDINARY_NAME, CLOUDINARY_API_KEY, CLOUDINARY_API_SECRET)

	if err != nil {
		fmt.Printf("Error initializing Cloudinary: %v", err)
		return nil, err
	}

	return cld, nil
}
