package main

import (
	"flag"
	"fmt"
	config "node_metrics_go/conf"
	"node_metrics_go/internal/etl"
	"node_metrics_go/internal/excelops"
	. "node_metrics_go/pkg/log"
	"os"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/q191201771/naza/pkg/bininfo"
	"github.com/xuri/excelize/v2"
)

func main() {

	config.GetVersion()

	// 存储最终指标
	var storeResults = etl.NewStoreResults()

	// prom客户端api
	queryApi := etl.ClientForProm(etl.Address)

	// 查询指标
	for label, sql := range etl.Promqls {
		go func(label, sql string, queryApi v1.API) {
			etl.QueryFromProm(label, sql, queryApi)
		}(label, sql, queryApi)
	}
	etl.WgReceiver.Add(1)
	// 转换数据
	go etl.ShuffleResult(len(etl.Promqls), &storeResults)
	etl.WgReceiver.Wait()

	writeResults := [][]string{}
	for _, sr := range storeResults {
		Log.Info("node of all metrics: ", sr.Print())
		writeResults = append(writeResults, sr.ConvertToSlice())
	}

	// 写入数据
	f := excelize.NewFile()
	index := f.NewSheet(excelops.Sheetname)

	excelops.WriteExcel(f, excelops.Sheetname, writeResults)
	f.SetActiveSheet(index)
	if err := f.SaveAs(excelops.SaveXlsx); err != nil {
		Log.Fatal("save excel file occur err, err info: ", err)
	}

	// 添加编译信息
	v := flag.Bool("v", false, "show bin info")
	flag.Parse()
	if *v {
		_, _ = fmt.Fprint(os.Stderr, bininfo.StringifyMultiLine())
		os.Exit(1)
	}
}
