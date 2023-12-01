package db

import (
	"context"
	"github.com/peifengll/casbin_demo/rbac/persistence/po"
	"gorm.io/gorm"
)

type CasbinRuleDao interface {
	QuerryAll(ctx context.Context) ([]po.CasbinRule, error)
	QueryRuleById(ctx context.Context, id int) (*po.CasbinRule, error)
	QueryRulesByIds(ctx context.Context, ids []int) ([]po.CasbinRule, error)
	QueryRuleByRule(ctx context.Context, ptype string, rule []string) ([]po.CasbinRule, error)
	//QueryRulesByPType(ctx context.Context, ptype string) ([]po.CasbinRule, error)
	//QueryRulesByEvaluator(ctx context.Context, evaluator string) ([]po.CasbinRule, error)
	//QueryRulesByResource(ctx context.Context, resource string) ([]po.CasbinRule, error)
	//QueryRulesByAction(ctx context.Context, action string) ([]po.CasbinRule, error)

	//	修改操作
	UpdateRuleById(ctx context.Context, id int, ptype string, ru []string) error
	UpdateRulesByIds(ctx context.Context, ids []int, ptype string, ru []string) error
	UpdateRulesByRule(ctx context.Context, ptype string, ru []string, rule []string) error
	// 增加  增加直接用 casbin的增加即可，批量增加？也可以
	AddRule(ctx context.Context, pytepe string, rule []string) error
	AddRules(ctx context.Context, ptype string, rules [][]string) error
	// 删除操作，提供按照id删，和按照规则删
	RemoveRuleById(ctx context.Context, id int) error
	RemoveRuleByIds(ctx context.Context, ids []int) error
	RemoveRuleByRule(ctx context.Context, ptype string, rule []string) error
	RemoveRuleByRules(ctx context.Context, ptype string, rules [][]string) error
}

func NewCasbinRuleDao(db *gorm.DB) *CasbinRuleDaoImpl {
	return &CasbinRuleDaoImpl{db}
}
