package excelops

import (
	_ "embed"

	"github.com/xuri/excelize/v2"
)

// var (
// 	//go:embed styles/title.json
// 	titleStyle string

// 	//go:embed styles/content.json
// 	contentStyle string

// 	// 标题内容
// 	titleContent = []string{
// 		"服务器IP",
// 		"CPU（使用率均值）",
// 		"CPU（使用率峰值）",
// 		"内存（使用率均值）",
// 		"内存（使用率峰值）",
// 		"磁盘（使用率）",
// 		"磁盘I/O（读速率均值）",
// 		"磁盘I/O（读速率峰值）",
// 		"磁盘I/O（写速率均值）",
// 		"磁盘I/O（写速率峰值）",
// 		"网络带宽（下行均值）",
// 		"网络带宽（上行均值）",
// 		"每秒上下文切换次数",
// 		"socket使用量",
// 	}
// )

func WriteExcel(f *excelize.File, sheetname string, data [][]string) {

	styleTitle := GetStyle(f, titleStyle)
	styleContent := GetStyle(f, contentStyle)

	// 写入标题
	for _, cell := range GetTitle(styleTitle, titleContent) {
		Formatting(cell, sheetname, f)
		Writing(cell, sheetname, f)
	}
	// 写入内容
	for _, row := range GetContent(styleContent, 2, data) {
		for _, cell := range row {
			// 需要动态插入值的单元格做空值写入
			Formatting(cell, sheetname, f)
			Writing(cell, sheetname, f)
		}
	}
}
