package redis

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/pajri/personal-backend/config"
)

var Client redis.Conn

func InitRedis() {
	//Connect
	var err error
	address := fmt.Sprintf("%s:%v", config.Config.Redis.Host, config.Config.Redis.Port)
	Client, err = redis.Dial("tcp", address)
	if err != nil {
		log.Fatal("error redis.Dial : ", err)
	}

	// response, err := c.Do("AUTH", config.Config.Redis.Password)
	// if err != nil {
	// 	log.Fatal("redis auth error : ", err)
	// }

	// fmt.Println("Redis connected ", response)
}
