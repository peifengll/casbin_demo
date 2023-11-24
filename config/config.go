package config

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

var (
	c    *TomConfig
	once sync.Once
)

type TomConfig struct {
	AppName string
	Mysql   MySQLConfig
	Redis   RedisConfig
}

type MySQLConfig struct {
	Host        string
	DbName      string
	Password    string
	Port        int
	TablePrefix string
	User        string
}

type RedisConfig struct {
	Addr string
	Type string
	Pass string
}

func InitConfig() error {
	viper.SetConfigName("conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath("config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")
	err := viper.ReadInConfig()
	if nil != err {
		return err
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		return err
	}
	return nil
}
func GetConfig() *TomConfig {
	once.Do(func() {
		err := InitConfig()
		if err != nil {
			log.Fatalf("viper.Unmarshal err: %v", err)
			return
		}

	})
	return c
}
