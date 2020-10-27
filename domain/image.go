package domain

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type Image struct {
	ImageID  string
	ImageURL string
}

type IImageRepository interface {
	SaveImage(image Image) error
}

type IImageUsecase interface {
	SaveImage(c *gin.Context, imageFile *multipart.FileHeader, email string) (string, error)
}
