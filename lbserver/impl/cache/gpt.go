package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/bgg/lbserver/impl/constant"
	"github.com/oldbai555/lbtool/log"
	"time"
)

// GetGptResult 获取gpt的结果
func GetGptResult(ctx context.Context, uuid string) (content string, err error) {
	content, err = rdb.Get(ctx, uuid).Result()
	if err != nil && err != redis.Nil {
		log.Errorf("err:%v", err)
		return
	}
	// 没错就结束咯
	if err == nil {
		err = rdb.Del(ctx, uuid).Err()
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
		return
	}
	// 拿不到就给个默认的话语
	return fmt.Sprintf(constant.SpeechQueueStartTemplate, uuid), nil
}

// SetGptResult 写入gpt的结果
func SetGptResult(ctx context.Context, uuid string, content string) error {
	err := rdb.Set(ctx, uuid, content, time.Hour).Err()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
