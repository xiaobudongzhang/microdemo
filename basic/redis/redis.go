package redis

import (
	"sync"

	"github.com/go-log/log"
	"github.com/go-redis/redis"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part2/basic/config"
)

var (
	client *redis.Client
	m      sync.RWMutex
	inited bool
)

func Init() {
	m.Lock()
	defer m.Unlock()

	if inited {
		log.Log("已经初始化过redis....")
		return
	}

	redisConfig := config.GetRedisConfig()

	if redisConfig != nil && redisConfig.GetEnabled() {
		log.Log("初始化redis...")

		if redisConfig.GetSentinelConfig() != nil && redisConfig.GetSentinelConfig().GetEnabled() {
			log.Log("初始化哨兵模式...")
			initSentinel(redisConfig)
		} else {
			log.Log("初始化普通模式...")
			initSingle(redisConfig)
		}

		log.Log("初始化redis 检测连接...")

		pong, err := client.Ping().Result()
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Log("ping %s", pong)
	}

	inited = true
}

func GetReis() *redis.Client {
	return client
}

func initSentinel(redisConfig config.RedisConfig) {
	client = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: redisConfig.GetSentinelConfig().GetMaster(),
		SentinelAddr: redisConfig.GetSentinelConfig().GetNodes(),
		DB: redisConfig.GetDBNum(),
		Password: redisConfig.GetPassword()
	})
}

func initSingle(redisConfig config.RedisConfig) {
	client = redis.NewClient(&redis.Options{
		Addr: redisConfig.GetConn(),
		Password:redisConfig.GetPassword(),
		DB:redisConfig.GetDBNum(),
	})
}
