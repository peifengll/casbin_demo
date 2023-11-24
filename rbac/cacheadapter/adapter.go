package cacheadapter

import (
	"context"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/casbin/casbin/v2/util"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"

	"strings"
)

type CacheAdapter struct {
	// CasbinIndicate 这个是键的值，
	//其先要查看这个看自己需要做什么操作
	// 在初始化的时候生成，使用uuid保证这个实例的是否需要更新？
	// 如果redis里边没有这个 键 ，那么一定就是要更新
	CasbinIndicate string
	tableName      string
	//databaseName   string
	//r_adapter      *redisadapter.Adapter
	redisCli *redis.Client
	db       *gorm.DB
}

func NewAdapter(db *gorm.DB, r *redis.Client, params ...interface{}) (*CacheAdapter, error) {
	a := &CacheAdapter{
		CasbinIndicate: Casbin_Indicate_Prefix,
		tableName:      defaultTableName,
		redisCli:       r,
		db:             db,
	}
	return a, nil
}

func (a *CacheAdapter) LoadPolicy(model model.Model) (err error) {
	ctx := context.Background()
	// 先看policy的key存不存在
	exists, err := a.redisCli.Exists(ctx, Policy_Key).Result()
	if err != nil {
		log.Fatalf("redis.Exists err:%v \n", err)
	}
	if exists == 1 {
		//	 key存在
		return a.loadPolicy(ctx, model, persist.LoadPolicyArray)
	} else {
		//	 查数据库，然后放进redis
		var lines []CasbinRule
		if err := a.db.Order("ID").Find(&lines).Error; err != nil {
			return err
		}
	}
	return nil
}

func (a *CacheAdapter) loadPolicy(ctx context.Context, model model.Model, handler func([]string, model.Model)) (err error) {
	// 0, -1 fetches all entries from the list
	rules, err := a.redisCli.LRange(ctx, Policy_Key, 0, -1).Result()
	if err != nil {
		return err
	}
	// Parse the rules from Redis
	for _, rule := range rules {
		handler(strings.Split(rule, ", "), model)
	}
	return
}

// AddPolicy adds a policy rule to the storage.
func (a *CacheAdapter) AddPolicy(_ string, ptype string, rule []string) (err error) {
	ctx := context.Background()
	err = a.addPolicyToMysql(ctx, ptype, rule)
	if err != nil {
		return err
	}
	err = a.addPolicyToRedis(ctx, buildRuleStr(ptype, rule))
	return
}

func (a *CacheAdapter) addPolicyToRedis(ctx context.Context, rule string) (err error) {
	if err = a.redisCli.RPush(ctx, Policy_Key, rule).Err(); err != nil {
		return err
	}
	return
}
func (a *CacheAdapter) addPolicyToMysql(ctx context.Context, ptype string, rule []string) (err error) {
	line := a.toPolicyLine(ptype, rule)
	err = a.db.Create(&line).Error
	return
}

func buildRuleStr(ptype string, rule []string) string {
	return ptype + ", " + util.ArrayToString(rule)
}

func (a *CacheAdapter) toPolicyLine(ptype string, rule []string) *CasbinRule {
	line := &CasbinRule{}

	line.Ptype = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}
