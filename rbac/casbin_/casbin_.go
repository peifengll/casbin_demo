package casbin_

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/peifengll/casbin_demo/rbac/cacheadapter"
	"github.com/peifengll/casbin_demo/rbac/manage"
	"gorm.io/gorm"
)

type Handler struct {
	enforcer *casbin.Enforcer
	c        *manage.Checker
}

func NewHandler(client *redis.Client, db *gorm.DB, prefix, path string) (*Handler, error) {
	uid := uuid.New()
	cacheada := cacheadapter.NewAdapter(db, client)
	checker := manage.NewChecker(client, uid.String(), prefix)
	e, err := casbin.NewEnforcer(path, cacheada)
	if err != nil {
		return nil, err
	}
	return &Handler{
		enforcer: e,
		c:        checker,
	}, nil
}

func (r *Handler) HasPermission(who Evaluator, what Resource, do Action) (bool, error) {
	// 先进行检查一致性
	err := r.c.CheckConsist(r.enforcer)
	if err != nil {
		return false, err
	}
	val, err := r.enforcer.Enforce(who, what, do)
	if err != nil {
		return false, err
	}
	return val, nil
}
