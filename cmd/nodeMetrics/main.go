package main

import (
	"node_metrics_go/global"
	"node_metrics_go/internal/etl"
	"node_metrics_go/internal/excelops"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

func main() {

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
