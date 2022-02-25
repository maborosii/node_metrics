package main

import (
	"fmt"
	"log"
	"node_metrics_go/cmd"
	"node_metrics_go/global"
	"node_metrics_go/internal/etl"
	"node_metrics_go/internal/excelops"
	"node_metrics_go/pkg/logger"
	"node_metrics_go/pkg/setting"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// func main() {
// 	defer log.Logger.Sync()
// 	test.Print()
// }

// func init() {
// 	// 传入配置文件路径，加载配置文件,
// 	if err := setting.InitConfig("conf", "monitor.toml"); err != nil {
// 		fmt.Printf("load config from file failed, err:%v\n", err)
// 		return
// 	}
// 	fmt.Println("config.ini配置加载成功", setting.Config.GetAddress())

// 	logConfig := setting.Config.GetLogConfig()
// 	log.InitLogger(logConfig)
// 	log.Logger.Debug("大家好，日志展示")
// }

func init() {
	// 获取配置文件
	var err error
	err = cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute err: %v", err)
	}
	fmt.Println(global.ConfigPath)
	// 初始化配置
	err = setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	// 初始化日志
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
}

func main() {

	// fmt.Println(global.MonitorSetting.GetAddress())
	// fmt.Println(global.MonitorSetting.GetTimeOut())
	// fmt.Println(*global.MonitorSetting.GetLogConfig())
	// fmt.Println(global.MonitorSetting.GetMonitorItems())
	// fmt.Println(global.MonitorSetting.GetOutputFileAndSheetName())
	// 存储最终指标
	var storeResults = etl.NewStoreResults()

	// prom客户端api
	queryApi := etl.ClientForProm(global.MonitorSetting.GetAddress())

	// 查询指标
	for label, sql := range global.MonitorSetting.GetMonitorItems() {
		go func(label, sql string, queryApi v1.API) {
			etl.QueryFromProm(label, sql, queryApi)
		}(label, sql, queryApi)
	}
	etl.WgReceiver.Add(1)
	// 转换数据
	go etl.ShuffleResult(len(global.MonitorSetting.GetMonitorItems()), &storeResults)
	etl.WgReceiver.Wait()

	writeResults := [][]string{}
	for _, sr := range storeResults {
		global.Logger.Info("get node of all metrics", zap.String("metrics", sr.Print()))
		writeResults = append(writeResults, sr.ConvertToSlice())
	}
	// 写入数据
	f := excelize.NewFile()
	sheetname, filename := global.MonitorSetting.GetOutputFileAndSheetName()
	index := f.NewSheet(sheetname)

	excelops.WriteExcel(f, sheetname, global.MonitorSetting.GetOutputTitle(), writeResults)
	f.SetActiveSheet(index)
	if err := f.SaveAs(filename); err != nil {
		global.Logger.Fatal("save excel file occur err, err info: ", zap.Error(err))
	}

}

func setupSetting() error {
	setting, err := setting.NewSetting(global.ConfigPath)
	if err != nil {
		return err
	}
	err = setting.ReadConfig(&global.MonitorSetting)
	if err != nil {
		return err
	}

	return nil
}

func setupLogger() error {
	logConfig := global.MonitorSetting.GetLogConfig()
	err := logger.InitLogger(logConfig)
	if err != nil {
		return err
	}
	return nil
}
