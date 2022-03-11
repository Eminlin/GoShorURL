package common

import "github.com/go-redis/redis"

var RedisClient *redis.Client

func init() {
	RedisClient = newRedisClient()
}

func newRedisClient() *redis.Client {
	c := AppConf.Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     c.HostPort,
		Password: c.Password,
		DB:       c.DB,
	})
	if err := RedisClient.Ping().Err(); err != nil {
		panic(err)
	}
	return RedisClient
}
