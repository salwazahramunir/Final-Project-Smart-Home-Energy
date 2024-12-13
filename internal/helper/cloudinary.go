package helper

import (
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
)

func InitCloudinary() (*cloudinary.Cloudinary, error) {
	CLOUDINARY_NAME, _ := GetToken("CLOUDINARY_NAME")
	CLOUDINARY_API_KEY, _ := GetToken("CLOUDINARY_API_KEY")
	CLOUDINARY_API_SECRET, _ := GetToken("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(CLOUDINARY_NAME, CLOUDINARY_API_KEY, CLOUDINARY_API_SECRET)

	if err != nil {
		fmt.Printf("Error initializing Cloudinary: %v", err)
		return nil, err
	}

	return cld, nil
}
