package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"mxshop-api/user-web/global"
)

var pathConfig = "user-web/config.yaml"

// Viper 初始化管理配置
func Viper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile(pathConfig)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.ServerConfig); err != nil {
			panic(fmt.Errorf("unmarshal to Conf failed, err:%v", err))
		}
	})
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(fmt.Errorf("unmarshal to Conf failed, err:%v", err))
	}
	return v
}
