package middleware

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
	"github.com/pajri/personal-backend/helper"
	"github.com/stretchr/stew/slice"
)

type AuthResponse struct {
	ErrorType string `json:"error_type"`
	Message   string `json:"error_message"`
}

func handleAuth(c *gin.Context, useCase domain.IAuthUsecase) bool {
	if !global.IsEnvDevelopment() && !slice.Contains(excludedFromAuth, c.FullPath()) {
		var accountID, email string

		authArr := c.Request.Header["Authorization"]
		if len(authArr) > 0 {
			token := authArr[0]

			/*start parse jwt*/
			jwtHelper := helper.JWTHelper{}
			parsedToken, err := jwtHelper.ParseToken(token)
			if err != nil {
				cerr, ok := err.(cerror.Error)
				if !ok {
					cerr = cerror.NewAndPrintWithTag("AUM00", err, global.FRIENDLY_MESSAGE)
				}

				if cerr.Type != cerror.TYPE_EXPIRED {
					resp := AuthResponse{
						ErrorType: "internal_server_error",
						Message:   cerr.FriendlyMessageWithTag(),
					}
					c.JSON(http.StatusInternalServerError, resp)
					return false
				} else {
					resp := AuthResponse{ErrorType: "token_expired"}
					c.JSON(http.StatusUnauthorized, resp)
					return false
				}
			}

			if parsedToken == nil {
				cerr := cerror.NewAndPrintWithTag("AUM01", errors.New("parsed token is nil"), global.FRIENDLY_MESSAGE)
				resp := AuthResponse{
					ErrorType: "internal_server_error",
					Message:   cerr.FriendlyMessageWithTag(),
				}
				c.JSON(http.StatusInternalServerError, resp)
				return false
			}

			claims := parsedToken.Claims.(jwt.MapClaims)
			accountID, email = claims["account_id"].(string), claims["email"].(string)
			/*end parse jwt*/

			/*start check from redis*/
			//check if access token exists
			accessToken, _ := helper.RedisHelper.Get(claims["access_uuid"].(string))
			if accessToken == "" {
				//token is expired
				resp := AuthResponse{ErrorType: "token_expired"}
				c.JSON(http.StatusUnauthorized, resp)
				return false
			}

			c.Set("account_id", accountID)
			c.Set("email", email)
		} else {
			_ = cerror.New("AUM02", errors.New("token_not_found"), "token_not_found") //only need to print the error
			resp := AuthResponse{
				ErrorType: "unauthorized",
				Message:   "token_not_found",
			}
			c.JSON(http.StatusUnauthorized, resp)
			return false
		}
	}

	return true
}
