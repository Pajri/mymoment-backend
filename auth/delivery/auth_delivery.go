package delivery

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
)

// #region type helper
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Message      []string `json:"message"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}

type SignUpRequest struct {
	Fullname        string `form:"full_name" binding:"required"`
	Email           string `form:"email" binding:"required"`
	Password        string `form:"password" binding:"required,min=10"`
	ConfirmPassword string `form:"confirm_password" binding:"required,eqfield=Password"`
}

type SignUpResponse struct {
	Message []string        `json:"message"`
	Account *domain.Account `json:"account"`
	Profile *domain.Profile `json:"profile"`
}

type VerifyResponse struct {
	Message string `json:"message"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordResponse struct {
	Message []string `json:"message"`
}

type ChangePasswordRequest struct {
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
}

type ChangePasswordResponse struct {
	Message []string `json:"message"`
}

// #endregion

type AuthHandler struct {
	useCase domain.IAuthUsecase
}

func NewAuthHandler(router *gin.Engine, authUsecase domain.IAuthUsecase) {
	handler := &AuthHandler{
		useCase: authUsecase,
	}

	router.POST("/api/auth/login", handler.Login)
	router.POST("/api/auth/signup", handler.SignUp)
	router.GET("/api/auth/verify_email", handler.VerifyEmail)
	router.POST("/api/auth/reset_password/", handler.ResetPassword)
	router.POST("/api/auth/change_password", handler.ChangePassword)
}

func (ah AuthHandler) Login(c *gin.Context) {
	var (
		request  LoginRequest
		response LoginResponse
	)

	err := c.ShouldBind(&request)
	if err != nil {
		cerr := cerror.NewAndPrintWithTag("ALG", err, global.FRIENDLY_MESSAGE)

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
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}
		/*end validation*/

		response.Message = []string{cerr.FriendlyMessageWithTag()}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	//populate
	var account domain.Account
	account.Email = request.Email
	account.Password = request.Password

	token, err := ah.useCase.Login(account)
	if err != nil {
		response := LoginResponse{
			Message: []string{err.(cerror.Error).FriendlyMessageWithTag()},
		}

		httpStatus := http.StatusInternalServerError
		if err.(cerror.Error).Err == sql.ErrNoRows {
			httpStatus = http.StatusNotFound
			response.Message = []string{global.FRIENDLY_INVALID_USNME_PASSWORD}
		} else if err.(cerror.Error).Type == cerror.TYPE_UNAUTHORIZED {
			httpStatus = http.StatusUnauthorized
		}

		c.JSON(httpStatus, response)
		return
	}

	response.AccessToken = token.AccessToken
	response.RefreshToken = token.RefreshToken
	c.JSON(http.StatusOK, response)
	return
}

