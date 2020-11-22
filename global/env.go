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
}

func InitWD() {
	var err error
	WD, err = os.Getwd()
	if err != nil {
		log.Fatal("init wd error : ", err)
	}
}

func EnvFileName() string {
	if IsEnvDevelopment() {
		return ".env"
	}

	return Env + ".env"
}

func IsEnvDevelopment() bool {
	return Env == "" || Env == "dev"
}
