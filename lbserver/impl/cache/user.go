package cache

import (
	"context"
	"github.com/oldbai555/lbtool/log"
	"time"
)

func SetLoginUserToken(ctx context.Context, key, val string, expiresIn time.Duration) error {
	err := rdb.Set(ctx, key, val, expiresIn).Err()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func GetLoginUserToken(ctx context.Context, key string) (string, error) {
	token, err := rdb.Get(ctx, key).Result()
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	return token, nil
}

func DelLoginUserToken(ctx context.Context, key string) error {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func SetLoginUser(ctx context.Context, key string, val interface{}, expiresIn time.Duration) error {
	err := rdb.SetJson(ctx, key, val, expiresIn)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func GetLoginUser(ctx context.Context, key string, out interface{}) error {
	err := rdb.GetJson(ctx, key, out)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func DelLoginUser(ctx context.Context, key string) error {
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
