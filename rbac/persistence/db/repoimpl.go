package db

import (
	"context"
	"github.com/peifengll/casbin_demo/rbac/convterter"
	"github.com/peifengll/casbin_demo/rbac/persistence/po"
	"gorm.io/gorm"
)

type CasbinRuleDaoImpl struct {
	db *gorm.DB
}

var _ CasbinRuleDao = &CasbinRuleDaoImpl{}

// QuerryAll 最好别用
func (r *CasbinRuleDaoImpl) QuerryAll(ctx context.Context) ([]po.CasbinRule, error) {
	lines := make([]po.CasbinRule, 0)
	err := r.db.Find(&lines).Error
	if err != nil {
		return nil, err
	}
	return lines, err
}
func (r *CasbinRuleDaoImpl) QueryRuleById(ctx context.Context, id int) (*po.CasbinRule, error) {
	line := &po.CasbinRule{}
	if err := r.db.First(&line, id).Error; err != nil {
		return nil, err
	}
	return line, nil
}
func (r *CasbinRuleDaoImpl) QueryRulesByIds(ctx context.Context, ids []int) ([]po.CasbinRule, error) {
	lines := make([]po.CasbinRule, 0)
	if err := r.db.Find(&lines, ids).Error; err != nil {
		return nil, err
	}
	return lines, nil
}
func (r *CasbinRuleDaoImpl) QueryRuleByRule(ctx context.Context, ptype string, rule []string) ([]po.CasbinRule, error) {
	line := convterter.ToPolicyPo(ptype, rule)
	lines := make([]po.CasbinRule, 0)
	if err := r.db.Model(&po.CasbinRule{}).Where(&line).Find(&lines).Error; err != nil {
		return nil, err
	}
	return lines, nil
}

func (r *CasbinRuleDaoImpl) UpdateRuleById(ctx context.Context, id int, ptype string, rule []string) error {
	newline := convterter.ToPolicyPo(ptype, rule)
	return r.db.Model(&po.CasbinRule{}).Where("id = ?", id).Updates(newline).Error
}
func (r *CasbinRuleDaoImpl) UpdateRulesByIds(ctx context.Context, ids []int, ptype string, rule []string) error {
	newline := convterter.ToPolicyPo(ptype, rule)
	return r.db.Model(&po.CasbinRule{}).Where("id in (?)", ids).Updates(newline).Error
}

func (r *CasbinRuleDaoImpl) UpdateRulesByRule(ctx context.Context,
	ptype string, oldRule, newPolicy []string) error {
	oldline := convterter.ToPolicyPo(ptype, oldRule)
	newline := convterter.ToPolicyPo(ptype, newPolicy)
	return r.db.Model(&po.CasbinRule{}).Where(&oldline).Updates(newline).Error
}

func (r *CasbinRuleDaoImpl) AddRule(ctx context.Context, pytepe string, rule []string) error {
	line := convterter.ToPolicyPo(pytepe, rule)
	return r.db.Create(&line).Error
}
func (r *CasbinRuleDaoImpl) AddRules(ctx context.Context, ptype string, rules [][]string) error {
	var lines []po.CasbinRule
	for _, rule := range rules {
		line := convterter.ToPolicyPo(ptype, rule)
		lines = append(lines, line)
	}
	return r.db.Create(&lines).Error
}

func (r *CasbinRuleDaoImpl) RemoveRuleById(ctx context.Context, id int) error {
	return r.db.Delete(&po.CasbinRule{}, id).Error
}
func (r *CasbinRuleDaoImpl) RemoveRuleByIds(ctx context.Context, ids []int) error {
	return r.db.Delete(&po.CasbinRule{}, ids).Error
}

func (r *CasbinRuleDaoImpl) RemoveRuleByRule(ctx context.Context, ptype string, rule []string) error {
	line := convterter.ToPolicyPo(ptype, rule)
	return r.rawDelete(r.db, line)
}

func (r *CasbinRuleDaoImpl) RemoveRuleByRules(ctx context.Context, ptype string, rules [][]string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, rule := range rules {
			line := convterter.ToPolicyPo(ptype, rule)
			if err := r.rawDelete(tx, line); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *CasbinRuleDaoImpl) rawDelete(db *gorm.DB, line po.CasbinRule) error {
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
