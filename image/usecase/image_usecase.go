package usecase

import (
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kennygrant/sanitize"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
)

type ImageUsecase struct {
	imageRepo domain.IImageRepository
}

func NewImageUsecase(imageRepository domain.IImageRepository) domain.IImageUsecase {
	return &ImageUsecase{
		imageRepo: imageRepository,
	}
}

func (iu ImageUsecase) SaveImage(c *gin.Context,
	imageFile *multipart.FileHeader, email string) (string, error) {
	//create filename
	fileExtension := filepath.Ext(imageFile.Filename)
	timestamp := time.Now().Format("20060102150405")
	filename := sanitize.BaseName(email+"_"+timestamp) + fileExtension
	path := "upload/images/" + filename

	//upload file
	err := c.SaveUploadedFile(imageFile, path)
	if err != nil {
		cerror := cerror.NewAndPrintWithTag("UIP00", err, global.FRIENDLY_MESSAGE)
		return "", cerror
	}

	//save image data to db
	var image domain.Image
	image.ImageID = uuid.New().String()
	image.ImageURL = "/" + path
	err = iu.imageRepo.SaveImage(image)
	if err != nil {
		return "", err
	}

	return image.ImageURL, nil
}
