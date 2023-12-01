package cacheadapter

import (
	"context"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/casbin/casbin/v2/util"
	"github.com/go-redis/redis/v8"
	"github.com/peifengll/casbin_demo/rbac/convterter"
	"github.com/peifengll/casbin_demo/rbac/persistence/po"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"

	"strings"
)

type CacheAdapter struct {
	// CasbinIndicate 这个是键的值，
	//其先要查看这个看自己需要做什么操作
	// 在初始化的时候生成，使用uuid保证这个实例的是否需要更新？
	// 如果redis里边没有这个 键 ，那么一定就是要更新
	//CasbinIndicate string
	tableName string
	//databaseName   string
	//r_adapter      *redisadapter.Adapter
	redisCli *redis.Client
	db       *gorm.DB
}

func NewAdapter(db *gorm.DB, r *redis.Client, params ...interface{}) *CacheAdapter {
	a := &CacheAdapter{
		//CasbinIndicate: Casbin_Indicate_Prefix,
		tableName: defaultTableName,
		redisCli:  r,
		db:        db,
	}
	return a
}

func (a *CacheAdapter) LoadPolicy(model model.Model) (err error) {
	ctx := context.Background()
	// 先看policy的key存不存在
	exists, err := a.redisCli.Exists(ctx, Policy_Key).Result()
	if err != nil {
		log.Fatalf("redis.Exists err:%v \n", err)
	}
	//fmt.Println("执行*1")
	if exists == 1 {
		LoadFormCounter.With(prometheus.Labels{"from": "redis"}).Inc()
		//fmt.Println("+++++")
		//	 key存在
		return a.loadPolicyRedis(ctx, model, persist.LoadPolicyArray)
	} else {
		//fmt.Println("-------")
		//	 查数据库，然后放进redis
		var lines []po.CasbinRule
		LoadFormCounter.With(prometheus.Labels{"from": "mysql"}).Inc()
		if err := a.db.Order("ID").Find(&lines).Error; err != nil {
			return err
		}
		// 放进redis，line是持有了数据库所有规则
		for _, v := range lines {

			err := loadPolicyLine(v, model)
			if err != nil {
				return err
			}
			err = a.addPolicyToRedis(ctx, a.toPolicyRuleStr(&v))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func loadPolicyLine(line po.CasbinRule, model model.Model) error {
	var p = []string{line.Ptype,
		line.V0, line.V1, line.V2,
		line.V3, line.V4, line.V5}

	index := len(p) - 1
	for p[index] == "" {
		index--
	}
	index += 1
	p = p[:index]
	persist.LoadPolicyArray(p, model)
	return nil
}

func (a *CacheAdapter) loadPolicyRedis(ctx context.Context, model model.Model, handler func([]string, model.Model)) (err error) {
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
	line := convterter.ToPolicyPo(ptype, rule)
	err = a.db.Create(&line).Error
	return
}

func (a *CacheAdapter) SavePolicy(model model.Model) error {
	// 直接删除redis和mysql，然后再往mysql添加即可
	err := a.savePolicy(model)
	return err
}
func (a *CacheAdapter) savePolicy(model model.Model) error {
	var err error
	tx := a.db.Clauses(dbresolver.Write).Begin()
	sql := fmt.Sprintf("truncate table %s", a.tableName)
	err = a.db.Exec(sql).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var lines []po.CasbinRule
	flushEvery := 1000
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			lines = append(lines, convterter.ToPolicyPo(ptype, rule))
			if len(lines) > flushEvery {
				if err := tx.Create(&lines).Error; err != nil {
					tx.Rollback()
					return err
				}
				lines = nil
			}
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			lines = append(lines, convterter.ToPolicyPo(ptype, rule))
			if len(lines) > flushEvery {
				if err := tx.Create(&lines).Error; err != nil {
					tx.Rollback()
					return err
				}
				lines = nil
			}
		}
	}

	if len(lines) > 0 {
		if err := tx.Create(&lines).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}
	// 现在来删掉redis
	err = a.delPolicyInRedis()
	if err != nil {
		return err
	}
	return err

}

func (a *CacheAdapter) delPolicyInRedis() (err error) {
	ctx := context.Background()
	if err = a.redisCli.Del(ctx, Policy_Key).Err(); err != nil {
		return err
	}
	return
}

// 注意redis 删除的时候要放四条规则进去，有一个是deny什么的
func (a *CacheAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	//	 同时删掉mysql跟redis
	line := convterter.ToPolicyPo(ptype, rule)
	err := a.rawDelete(a.db, line)
	if err != nil {
		return err
	}
	err = a.removePolicyInRedis(buildRuleStr(ptype, rule))
	//fmt.Println(buildRuleStr(ptype, rule))
	return err
}

func (a *CacheAdapter) removePolicyInRedis(rule string) (err error) {
	ctx := context.Background()
	if err = a.redisCli.LRem(ctx, Policy_Key, 1, rule).Err(); err != nil {
		return err
	}
	return
}

func (a *CacheAdapter) rawDelete(db *gorm.DB, line po.CasbinRule) error {
	queryArgs := []interface{}{line.Ptype}

	queryStr := "ptype = ?"
	if line.V0 != "" {
		queryStr += " and v0 = ?"
		queryArgs = append(queryArgs, line.V0)
	}
	if line.V1 != "" {
		queryStr += " and v1 = ?"
		queryArgs = append(queryArgs, line.V1)
	}
	if line.V2 != "" {
		queryStr += " and v2 = ?"
		queryArgs = append(queryArgs, line.V2)
	}
	if line.V3 != "" {
		queryStr += " and v3 = ?"
		queryArgs = append(queryArgs, line.V3)
	}
	if line.V4 != "" {
		queryStr += " and v4 = ?"
		queryArgs = append(queryArgs, line.V4)
	}
	if line.V5 != "" {
		queryStr += " and v5 = ?"
		queryArgs = append(queryArgs, line.V5)
	}
	args := append([]interface{}{queryStr}, queryArgs...)
	err := db.Delete(&po.CasbinRule{}, args...).Error
	return err
}

// RemoveFilteredPolicy 从持久层删除符合筛选条件的policy规则
// 删除了之后可以考虑直接清除redis，或是也一个一个的找到并删除
func (a *CacheAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := &po.CasbinRule{}
	line.Ptype = ptype
	if fieldIndex == -1 {
		return a.rawDelete(a.db, *line)
	}
	err := checkQueryField(fieldValues)
	if err != nil {
		return err
	}
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
	}
	err = a.rawDelete(a.db, *line)
	if err != nil {
		return err
	}
	err = a.delPolicyInRedis()
	return err
}

