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
	Username string
}

func (rc *RedisConfig) StartRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     rc.Addr,
		Password: rc.Password,
		DB:       rc.DB,
		Username: rc.Username,
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
		Addr:     "redis-11068.c261.us-east-1-4.ec2.redns.redis-cloud.com:11068",
		Username: "default",
		Password: "DM4iBq2pTYNQkweBZmwqGKyDYrj872M8",
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
