package config

import (
	"context"
	"log"
	"os"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to redis")
	}

	log.Println("Connected to redis")
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
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
