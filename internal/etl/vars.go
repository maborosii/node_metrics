package etl

import (
	"sync"
)

var metricsChan = make(chan *QueryResult)
var notifyChan = make(chan struct{})
var WgReceiver sync.WaitGroup

// 正则表达式匹配模式 --> 筛选出ip和value
var Pattern = `(?m)ip="(.*?)".*\s=>\s*(\d*\.?\d{0,2}).*$`
