package Test

import (
	"fmt"
	"github.com/casbin/casbin/v2"
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

func TestInitPolicy(t *testing.T) {
	Init()
	Init()
	// 这个model初始化估计光这样不行，看下边好像还有啥子
	e, err := casbin.NewEnforcer("D:\\code\\go\\trys\\casbin_demo\\config\\rbac_model.conf", cacheada)
	//fmt.Println("这之前执行过嘛")
	if err != nil {
		log.Fatalf("new Enforcer err ,%v \n", err)
	}
	//err = e.LoadPolicy()
	if err != nil {
		log.Fatalf("load policy is err,%+v\n", err)
		return
	}
	log.Println("能到这里就是没问题，往下就是测试验证看看")
	e.AddPolicy("alice", "data1", "read", "allow")
	e.AddPolicy("alice", "data1", "write", "allow")
	e.AddPolicy("alice", "data1", "cry", "allow")
	e.AddPolicy("alice", "data2", "cry", "allow")
	e.AddPolicy("alice", "data2", "read", "allow")
	e.AddPolicy("alice", "data2", "write", "allow")
	e.AddPolicy("alice", "data3", "cry", "allow")
}

func TestAdd(t *testing.T) {
	Init()
	k := []string{"alice", "data1", "read", "deny"}
	//k := []string{"alice", "data1", "write", "allow"}
	err := cacheada.AddPolicy("", "p", k)
	if err != nil {
		log.Fatalf("err is ,%+v \n", err)
	}
}

func TestLoad(t *testing.T) {
	Init()
	// 这个model初始化估计光这样不行，看下边好像还有啥子
	e, err := casbin.NewEnforcer("D:\\code\\go\\trys\\casbin_demo\\config\\rbac_model.conf", cacheada)
	if err != nil {
		log.Fatalf("new Enforcer err ,%v \n", err)
	}
	err = e.LoadPolicy()
	if err != nil {
		log.Fatalf("load policy is err,%+v\n", err)
		return
	}
	log.Println("能到这里就是没问题，往下就是测试验证看看")
	//e.AddPolicy("alice", "data1", "read", "allow")
	//e.AddPolicy("alice", "data1", "write", "allow")

	enforce, err := e.Enforce("alice", "data1", "read")
	if err != nil {
		log.Fatalf(" 策略检查报错，%+v\n", err)
	}
	log.Printf("检查结果是, %+#v\n", enforce)
}

func TestRemove(t *testing.T) {
	Init()
	// 这个model初始化估计光这样不行，看下边好像还有啥子
	e, err := casbin.NewEnforcer("D:\\code\\go\\trys\\casbin_demo\\config\\rbac_model.conf", cacheada)
	if err != nil {
		log.Fatalf("new Enforcer err ,%v \n", err)
	}
	err = e.LoadPolicy()
	if err != nil {
		log.Fatalf("load policy is err,%+v\n", err)
		return
	}
	log.Println("能到这里就是没问题，往下就是测试验证看看")
	enforce, err := e.Enforce("alice", "data1", "write")
	log.Printf("未删除该规则时的检测结果是, %+#v\n", enforce)

	_, err = e.RemovePolicy("alice", "data1", "write", "allow")
	if err != nil {
		log.Fatalf("removePolicy is errr: %+v \n", err)
	}
	enforce, err = e.Enforce("alice", "data1", "write")
	if enforce {
		log.Fatalf("使用remove之后的检测结果是, %+#v\n", enforce)
	}
	k := []string{"alice", "data1", "write", "allow"}
	_ = cacheada.AddPolicy("", "p", k)

}

// 在初始化的时候会自动进行一次load，
func TestSave(t *testing.T) {

	Init()
	// 这个model初始化估计光这样不行，看下边好像还有啥子
	e, err := casbin.NewEnforcer("D:\\code\\go\\trys\\casbin_demo\\config\\rbac_model.conf", cacheada)
	//fmt.Println("这之前执行过嘛")
	if err != nil {
		log.Fatalf("new Enforcer err ,%v \n", err)
	}
	//err = e.LoadPolicy()
	if err != nil {
		log.Fatalf("load policy is err,%+v\n", err)
		return
	}
	//e.EnableAutoSave(false)

	log.Println("能到这里就是没问题，往下就是测试验证看看")
	_, err = e.AddPolicy("k996", "data1", "read", "allow")
	if err != nil {
		log.Fatalf("addpolicy is err :%+#v\n", err)
	}
	_, err = e.AddPolicy("k96", "data1", "write", "allow")
	if err != nil {
		return
	}
	if err != nil {
		log.Fatalf("addpolicy is err :%+#v\n", err)
	}
	err = e.SavePolicy()
	if err != nil {
		log.Fatalf("err is  %+#v", err)
	}
}

func TestRemoveFilter(t *testing.T) {
	Init()
	// 这个model初始化估计光这样不行，看下边好像还有啥子
	e, err := casbin.NewEnforcer("D:\\code\\go\\trys\\casbin_demo\\config\\rbac_model.conf", cacheada)
	//fmt.Println("这之前执行过嘛")
	if err != nil {
		log.Fatalf("new Enforcer err ,%v \n", err)
	}
	//err = e.LoadPolicy()
	if err != nil {
		log.Fatalf("load policy is err,%+v\n", err)
		return
	}

	log.Println("能到这里就是没问题，往下就是测试验证看看")
	_, err = e.RemoveFilteredPolicy(0, "alice", "data1")
	// fieldIndex 是alice放的位置，也就是v0
	if err != nil {
		log.Fatalf("测试失败 %+v", err)
	}
}
