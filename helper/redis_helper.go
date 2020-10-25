package helper

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/pajri/personal-backend/adapter/cerror"
	"github.com/pajri/personal-backend/config"
	"github.com/pajri/personal-backend/global"
)

type IRedis interface {
	Set(key string, value interface{}, exp int64) error
	Get(key string) (string, error)
}

type Redis struct {
	Client redis.Conn
}

var RedisHelper IRedis

func InitRedis() {
	RedisHelper = NewRedisHelper()
}

func NewRedisHelper() IRedis {
	//Connect
	var err error
	address := fmt.Sprintf("%s:%v", config.Config.Redis.Host, config.Config.Redis.Port)
	client, err := redis.Dial("tcp", address)
	if err != nil {
		log.Fatal("error redis.Dial : ", err)
	}

	RedisHelper = Redis{Client: client}

	// response, err := c.Do("AUTH", config.Config.Redis.Password)
	// if err != nil {
	// 	log.Fatal("redis auth error : ", err)
	// }

	// fmt.Println("Redis connected ", response)

	return RedisHelper
}

func (rh Redis) Set(key string, value interface{}, exp int64) error {
	_, err := rh.Client.Do("SET", key, value)
	if err != nil {
		return cerror.NewAndPrintWithTag("SRV00", err, global.FRIENDLY_MESSAGE)
	}

	if exp != 0 {
		_, err = rh.Client.Do("EXPIREAT", key, exp)
		if err != nil {
			return cerror.NewAndPrintWithTag("SRV01", err, global.FRIENDLY_MESSAGE)
		}
	}
	return nil
}

func (rh Redis) Get(key string) (string, error) {
	value, err := redis.String(rh.Client.Do("GET", key))
	if err != nil {
		return "", cerror.NewAndPrintWithTag("GRV00", err, global.FRIENDLY_MESSAGE)
	}
	return value, nil
}
