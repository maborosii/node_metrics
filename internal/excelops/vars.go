package excelops

import _ "embed"

var (
	// SaveXlsx  = "巡检报告.xlsx"
	// Sheetname = "Sheet1"

	//go:embed styles/title.json
	titleStyle string

	//go:embed styles/content.json
	contentStyle string

	// 标题内容
	// titleContent = []string{
	// 	"服务器IP",
	// 	"CPU（使用率均值）",
	// 	"CPU（使用率峰值）",
	// 	"内存（使用率均值）",
	// 	"内存（使用率峰值）",
	// 	"磁盘（使用率）",
	// 	"磁盘I/O（读速率均值）",
	// 	"磁盘I/O（读速率峰值）",
	// 	"磁盘I/O（写速率均值）",
	// 	"磁盘I/O（写速率峰值）",
	// 	"网络带宽（下行均值）",
	// 	"网络带宽（上行均值）",
	// 	"每秒上下文切换次数",
	// 	"socket使用量",
	// }
)
