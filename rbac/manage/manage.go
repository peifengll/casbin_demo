package manage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Handler struct {
	tableName string
	redisCli  *redis.Client
	db        *gorm.DB
	prefix    string
}

// 先完成一个删除功能，然后看会不会重新加载
//func (h *Handler) DelPolicy() error {
//
//}

// 在完成操作之后
func (h *Handler) AfterOp() error {
	ctx := context.Background()
	// 目前就是删除所有的preifix为前缀的东西
	return h.deleteKeysByPrefix(ctx)
}
func (h *Handler) deleteKeysByPrefix(ctx context.Context) error {
	var cursor uint64
	keys, _, err := h.redisCli.Scan(ctx, cursor, h.prefix+"*", 0).Result()
	if err != nil {
		return err
	}
	// 删除匹配的键
	if len(keys) > 0 {
		err = h.redisCli.Del(ctx, keys...).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
