package redisx

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/tiny-sky/Tvote/log"
)

type Redis struct {
	Client *redis.Client
}

var Rdb Redis

type Settings struct {
	Host            string `json:"host" yaml:"host"`
	Password        string `json:"password" yaml:"password"`
	DB              int    `json:"db" yaml:"db"`
	MaxRetries      int    `json:"maxRetries" yaml:"maxRetries"`
	MinRetryBackoff int    `json:"minRetryBackoff" yaml:"minRetryBackoff"`
	MaxRetryBackoff int    `json:"maxRetryBackoff" yaml:"maxRetryBackoff"`
	PoolSize        int    `json:"poolSize" yaml:"poolSize"`
	MinIdleConns    int    `json:"minIdleConns" yaml:"minIdleConns"`
}

func (s *Settings) Init() {
	options := &redis.Options{
		Addr:            s.Host,
		Password:        s.Password,
		DB:              s.DB,
		MaxRetries:      s.MaxRetries,
		MinRetryBackoff: time.Duration(s.MinRetryBackoff) * time.Second,
		MaxRetryBackoff: time.Duration(s.MaxRetryBackoff) * time.Second,
		PoolSize:        s.PoolSize,
		MinIdleConns:    s.MinIdleConns,
	}

	Rdb.Client = redis.NewClient(options)

	_, err := Rdb.Client.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
