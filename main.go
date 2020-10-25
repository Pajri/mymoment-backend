package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/config"
	"github.com/pajri/personal-backend/db"
	"github.com/pajri/personal-backend/domain"
	"github.com/pajri/personal-backend/helper"
	"github.com/stretchr/stew/slice"

	_postDelivery "github.com/pajri/personal-backend/post/delivery"
	_postRepository "github.com/pajri/personal-backend/post/repository/mysql"
	_postUsecase "github.com/pajri/personal-backend/post/usecase"

	_accountRepository "github.com/pajri/personal-backend/account/repository/mysql"
	_authDelivery "github.com/pajri/personal-backend/auth/delivery"
	_authUsecase "github.com/pajri/personal-backend/auth/usecase"
	_profileRepository "github.com/pajri/personal-backend/profile/repository/mysql"
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

	postRepo := _postRepository.NewMySqlPostRepository(dbConn)
	postUsecase := _postUsecase.NewPostUseCase(postRepo)

	accountRepo := _accountRepository.NewMySqlAccountRepository(dbConn)
	profileRepo := _profileRepository.NewMySqlProfileRepository(dbConn)
	mailHelper := helper.NewEmailHelper()
	authUsecase := _authUsecase.NewAuthUsecase(accountRepo, profileRepo, mailHelper)

	r.Use(middleware(authUsecase))
	_postDelivery.NewPostHandler(r, postUsecase)
	_authDelivery.NewAuthHandler(r, authUsecase)

	r.Run()
}

func middleware(useCase domain.IAuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("here")
		if !slice.Contains(excludedFromAuth, c.FullPath()) {
			authArr := c.Request.Header["Authorization"]
			if len(authArr) > 0 {
				token := authArr[0]

				/*start parse jwt*/
				parsedToken, err := useCase.ParseJWT(token)
				if err != nil {
					err.(cerror.Error).PrintErrorWithTag()
					response := struct {
						Message string `json:"message"`
					}{
						"Authentication was not succesful [AUM00]",
					}

					c.JSON(http.StatusUnauthorized, response)
					c.Abort()
					return
				}
				/*end parse jwt*/

				/*start check from redis*/
				accessToken, err := helper.RedisHelper.Get(parsedToken["access_uuid"].(string))
				if err != nil || accessToken == "" {
					if err != nil {
						err.(cerror.Error).PrintErrorWithTag()
					}

					response := struct {
						Message string `json:"message"`
					}{
						"Authentication was not succesful [AUM01]",
					}

					c.JSON(http.StatusUnauthorized, response)
					c.Abort()
					return
				}

			} else {
				log.Println("token not found")
			}
		}
		c.Next()
	}
}
