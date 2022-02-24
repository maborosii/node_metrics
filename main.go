package main

import (
	"fmt"
	"node_metrics_go/internal/test"
	"node_metrics_go/pkg/log"
	"node_metrics_go/setting"
)

func main() {
	defer log.Logger.Sync()
	test.Print()
}

func init() {
	// 传入配置文件路径，加载配置文件,
	if err := setting.InitConfig("conf", "monitor.toml"); err != nil {
		fmt.Printf("load config from file failed, err:%v\n", err)
		return
	}
	fmt.Println("config.ini配置加载成功", setting.Config.GetAddress())

	logConfig := setting.Config.GetLogConfig()
	log.InitLogger(logConfig)
	log.Logger.Debug("大家好，日志展示")
}
