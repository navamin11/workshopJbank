package configs

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

func ConnectRedis() *redis.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Error loading .env file")
	}
	// Conncet to Redis
	RdHost := os.Getenv("REDIS_HOST")
	RdPassword := os.Getenv("REDIS_PASSWORD")
	RdIndex := 1 // Redis support 16 database. You can switch a DB using an integer starting from 0 to 15
	RdPort := os.Getenv("REDIS_PORT")

	addConnect := fmt.Sprintf("%s:%s", RdHost, RdPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addConnect,
		Password: RdPassword,
		DB:       RdIndex,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Error(err)
	}

	return rdb
}
