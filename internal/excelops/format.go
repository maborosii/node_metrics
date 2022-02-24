/*
	size： excelize/excel

	列 		 A		 	 B		 	 C	     	  D	 	 	 E
	width 	8.38/10	    22.88/27    42.25/48	 8.63/9.37	14/14.7
	...

	行	height
	1	31
	2	15
	3   16.5
	.	.
	.	.
	.	.
	37	16.5
*/

package excelops

import (
	"fmt"
	. "node_metrics_go/pkg/log"

	"github.com/xuri/excelize/v2"
)

// 初始化单元格格式
func GetStyle(f *excelize.File, styleconfig string) int {
	style, err := f.NewStyle(styleconfig)
	if err != nil {
		Log.Fatal("when init style of cell occur error, error info: ", err)
	}
	return style
}

// 初始化标题格式
func GetTitle(style int, titleContent []string) []*Cell {
	title := []*Cell{}
	startX := 'a'
	for _, v := range titleContent {
		singleTitleCell := NewCell(false, []string{fmt.Sprintf("%c", startX)}, []int{1}, []float64{27}, []float64{31}, WithCellStyle(style), WithCellContent(v))
		title = append(title, singleTitleCell)
		startX++
	}
	return title
}

// 初始化内容格式
func GetContent(style, start_index_y int, dataContent [][]string) [][]*Cell {

	data := [][]*Cell{}
	for rowNum, row := range dataContent {
		startX := 'a'
		singleDataRow := []*Cell{}
		for _, v := range row {
			singleContentCell := NewCell(false, []string{fmt.Sprintf("%c", startX)}, []int{start_index_y + rowNum}, []float64{27}, []float64{31}, WithCellStyle(style), WithCellContent(v))
			singleDataRow = append(singleDataRow, singleContentCell)
			startX++
		}
		data = append(data, singleDataRow)
	}
	return data
}
