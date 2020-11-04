package global

import "os"

var Env string

func InitEnv() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	Env = env
}

func IsEnvDevelopment() bool {
	return Env == "" || Env == "dev"
}
