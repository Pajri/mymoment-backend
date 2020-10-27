package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/config"
	"github.com/pajri/personal-backend/db"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/global"
	"github.com/pajri/personal-backend/helper"
	"github.com/stretchr/stew/slice"

	_postDelivery "github.com/pajri/personal-backend/post/delivery"
	_postRepository "github.com/pajri/personal-backend/post/repository/mysql"
	_postUsecase "github.com/pajri/personal-backend/post/usecase"

	_accountRepository "github.com/pajri/personal-backend/account/repository/mysql"
	_authDelivery "github.com/pajri/personal-backend/auth/delivery"
	_authUsecase "github.com/pajri/personal-backend/auth/usecase"
	_profileRepository "github.com/pajri/personal-backend/profile/repository/mysql"

	_imageDelivery "github.com/pajri/personal-backend/image/delivery"
	_imageRepository "github.com/pajri/personal-backend/image/repository/mysql"
	_imageUsecase "github.com/pajri/personal-backend/image/usecase"
)

var excludedFromAuth = []string{
	"/api/auth/login",
	"/api/auth/signup",
	"/api/auth/verify_email",
	"/api/auth/reset_password/",
	"/api/auth/change_password",
}

func main() {
	//init config
	config.InitConfig()

	/* start init db*/
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("unable to connect to db : ", err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal("error while pinging db : ", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	/*end init db*/

	/*start load env variable*/
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file : ", err)
	}
	/*end load env variable*/

	/*start init redis*/
	helper.InitRedis()
	defer helper.RedisHelper.(helper.Redis).Client.Close()
	/*end init redis*/

	r := gin.Default()
	//setup helper
	mailHelper := helper.NewEmailHelper()

	//setup repo and usecase
	postRepo := _postRepository.NewMySqlPostRepository(dbConn)
	postUsecase := _postUsecase.NewPostUseCase(postRepo)

	profileRepo := _profileRepository.NewMySqlProfileRepository(dbConn)

	accountRepo := _accountRepository.NewMySqlAccountRepository(dbConn)

	authUsecase := _authUsecase.NewAuthUsecase(accountRepo, profileRepo, mailHelper)

	imageRepo := _imageRepository.NewMySqlImageRepository(dbConn)
	imageUsecase := _imageUsecase.NewImageUsecase(imageRepo)

	r.Use(middleware(authUsecase, accountRepo))
	_postDelivery.NewPostHandler(r, postUsecase)
	_authDelivery.NewAuthHandler(r, authUsecase)
	_imageDelivery.NewImageHandler(r, imageUsecase)

	r.Run()
}

func middleware(useCase domain.IAuthUsecase, accountRepo domain.IAccountRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !slice.Contains(excludedFromAuth, c.FullPath()) {
			authArr := c.Request.Header["Authorization"]
			var accountID, email string

			if len(authArr) > 0 {
				token := authArr[0]

				/*start parse jwt*/
				_, parsedToken, err := useCase.ParseJWT(token)
				if err != nil {
					cerr, ok := err.(cerror.Error)
					if !ok {
						cerr = cerror.NewAndPrintWithTag("AUM00", err, global.FRIENDLY_MESSAGE)
					}

					if cerr.Type != cerror.TYPE_EXPIRED {
						authNotSuccessfulResponse(c, 0, cerr)
						return
					}
				}

				if parsedToken != nil {
					accountID = parsedToken["account_id"].(string)
					email = parsedToken["email"].(string)
				}
				/*end parse jwt*/

				/*start check from redis*/
				//check if access token exists
				accessToken, _ := helper.RedisHelper.Get(parsedToken["access_uuid"].(string))

				//check refresh token
				expire := parsedToken["exp"].(float64)
				unixInt := int64(expire)
				expTime := time.Unix(unixInt, 0)
				if accessToken == "" || time.Now().After(expTime) { //token is expired
					refreshToken, err := c.Request.Cookie("refresh_token")
					if err != nil {
						cerr, ok := err.(cerror.Error)
						if !ok {
							cerr = cerror.NewAndPrintWithTag("AUM02", err, global.FRIENDLY_MESSAGE)
						}
						authNotSuccessfulResponse(c, 0, cerr)
						return
					}

					parsedRfToken, _, err := useCase.ParseJWT(refreshToken.Value)
					if err != nil {
						cerr, ok := err.(cerror.Error)
						if !ok {
							cerr = cerror.NewAndPrintWithTag("AUM03", err, global.FRIENDLY_MESSAGE)
						}

						authNotSuccessfulResponse(c, 0, cerr)
						return
					}

					filter := domain.AccountFilter{
						AccountID: parsedRfToken.Claims.(jwt.MapClaims)["account_id"].(string),
					}
					account, err := accountRepo.GetAccount(filter)
					if err != nil || account == nil {
						cerr := cerror.NewAndPrintWithTag("AUM04", err, global.FRIENDLY_MESSAGE)
						authNotSuccessfulResponse(c, 0, cerr)
						return
					}

					accessTokenClaims := jwt.MapClaims{}
					accessTokenClaims["authorized"] = true
					accessTokenClaims["account_id"] = account.AccountID
					accessTokenClaims["access_uuid"] = uuid.New().String()
					accessTokenClaims["email"] = account.Email
					accessTokenClaims["exp"] = time.Now().Add(15 * time.Minute).Unix()

					accountID = account.AccountID
					email = account.Email

					refreshTokenClaims := jwt.MapClaims{}
					refreshTokenClaims["account_id"] = account.AccountID
					refreshTokenClaims["refresh_uuid"] = uuid.New().String()
					refreshTokenClaims["exp"] = time.Now().Add(1 * time.Hour).Unix()

					jwtHelper := helper.JWTHelper{}
					token, err := jwtHelper.CreateTokenPair(accessTokenClaims, refreshTokenClaims)
					if err != nil {
						cerr := cerror.NewAndPrintWithTag("AUM05", err, global.FRIENDLY_MESSAGE)
						authNotSuccessfulResponse(c, 0, cerr)
						return
					}

					tokenByte, _ := json.Marshal(token)
					tokenString := string(tokenByte)
					c.SetCookie("token", tokenString, 0, "", "", false, false)
				} else {
					if accessToken != token { //compare token from redis with token from header
						cerr := cerror.NewAndPrintWithTag("AUM06",
							errors.New("token from header is different from token from redis"), "")
						authNotSuccessfulResponse(c, 0, cerr)
						return
					}
				}

				c.Set("account_id", accountID)
				c.Set("email", email)
			} else {
				log.Println("token not found")
			}
		}
		c.Next()
	}
}

func authNotSuccessfulResponse(c *gin.Context, status int, err cerror.Error) {
	if status == 0 {
		status = http.StatusUnauthorized
	}
	response := struct {
		Message string `json:"message"`
	}{err.FriendlyMessageWithTag()}

	c.JSON(status, response)
	c.Abort()
	return
}
