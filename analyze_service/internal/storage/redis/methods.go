package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	SetTTL(ctx context.Context, duration time.Duration, key string, value any) error
	Get(ctx context.Context, key string, value any) error
	Del(ctx context.Context, key string) error
	DelAll(ctx context.Context, key string) error
}

func (r *RedisDB) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling data before writing to cache: %w", err)
	}

	cmd := r.rdb.Set(ctx, key, data, expiration)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

//func (r *RedisDB) SetTTL(ctx context.Context, duration time.Duration, key string, value any) error {
//	data, err := json.Marshal(value)
//	if err != nil {
//		return fmt.Errorf("error marshaling data before writing to cache: %w", err)
//	}
//
//	cmd := r.rdb.SetEX(ctx, key, data, duration)
//	if err := cmd.Err(); err != nil {
//		return err
//	}
//
//	return nil
//}

func (r *RedisDB) Del(ctx context.Context, key string) error {
	cmd := r.rdb.Del(ctx, key)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisDB) DelAll(ctx context.Context, key string) error {
	cmd := r.rdb.Keys(ctx, key)
	if err := cmd.Err(); err != nil {
		return err
	}

	for _, key := range cmd.Val() {
		if c := r.rdb.Del(ctx, key); c.Err() != nil {
			return fmt.Errorf("error while deleting keys: %w", c.Err())
		}
	}

	return nil
}

var ErrRedisNoResult = errors.New("redis no result")

func (r *RedisDB) Get(ctx context.Context, key string, value any) error {
	var data string

	err := r.rdb.Get(ctx, key).Scan(&data)
	if errors.Is(err, redis.Nil) {
		return ErrRedisNoResult
	}

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), value)
	if err != nil {
		return fmt.Errorf("error unmarshal data after reading from cache: %w", err)
	}

	return nil
}
