package cacheadapter

import "github.com/prometheus/client_golang/prometheus"

const (
	// 在mysql数据库里
	defaultDatabaseName = "casbin"
	defaultTableName    = "casbin_rule"

	// 策略再redis中的key
	Policy_Key = "casbin:policy"

	// 前缀，代表他要检测redis的哪一个键
	Casbin_Indicate_Prefix = "casbin:op"

	//往下都是值，也就是需要的操作

	// Casbin_Policy_Ad
	//
	//d 添加列表末尾的几行进入内存，addplicy命令
	// 像deny就是直接加一条就可以了
	Casbin_Policy_Add = "add:"
	// Casbin_Policy_Del 删除谁的什么权限，不能做就重新加载
	Casbin_Policy_Del = "del:"
	// Casbin_Policy_Load 重新加载
	Casbin_Policy_Load = "load"
)

const NS = "PolicyLoad"

const createtablesql = `
CREATE TABLE IF NOT EXISTS  casbin_rule 
(
     id     bigint unsigned NOT NULL AUTO_INCREMENT,
     ptype  varchar(100) DEFAULT NULL,
     v0     varchar(40)  DEFAULT NULL,
     v1     varchar(40)  DEFAULT NULL,
     v2     varchar(40)  DEFAULT NULL,
     v3     varchar(40)  DEFAULT NULL,
     v4     varchar(40)  DEFAULT NULL,
     v5     varchar(40)  DEFAULT NULL,
    PRIMARY KEY ( id ),
    UNIQUE KEY  idx_casbin_rule  ( ptype ,  v0 ,  v1 ,  v2 ,  v3 ,  v4 ,  v5 ),
    UNIQUE KEY  unique_index  ( v0 ,  v1 ,  v2 ,  v3 ,  v4 ,  v5 )
) ENGINE = InnoDB
  AUTO_INCREMENT = 20014
  DEFAULT CHARSET = utf8mb4;
 `

var (
	LoadFormCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: NS,
			Name:      "policy_load_total",
			Help:      "A counter for policy load from mysql and redis.",
		},
		[]string{"from", "status"},
	)
	LoadTimeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   NS,
			Subsystem:   "",
			Name:        "load_seconds",
			Help:        "Histogram of loading policy response latency in seconds.",
			ConstLabels: nil,
			Buckets:     []float64{0.05, 0.1, 0.15, 0.2, 0.25, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 2.0, 5.0},
		}, []string{"from"})
)

func GetCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		LoadFormCounter,
		LoadTimeHistogram,
	}
}
