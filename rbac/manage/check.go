package manage

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
)

type Checker struct {
	redisCli *redis.Client
	uuid     string
	prefix   string
}

func NewChecker(redisCli *redis.Client, uuid, prefix string) *Checker {
	return &Checker{
		redisCli: redisCli,
		uuid:     uuid,
		prefix:   prefix,
	}
}

// CheckConsist Check for consistency
func (c *Checker) CheckConsist(e *casbin.Enforcer) error {
	var err error
	key := c.prefix + c.uuid
	ctx := context.Background()
	val, err := c.redisCli.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return c.loadAndSet(ctx, e, key)
		} else {
			return err
		}
	} else if val != "load" {
		return c.loadAndSet(ctx, e, key)
	}
	return nil
}
func (c *Checker) loadAndSet(ctx context.Context, e *casbin.Enforcer, key string) error {
	err2 := e.LoadPolicy()
	if err2 != nil {
		return err2
	}
	err2 = c.redisCli.Set(ctx, key, "load", 0).Err()
	if err2 != nil {
		return err2
	}
	return nil
}