func (ah AuthHandler) SignUp(c *gin.Context) {
	var (
		request  SignUpRequest
		response SignUpResponse
	)

	err := c.ShouldBindWith(&request, binding.Form)
	if err != nil {
		cerr := cerror.NewAndPrintWithTag("ASU00", err, global.FRIENDLY_MESSAGE)

		/*start validation*/
		valError := err.(validator.ValidationErrors)
		if valError != nil {
			for _, elem := range valError {
				fieldName := elem.Field()
				field, _ := reflect.TypeOf(&request).Elem().FieldByName(fieldName)
				jsonField, _ := field.Tag.Lookup("form")

				switch elem.Tag() {
				case "required":
					msg := fmt.Sprintf(global.ERR_REQUIRED_FORMATTER, jsonField)
					response.Message = append(response.Message, msg)
					break

				case "eqfield":
					msg := fmt.Sprintf(global.ERR_DIFFERENT_FORMATTER, jsonField, "password")
					response.Message = append(response.Message, msg)
					break

				case "min":
					msg := fmt.Sprintf(global.ERR_MIN_CHAR, jsonField, elem.Param())
					response.Message = append(response.Message, msg)
					fmt.Println(msg)
					break
				}

			}

			c.JSON(http.StatusBadRequest, response)
			return
		}
		/*end validation*/
		response.Message = []string{cerr.FriendlyMessageWithTag()}

		c.JSON(http.StatusBadRequest, response)
		return
	}

	//populate request based on domain
	var account domain.Account
	account.Password = request.Password
	account.Email = request.Email

	var profile domain.Profile
	profile.FullName = request.Fullname

	//create account
	createdAccount, createdProfile, err := ah.useCase.SignUp(account, profile)
	if err != nil {
		err.(cerror.Error).PrintErrorWithTag()
		response.Message = []string{err.(cerror.Error).FriendlyMessageWithTag()}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Account = createdAccount
	response.Profile = createdProfile
	c.JSON(http.StatusCreated, response)
}

func (ah AuthHandler) VerifyEmail(c *gin.Context) {
	var (
		emailToken string
		response   VerifyResponse
	)

	query := c.Request.URL.Query()
	if len(query) > 0 && query["token"] != nil && len(query["token"]) > 0 {
		emailToken = query["token"][0]
		err := ah.useCase.VerifyEmail(emailToken)
		if err != nil {
			_, ok := err.(cerror.Error)
			if ok {
				cerr := err.(cerror.Error)
				response.Message = cerr.FriendlyMessage
			}

			cerr := cerror.NewAndPrintWithTag("VEH00", err, global.FRIENDLY_MESSAGE)
			response.Message = cerr.FriendlyMessage
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		c.JSON(http.StatusOK, response)
		return
	}

	response.Message = global.FRIENDLY_TOKEN_REQUIRED
	c.JSON(http.StatusBadRequest, response)
}

func (ah AuthHandler) ResetPassword(c *gin.Context) {
	var (
		request  ResetPasswordRequest
		response ResetPasswordResponse
	)

	err := c.ShouldBind(&request)
	if err != nil {
		cerr := cerror.NewAndPrintWithTag("RPH00", err, global.FRIENDLY_MESSAGE)

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
				case "email":
					msg := global.FRIENDLY_INVALID_EMAIL_FORMAT
					response.Message = append(response.Message, msg)
					break
				}
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}
		/*end validation*/

		response.Message = []string{cerr.FriendlyMessageWithTag()}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	err = ah.useCase.ResetPassword(request.Email)
	if err != nil {
		cerr, ok := err.(cerror.Error)
		if ok {
			if cerr.Type != cerror.TYPE_NOT_FOUND {
				response.Message = append(response.Message, cerr.FriendlyMessageWithTag())
				c.JSON(http.StatusInternalServerError, response)
				return
			}
		}
	}

	c.JSON(http.StatusOK, response)
	return
}

func (ah AuthHandler) ChangePassword(c *gin.Context) {
	var (
		resetPasswordToken string
		request            ChangePasswordRequest
		response           ChangePasswordResponse
	)

	err := c.ShouldBind(&request)
	if err != nil {
		cerr := cerror.NewAndPrintWithTag("CPH00", err, global.FRIENDLY_MESSAGE)

		/*start form validation*/
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

				case "eqfield":
					msg := fmt.Sprintf(global.ERR_DIFFERENT_FORMATTER, jsonField, "password")
					response.Message = append(response.Message, msg)
					break
				}
			}

			c.JSON(http.StatusBadRequest, response)
			return
		}
		/*end form validation*/
		response.Message = []string{cerr.FriendlyMessageWithTag()}

		c.JSON(http.StatusBadRequest, response)
		return
	}

	query := c.Request.URL.Query()
	if len(query) > 0 && query["token"] != nil && len(query["token"]) > 0 { //token validation
		resetPasswordToken = query["token"][0]
		err := ah.useCase.ChangePassword(resetPasswordToken, request.Password)
		if err != nil {
			cerr, ok := err.(cerror.Error)
			if ok {
				response.Message = []string{cerr.FriendlyMessageWithTag()}
			} else {
				cerr := cerror.NewAndPrintWithTag("CPH01", err, global.FRIENDLY_MESSAGE)
				response.Message = []string{cerr.FriendlyMessageWithTag()}
			}

			c.JSON(http.StatusInternalServerError, response)
			return
		}

		c.JSON(http.StatusOK, response)
		return

	}

	response.Message = []string{global.FRIENDLY_TOKEN_REQUIRED}
	c.JSON(http.StatusBadRequest, response)
	return
}
