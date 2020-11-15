package delivery

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
)

type GetProfileResponse struct {
	Message  string `json:"message"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name" binding:"required"`
}

type UpdateProfileResponse struct {
	Message []string `json:"message"`
}

type ProfileHandler struct {
	useCase domain.IProfileUsecase
}

func NewProfileHandler(router *gin.Engine,
	profileUsecase domain.IProfileUsecase) {
	handler := &ProfileHandler{
		useCase: profileUsecase,
	}

	router.GET("/api/profile", handler.GetProfile)
	router.POST("/api/profile/update", handler.UpdateProfile)
}

func (ph ProfileHandler) GetProfile(c *gin.Context) {
	var (
		accountID string = c.GetString("account_id")
		response  GetProfileResponse
	)

	profileInput := domain.Profile{AccountID: accountID}
	profile, err := ph.useCase.GetProfile(profileInput)
	if err != nil {
		var cerr cerror.Error
		cerr, ok := err.(cerror.Error)
		if !ok {
			cerr = cerror.NewAndPrintWithTag("GPH00", err, global.FRIENDLY_MESSAGE)
		}

		response.Message = cerr.FriendlyMessageWithTag()
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.FullName = profile.FullName
	response.Email = profile.Account.Email
	c.JSON(http.StatusOK, response)
	return
}

func (ph ProfileHandler) UpdateProfile(c *gin.Context) {
	var (
		request   UpdateProfileRequest
		response  UpdateProfileResponse
		accountID string = c.GetString("account_id")
	)

	//validate input
	err := c.ShouldBind(&request)
	if err != nil {
		cerr := cerror.NewAndPrintWithTag("UPH00", err, global.FRIENDLY_MESSAGE)

		/*start validation*/
		valError := err.(validator.ValidationErrors)
		if valError != nil {
			for _, elem := range valError {
				fieldName := elem.Field()
				field, _ := reflect.TypeOf(&request).Elem().FieldByName(fieldName)
				jsonField, _ := field.Tag.Lookup("json")

				switch elem.Tag() {
				case "required":
					msg := fmt.Sprintf(global.ERR_REQUIRED_FORMATTER, jsonField)
					response.Message = append(response.Message, msg)
					break
				}

				c.JSON(http.StatusBadRequest, response)
				return
			}
			/*end validation*/
			response.Message = []string{cerr.FriendlyMessageWithTag()}

			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	//populate input
	var profile domain.Profile
	profile.AccountID = accountID
	profile.FullName = request.FullName

	//update profile
	err = ph.useCase.UpdateProfile(profile)
	if err != nil {
		cerr, ok := err.(cerror.Error)
		if !ok {
			cerr = cerror.NewAndPrintWithTag("UPH01", err, global.FRIENDLY_MESSAGE)
		}

		msg := cerr.FriendlyMessageWithTag()
		response.Message = append(response.Message, msg)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusNoContent, response)
	return

}
