package cacheadapter

const (
	// 在mysql数据库里
	defaultDatabaseName = "casbin"
	defaultTableName    = "casbin_rule"

	// 策略再redis中的key
	Policy_Key = "casbin:policy"

	// 前缀，代表他要检测redis的哪一个键
	Casbin_Indicate_Prefix = "casbin:op"

	//往下都是值，也就是需要的操作

	// Casbin_Policy_Add 添加列表末尾的几行进入内存，addplicy命令
	// 像deny就是直接加一条就可以了
	Casbin_Policy_Add = "add:"
	// Casbin_Policy_Del 删除谁的什么权限，不能做就重新加载
	Casbin_Policy_Del = "del:"
	// Casbin_Policy_Load 重新加载
	Casbin_Policy_Load = "load"
)
