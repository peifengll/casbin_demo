package Test

import (
	"fmt"
	"github.com/casbin/casbin/v2/model"
	"github.com/go-redis/redis/v8"
	"github.com/peifengll/casbin_demo/config"
	"github.com/peifengll/casbin_demo/rbac/cacheadapter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
)

var cacheada *cacheadapter.CacheAdapter

func Init() {
	co := config.GetConfig()
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=true",
		co.Mysql.User,
		co.Mysql.Password,
		co.Mysql.Host,
		co.Mysql.Port,
		co.Mysql.DbName,
	)
	client := redis.NewClient(&redis.Options{
		Addr:     co.Redis.Addr,
		Password: co.Redis.Pass,
		DB:       0,
	})
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db 加载失败, %v", err)
	}
	cacheada, err = cacheadapter.NewAdapter(db, client)
	if err != nil {
		log.Fatalf("cacheadapter  加载失败, %v", err)
	}

}

func TestAdd(t *testing.T) {
	Init()
	k := []string{"kill", "data1", "read", "deny"}
	err := cacheada.AddPolicy("", "p", k)
	if err != nil {
		log.Fatalf("err is ,%+v \n", err)
	}
}

func TestLoad(t *testing.T) {
	Init()
	m := model.Model{}
	err := cacheada.LoadPolicy(m)
	if err != nil {
		log.Fatalf("err")
		return
	}
}
