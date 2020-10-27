package delivery

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
	"github.com/stretchr/stew/slice"
)

type UploadImageResponse struct {
	Message  string `json:"message"`
	ImageURL string `json:"image_url"`
}

type ImageHandler struct {
	useCase domain.IImageUsecase
}

func NewImageHandler(router *gin.Engine, imageUsecase domain.IImageUsecase) {
	handler := &ImageHandler{
		useCase: imageUsecase,
	}

	router.POST("/api/image", handler.SaveImage)
}

func (ih ImageHandler) SaveImage(c *gin.Context) {
	var response UploadImageResponse
	var email string = c.GetString("email")

	//get image
	imageFile, err := c.FormFile("image")
	if imageFile != nil && err != nil {
		cerr := cerror.NewAndPrintWithTag("UIP00", err, global.FRIENDLY_MESSAGE)
		response.Message = cerr.FriendlyMessageWithTag()
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	/*start validate image*/
	if imageFile == nil {
		cerr := cerror.NewAndPrintWithTag("UIP01", err, global.FRIENDLY_IMAGE_REQUIRED)
		response.Message = cerr.FriendlyMessageWithTag()
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var imageHeader string
	if len(imageFile.Header["Content-Type"]) > 0 {
		imageHeader = imageFile.Header["Content-Type"][0]
	}

	if !slice.Contains(global.AllowedMIME, imageHeader) {
		errorMessage := fmt.Sprintf(global.ERR_IMAGE_NOT_ALLOWED, imageHeader)
		friendlyMessage := fmt.Sprintf(global.FRIENDLY_IMAGE_NOT_ALLOWED, imageHeader)
		cerr := cerror.NewAndPrintWithTag("UIP02", errors.New(errorMessage), friendlyMessage)
		response.Message = cerr.FriendlyMessageWithTag()
		c.JSON(http.StatusBadRequest, response)
		return
	}
	/*end validate image*/

	url, err := ih.useCase.SaveImage(c, imageFile, email)
	if err != nil {
		response.Message = err.(cerror.Error).FriendlyMessageWithTag()
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.ImageURL = url
	c.JSON(http.StatusOK, response)
	return
}
