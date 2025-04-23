package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis - contains parameters for connecting to Redis
type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

type RedisDB struct {
	rdb *redis.Client
}

// Connect - creates a connection to Redis
func (r *Redis) Connect() (*RedisDB, error) {
	cl := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       0, // use default DB

		MinIdleConns: 3,
	})

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Minute))

	_, err := cl.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisDB{rdb: cl}, nil

}
