package manage

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/peifengll/casbin_demo/rbac/cacheadapter"
	"github.com/peifengll/casbin_demo/rbac/persistence/db"
	"github.com/peifengll/casbin_demo/rbac/persistence/po"
	"gorm.io/gorm"
)

type Manager struct {
	//tableName string
	policyDao db.CasbinRuleDao
	enforcer  *casbin.Enforcer
	redisCli  *redis.Client
	//db       *gorm.DB
	// 前缀为这个的，键(是各个服务的标记)
	prefix string
}

func NewManager(client *redis.Client, d *gorm.DB, prefix, path string) (*Manager, error) {
	var err error
	cacheada := cacheadapter.NewAdapter(d, client)
	e, err := casbin.NewEnforcer(path, cacheada)
	if err != nil {
		return nil, err
	}
	return &Manager{
		enforcer:  e,
		redisCli:  client,
		prefix:    prefix,
		policyDao: db.NewCasbinRuleDao(d),
	}, nil
}

// 管理服务无非就是增删查改
func (h *Manager) QuerryAll(ctx context.Context) ([]po.CasbinRule, error) {
	return h.policyDao.QuerryAll(ctx)
}
func (h *Manager) QueryRuleById(ctx context.Context, id int) (*po.CasbinRule, error) {
	return h.policyDao.QueryRuleById(ctx, id)
}
func (h *Manager) QueryRulesByIds(ctx context.Context, ids []int) ([]po.CasbinRule, error) {
	return h.policyDao.QueryRulesByIds(ctx, ids)
}
func (h *Manager) QueryRuleByRule(ctx context.Context, ptype string, rule []string) ([]po.CasbinRule, error) {
	return h.policyDao.QueryRuleByRule(ctx, ptype, rule)
}

//QueryRulesByPType(ctx context.Context, ptype string) ([]po.CasbinRule, error)
//QueryRulesByEvaluator(ctx context.Context, evaluator string) ([]po.CasbinRule, error)
//QueryRulesByResource(ctx context.Context, resource string) ([]po.CasbinRule, error)
//QueryRulesByAction(ctx context.Context, action string) ([]po.CasbinRule, error)

// 修改操作
func (h *Manager) UpdateRuleById(ctx context.Context, id int, ptype string, ru []string) error {
	if err := h.policyDao.UpdateRuleById(ctx, id, ptype, ru); err != nil {
		return err
	}
	return h.AfterOp()
}
func (h *Manager) UpdateRulesByIds(ctx context.Context, ids []int, ptype string, ru []string) error {
	if err := h.policyDao.UpdateRulesByIds(ctx, ids, ptype, ru); err != nil {
		return err
	}
	return h.AfterOp()
}
func (h *Manager) UpdateRulesByRule(ctx context.Context, ptype string, ru []string, rule []string) error {
	if err := h.policyDao.UpdateRulesByRule(ctx, ptype, ru, rule); err != nil {
		return err
	}
	return h.AfterOp()
}

// 增加  增加直接用 casbin的增加即可，批量增加？也可以
func (h *Manager) AddRule(ctx context.Context, pytepe string, rule []string) error {
	if err := h.policyDao.AddRule(ctx, pytepe, rule); err != nil {
		return err
	}
	return h.AfterOp()
}
func (h *Manager) AddRules(ctx context.Context, ptype string, rules [][]string) error {
	err := h.AddRules(ctx, ptype, rules)
	if err != nil {
		return err
	}
	return h.AfterOp()
}

// 删除操作，提供按照id删，和按照规则删
func (h *Manager) RemoveRuleById(ctx context.Context, id int) error {
	if err := h.policyDao.RemoveRuleById(ctx, id); err != nil {
		return err
	}
	return h.AfterOp()

}
func (h *Manager) RemoveRuleByIds(ctx context.Context, ids []int) error {
	err := h.policyDao.RemoveRuleByIds(ctx, ids)
	if err != nil {
		return err
	}
	return h.AfterOp()

}
func (h *Manager) RemoveRuleByRule(ctx context.Context, ptype string, rule []string) error {
	err := h.RemoveRuleByRule(ctx, ptype, rule)
	if err != nil {
		return err
	}
	return h.AfterOp()

}
func (h *Manager) RemoveRuleByRules(ctx context.Context, ptype string, rules [][]string) error {
	err := h.RemoveRuleByRules(
		ctx,
		ptype,
		rules,
	)
	if err != nil {
		return err
	}
	return h.AfterOp()
}

// 在完成操作之后
func (h *Manager) AfterOp() error {
	ctx := context.Background()
	// 目前就是删除所有的preifix为前缀的东西
	return h.deleteKeysByPrefix(ctx)
}

func (h *Manager) deleteKeysByPrefix(ctx context.Context) error {
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
