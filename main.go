package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pajri/personal-backend/config"
	"github.com/pajri/personal-backend/db"
	"github.com/pajri/personal-backend/global"
	"github.com/pajri/personal-backend/helper"
	"github.com/pajri/personal-backend/middleware"

	_postDelivery "github.com/pajri/personal-backend/post/delivery"
	_postRepository "github.com/pajri/personal-backend/post/repository/mysql"
	_postUsecase "github.com/pajri/personal-backend/post/usecase"

	_accountRepository "github.com/pajri/personal-backend/account/repository/mysql"
	_authDelivery "github.com/pajri/personal-backend/auth/delivery"
	_authUsecase "github.com/pajri/personal-backend/auth/usecase"

	_profileDelivery "github.com/pajri/personal-backend/profile/delivery"
	_profileRepository "github.com/pajri/personal-backend/profile/repository/mysql"
	_profileUsecase "github.com/pajri/personal-backend/profile/usecase"

	_imageDelivery "github.com/pajri/personal-backend/image/delivery"
	_imageRepository "github.com/pajri/personal-backend/image/repository/mysql"
	_imageUsecase "github.com/pajri/personal-backend/image/usecase"
)

func main() {
	//init config
	config.InitConfig()

	global.InitEnv()
	global.InitWD()

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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.Config.FEHost},
		AllowMethods:     []string{"GET", "POST"},
		AllowCredentials: true,
	}))
	r.Static("/upload/images/", "./upload/images")

	//setup helper
	mailHelper := helper.NewEmailHelper()

	//setup repo and usecase
	imageRepo := _imageRepository.NewMySqlImageRepository(dbConn)
	imageUsecase := _imageUsecase.NewImageUsecase(imageRepo)

	postRepo := _postRepository.NewMySqlPostRepository(dbConn)
	postUsecase := _postUsecase.NewPostUseCase(postRepo, imageRepo)

	accountRepo := _accountRepository.NewMySqlAccountRepository(dbConn)

	profileRepo := _profileRepository.NewMySqlProfileRepository(dbConn)
	profileUsecase := _profileUsecase.NewProfileUsecase(accountRepo, profileRepo)

	authUsecase := _authUsecase.NewAuthUsecase(accountRepo, profileRepo, mailHelper)

	r.Use(middleware.Middleware(authUsecase))
	_postDelivery.NewPostHandler(r, postUsecase)
	_authDelivery.NewAuthHandler(r, authUsecase)
	_imageDelivery.NewImageHandler(r, imageUsecase)
	_profileDelivery.NewProfileHandler(r, profileUsecase)

	r.Run(":5000")

}
