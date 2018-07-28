package redis

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/penggy/EasyGoLib/utils"
)

var Client *redis.Client

func Init() (err error) {
	sec := utils.Conf().Section("redis")
	host := sec.Key("host").MustString("localhost")
	port := sec.Key("port").MustInt(6379)
	auth := sec.Key("auth").MustString("")
	db := sec.Key("db").MustInt(0)

	log.Printf("redis server --> redis://%s:%d/db%d", host, port, db)

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: auth,
		DB:       db,
	})
	if _, e := Client.Ping().Result(); e != nil {
		err = fmt.Errorf("redis connect failed, %v", e)
		return
	}
	return
}

func Close() (err error) {
	if Client != nil {
		err = Client.Close()
		if err != nil {
			return
		}
		Client = nil
	}
	return
}
