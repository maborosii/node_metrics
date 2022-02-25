package cmd

import (
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
		// fmt.Println(global.MonitorSetting.GetAddress())

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
