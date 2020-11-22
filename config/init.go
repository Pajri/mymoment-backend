package config

import (
	"log"

	"github.com/pajri/personal-backend/global"
	"github.com/tkanos/gonfig"
)

var Config Configuration

func InitConfig() {
	config := Configuration{}
	err := gonfig.GetConf("config/"+configFileName(), &config)
	if err != nil {
		log.Fatal("unable to read config : ", err)
	}

	Config = config
}

func configFileName() string {
	if global.IsEnvDevelopment() {
		return "config.json"
	}

	return "config." + global.Env + ".json"
}
