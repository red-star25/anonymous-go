package config

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func (rc *RedisConfig) StartRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     rc.Addr,
		Password: rc.Password,
		DB:       rc.DB,
	})

	RedisClient = client

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	ping, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to redis")
	}

	log.Println(ping)
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}

func SetRedisToken(key string, value string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	err := RedisClient.Set(ctx, key, value, time.Second*180).Err()
	if err != nil {
		log.Fatal(err)
	}
}
