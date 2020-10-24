package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pajri/personal-backend/config"
	"github.com/pajri/personal-backend/db"
	"github.com/pajri/personal-backend/helper"
	"github.com/pajri/personal-backend/redis"
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
	redis.InitRedis()
	defer redis.Client.Close()
	/*end init redis*/

	r := gin.Default()
	r.Use(middleware())

	postRepo := _postRepository.NewMySqlPostRepository(dbConn)
	postUsecase := _postUsecase.NewPostUseCase(postRepo)
	_postDelivery.NewPostHandler(r, postUsecase)

	accountRepo := _accountRepository.NewMySqlAccountRepository(dbConn)
	profileRepo := _profileRepository.NewMySqlProfileRepository(dbConn)
	mailHelper := helper.NewEmailHelper()
	authUsecase := _authUsecase.NewAuthUsecase(accountRepo, profileRepo, mailHelper)
	_authDelivery.NewAuthHandler(r, authUsecase)

	r.Run()
}

func middleware() gin.HandlerFunc {
	secret := os.Getenv("JWT_SECRET")
	fmt.Println(secret)
	return func(c *gin.Context) {
		if !slice.Contains(excludedFromAuth, c.FullPath()) {

		}
		c.Next()
	}
}
