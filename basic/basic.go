package basic

import (
	"user-srv/basic/config"
	"user-srv/basic/db"

	"github.com/micro-in-cn/tutorials/microservice-in-micro/part2/basic/redis"
)

func Init() {
	config.Init()
	db.Init()
	redis.Init()
}
