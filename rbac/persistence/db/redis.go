package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/peifengll/casbin_demo/rbac/cacheadapter"
)

func DelPolicyInRedis(redisCli *redis.Client) (err error) {
	ctx := context.Background()
	if err = redisCli.Del(ctx, cacheadapter.Policy_Key).Err(); err != nil {
		return err
	}
	return
}
