package cmd

import (
	"fmt"
	"node_metrics_go/global"
	"strings"

	"github.com/spf13/cobra"
)

var nodeConfig string

// 定义命令行的主要参数
var nodeCmd = &cobra.Command{
	Use:   "node",   // 子命令的标识
	Short: "获取主机指标", // 简短帮助说明
	Long:  nodeDesc, // 详细帮助说明
	Run: func(cmd *cobra.Command, args []string) {
		// 主程序，获取自定义配置文件

		global.ConfigPath = nodeConfig
		fmt.Println(global.MonitorSetting.GetAddress())
		// 存储最终指标
		// var storeResults = etl.NewStoreResults()

		// // prom客户端api
		// queryApi := etl.ClientForProm(global.MonitorSetting.GetAddress())

		// // 查询指标
		// for label, sql := range global.MonitorSetting.GetMonitorItems() {
		// 	go func(label, sql string, queryApi v1.API) {
		// 		etl.QueryFromProm(label, sql, queryApi)
		// 	}(label, sql, queryApi)
		// }
		// etl.WgReceiver.Add(1)
		// // 转换数据
		// go etl.ShuffleResult(len(global.MonitorSetting.GetMonitorItems()), &storeResults)
		// etl.WgReceiver.Wait()

		// writeResults := [][]string{}
		// for _, sr := range storeResults {
		// 	global.Logger.Info("get node of all metrics", zap.String("metrics", sr.Print()))
		// 	writeResults = append(writeResults, sr.ConvertToSlice())
		// }
		// // 写入数据
		// f := excelize.NewFile()
		// sheetname, filename := global.MonitorSetting.GetOutputFileAndSheetName()
		// index := f.NewSheet(sheetname)

		// excelops.WriteExcel(f, sheetname, global.MonitorSetting.GetOutputTitle(), writeResults)
		// f.SetActiveSheet(index)
		// if err := f.SaveAs(filename); err != nil {
		// 	global.Logger.Fatal("save excel file occur err, err info: ", zap.Error(err))
		// }
	},
}
var nodeDesc = strings.Join([]string{
	"该子命令支持获取主机指标，流程如下：",
	"1：从prometues获取节点指标",
	"2：将获取的指标进行转换",
	"3：输出为 地市_巡检报告.xlsx",
}, "\n")

// 用于执行main函数前初始化这个源文件里的变量
func init() {
	// 绑定命令行输入，绑定一个参数
	// 参数分别表示，绑定的变量，参数长名(--str)，参数短名(-s)，默认内容，帮助信息
	nodeCmd.Flags().StringVarP(&nodeConfig, "config", "c", "configs", "请选择配置文件")
}
