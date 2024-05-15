package providers

import (
	"github.com/gomodule/redigo/redis"
)

var (
	Pool *redis.Pool
)

func init() {
	Pool = newPool("redis://cache:6379/0")	
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(server)
		},
	}
}
