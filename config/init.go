package config

import (
	"log"

	"github.com/tkanos/gonfig"
)

var Config Configuration

func InitConfig() {
	config := Configuration{}
	err := gonfig.GetConf("config/config.json", &config)
	if err != nil {
		log.Fatal("unable to read config : ", err)
	}

	Config = config
}