func checkQueryField(fieldValues []string) error {
	for _, fieldValue := range fieldValues {
		if fieldValue != "" {
			return nil
		}
	}
	return errors.New("the query field cannot all be empty string (\"\"), please check")
}

func buildRuleStr(ptype string, rule []string) string {
	return ptype + ", " + util.ArrayToString(rule)
}

//
//func (a *CacheAdapter) toPolicyLine(ptype string, rule []string) po.CasbinRule {
//	line := &po.CasbinRule{}
//
//	line.Ptype = ptype
//	if len(rule) > 0 {
//		line.V0 = rule[0]
//	}
//	if len(rule) > 1 {
//		line.V1 = rule[1]
//	}
//	if len(rule) > 2 {
//		line.V2 = rule[2]
//	}
//	if len(rule) > 3 {
//		line.V3 = rule[3]
//	}
//	if len(rule) > 4 {
//		line.V4 = rule[4]
//	}
//	if len(rule) > 5 {
//		line.V5 = rule[5]
//	}
//
//	return *line
//}

func (a *CacheAdapter) toPolicyRuleStr(line *po.CasbinRule) (res string) {
	if line.V5 != "" {
		res = buildRuleStr(line.Ptype, []string{line.V0, line.V1, line.V2, line.V3, line.V4, line.V5})
	} else if line.V4 != "" {
		res = buildRuleStr(line.Ptype, []string{line.V0, line.V1, line.V2, line.V3, line.V4})
	} else if line.V3 != "" {
		res = buildRuleStr(line.Ptype, []string{line.V0, line.V1, line.V2, line.V3})
	} else if line.V2 != "" {
		res = buildRuleStr(line.Ptype, []string{line.V0, line.V1, line.V2})
	} else if line.V1 != "" {
		res = buildRuleStr(line.Ptype, []string{line.V0, line.V1})
	} else if line.V0 != "" {
		res = buildRuleStr(line.Ptype, []string{line.V0})
	}
	return

}

// AddMorePolicy （will delete） 增加1w个策略
func (a *CacheAdapter) AddMorePolicy() {
	for i := 0; i < 10000; i++ {
		line := &po.CasbinRule{
			Ptype: "p",
			V0:    "zhangsan",
			V1:    fmt.Sprintf("data%d", i),
			V2:    "read",
			V3:    "allow",
		}
		a.db.Create(&line)
	}
}
