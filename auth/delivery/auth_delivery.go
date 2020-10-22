package delivery

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
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
	Message []string `json:"message"`
	Token   string   `json:"token"`
}

type SignUpRequest struct {
	Fullname        string `json:"full_name" form:"bambang_sadikin" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Passowrd        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Passowrd"`
}

type SignUpResponse struct {
	Message []string        `json:"message"`
	Account *domain.Account `json:"account"`
	Profile *domain.Profile `json:"profile"`
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

	response.Token = token
	c.JSON(http.StatusOK, response)
	return
}

func (ah AuthHandler) SignUp(c *gin.Context) {
	var (
		request  SignUpRequest
		response SignUpResponse
	)

	err := c.ShouldBind(&request)
	if err != nil {
		cerr := cerror.NewAndPrintWithTag("ASU00", err, global.FRIENDLY_MESSAGE)

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

				case "eqfield":
					msg := fmt.Sprintf(global.ERR_DIFFERENT_FORMATTER, jsonField, "password")
					response.Message = append(response.Message, msg)
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
	account.Password = request.Passowrd
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
