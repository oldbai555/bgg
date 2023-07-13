package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	webtool "github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/lbtool/log"
)

type Cache struct {
	*redis.Client
}

type HashItem struct {
	k string
	v interface{}
}

func (h *HashItem) String() string {
	return fmt.Sprintf("HashItem{k:%v,v:%v}", h.k, h.v)
}

func NewCache(r *webtool.RedisConf, prefix string) *Cache {
	return &Cache{
		Client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
			Password: r.Password,
			DB:       r.Database,
		}),
	}
}

func NewHashItem(k string, v interface{}) *HashItem {
	return &HashItem{k: k, v: v}
}

func (c *Cache) HSet(ctx context.Context, hk string, items ...*HashItem) error {
	var args []interface{}

	for i := range items {
		item := items[i]
		itemBytes, err := json.Marshal(item.v)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		args = append(args, fmt.Sprintf("%v", item.k), itemBytes)
	}

	err := c.Client.HSet(ctx, hk, args...).Err()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (c *Cache) HDel(ctx context.Context, hk string) error {
	err := c.Client.HDel(ctx, hk).Err()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (c *Cache) HGet(ctx context.Context, hk string, subKey string, out interface{}) error {
	res := c.Client.HGet(ctx, hk, subKey)
	if res.Err() != nil {
		log.Errorf("err:%v", res.Err())
		return res.Err()
	}
	err := json.Unmarshal([]byte(res.String()), &out)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (c *Cache) HGetAll(ctx context.Context, hk string, fn func(ctx context.Context, item *HashItem) error) error {
	result, err := c.Client.HGetAll(ctx, hk).Result()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	for k, v := range result {
		err := fn(ctx, NewHashItem(k, v))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}
