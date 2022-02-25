package excelops

import (
	"node_metrics_go/global"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type Cell struct {
	// Xfile   *excelize.File
	IsMerge bool
	Xzone   []string
	Yzone   []int
	Xwidth  []float64
	Yheight []float64
	// Format  *excelize.Style
	Format  int
	Content string
}

type Write2Xlsx interface {
	SetWidth(sheetname string, f *excelize.File)
	SetHeight(sheetname string, f *excelize.File)
	MergeCell(sheetname string, f *excelize.File)
	SetFormat(sheetname string, f *excelize.File)
	SetValue(sheetname string, f *excelize.File)
}
type Option func(*Cell)

func WithCellStyle(style int) Option {
	return func(c *Cell) {
		c.Format = style
	}
}
func WithCellContent(content string) Option {
	return func(c *Cell) {
		c.Content = content
	}
}

func NewCell(ismerge bool, xzone []string, yzone []int, xwidth, yheight []float64, options ...Option) *Cell {
	cell := &Cell{IsMerge: ismerge, Xzone: xzone, Yzone: yzone, Xwidth: xwidth, Yheight: yheight}
	for _, option := range options {
		option(cell)
	}
	return cell
}

func (cl *Cell) SetWidth(sheetname string, f *excelize.File) {
	for i, strX := range cl.Xzone {
		if err := f.SetColWidth(sheetname, strX, strX, cl.Xwidth[i]); err != nil {
			global.Logger.Fatal("when setting width of cell occur error, error info: ", zap.Error(err))
		}
	}
}

func (cl *Cell) SetHeight(sheetname string, f *excelize.File) {
	for i, intY := range cl.Yzone {
		if err := f.SetRowHeight(sheetname, intY, cl.Yheight[i]); err != nil {
			global.Logger.Fatal("when setting height of cell occur error, error info: ", zap.Error(err))
		}
	}
}

func (cl *Cell) MergeCell(sheetname string, f *excelize.File) {
	if cl.IsMerge {
		if err := f.MergeCell(sheetname, cl.Xzone[0]+strconv.Itoa(cl.Yzone[0]),
			cl.Xzone[len(cl.Xzone)-1]+strconv.Itoa(cl.Yzone[len(cl.Yzone)-1])); err != nil {
			global.Logger.Fatal("when setting merge of cell occur error, error info: ", zap.Error(err))
		}
	}
}

func (cl *Cell) SetFormat(sheetname string, f *excelize.File) {
	err := f.SetCellStyle(sheetname, cl.Xzone[0]+strconv.Itoa(cl.Yzone[0]),
		cl.Xzone[len(cl.Xzone)-1]+strconv.Itoa(cl.Yzone[len(cl.Yzone)-1]), cl.Format)
	if err != nil {
		global.Logger.Fatal("when setting format of cell occur error, error info: ", zap.Error(err))
	}
}
func (cl *Cell) SetValue(sheetname string, f *excelize.File) {
	err := f.SetCellValue(sheetname, cl.Xzone[0]+strconv.Itoa(cl.Yzone[0]), cl.Content)
	if err != nil {
		global.Logger.Fatal("when setting format of cell occur error, error info: ", zap.Error(err))
	}
}

func Formatting(w2x Write2Xlsx, sheetname string, f *excelize.File) {
	w2x.SetWidth(sheetname, f)
	w2x.SetHeight(sheetname, f)
	w2x.MergeCell(sheetname, f)
	w2x.SetFormat(sheetname, f)
}

func Writing(w2x Write2Xlsx, sheetname string, f *excelize.File) {

	w2x.SetValue(sheetname, f)
}
