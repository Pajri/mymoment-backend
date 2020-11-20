package global

import (
	"log"
	"os"
)

var (
	Env string
	WD  string
)

func InitEnv() {
	env := os.Getenv("PERSONAL_ENV")
	if env == "" {
		env = "dev"
	}

	Env = env
	// Env = "sta"
}

func InitWD() {
	var err error
	WD, err = os.Getwd()
	if err != nil {
		log.Fatal("init wd error : ", err)
	}
}

func IsEnvDevelopment() bool {
	return Env == "" || Env == "dev"
}
