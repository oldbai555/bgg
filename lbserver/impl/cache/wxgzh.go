package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lbtool/log"
	"time"
)

// GetWxGzhAccessToken 获取 accessToken
func GetWxGzhAccessToken(ctx context.Context, key string) (string, error) {
	accessToken, err := rdb.Get(ctx, key).Result()

	if err != nil && err != redis.Nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	if err == nil && accessToken != "" {
		return accessToken, nil
	}
	return "", nil
}

// SetWxGzhAccessToken 设置 accessToken
func SetWxGzhAccessToken(ctx context.Context, key, val string, expiresIn time.Duration) error {
	err := rdb.Set(ctx, key, val, expiresIn).Err()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
