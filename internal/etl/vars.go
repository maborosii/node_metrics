package etl

import (
	"sync"
	"time"
)

var metricsChan = make(chan *QueryResult)
var notifyChan = make(chan struct{})
var WgReceiver sync.WaitGroup

// var ListStoreResult = NewStoreResults()

// 正则表达式匹配模式 --> 筛选出ip和value
var Pattern = `(?m)ip="(.*?)".*\s=>\s*(\d*\.?\d{0,2}).*$`

var Address = "http://prometheus:9090"
var timeOut = 20 * time.Second
var Promqls = map[string]string{
	"cpu_usage_avg_percents":         "(1-avg(rate(node_cpu_seconds_total{mode=\"idle\"}[5m]))by(ip))*100",
	"cpu_usage_max_percents":         "(1-min_over_time(avg(rate(node_cpu_seconds_total{mode=\"idle\"}[5m]))by(ip)[24h:1s]))*100",
	"mem_usage_avg_percents":         "(1-avg_over_time(node_memory_MemAvailable_bytes[24h])/node_memory_MemTotal_bytes)*100",
	"mem_usage_max_percents":         "(1-min_over_time(node_memory_MemAvailable_bytes[24h])/node_memory_MemTotal_bytes)*100",
	"rootdir_disk_usage_percents":    "(1-node_filesystem_free_bytes{mountpoint=\"/\",fstype=~\"xfs|ext4\"}/node_filesystem_size_bytes{mountpoint=\"/\",fstype=~\"xfs|ext4\"})*100",
	"disk_read_speed_avg_KB_persec":  "avg_over_time(rate(node_disk_read_bytes_total{device=\"vdb\"}[5m])[24h:1s])/1024",
	"disk_read_speed_max_KB_persec":  "max_over_time(rate(node_disk_read_bytes_total{device=\"vdb\"}[5m])[24h:1s])/1024",
	"disk_write_speed_avg_KB_persec": "avg_over_time(rate(node_disk_written_bytes_total{device=\"vdb\"}[5m])[24h:1s])/1024",
	"disk_write_speed_max_KB_persec": "max_over_time(rate(node_disk_written_bytes_total{device=\"vdb\"}[5m])[24h:1s])/1024",
	"network_in_speed_MB_persec":     "avg_over_time(rate(node_network_receive_bytes_total{device=\"eth0\"}[5m])[24h:1s])/1024/1024",
	"network_out_speed_MB_persec":    "avg_over_time(rate(node_network_transmit_bytes_total{device=\"eth0\"}[5m])[24h:1s])/1024/1024",
	"context_switches_persec_K":      "rate(node_context_switches_total[5m])/1000",
	"socket_nums_K":                  "node_sockstat_sockets_used/1000",
}
