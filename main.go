// Copyright 2019 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License"); // you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package v1_test provides examples making requests to Prometheus using the
// Golang client.
package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// 正则表达式匹配模式 --> 筛选出ip和value
var pattern = `(?m)ip="(.*?)".*\s=>\s*(\d*\.\d{2}).*$`

// 保存查询结果, --> 给查询结果打标签
type QueryResult struct {
	label string
	value model.Value
}

func NewQueryResult() func(label string, value model.Value) *QueryResult {
	return func(label string, value model.Value) *QueryResult {
		return &QueryResult{label: label, value: value}
	}
}
func (q *QueryResult) Print() {
	fmt.Printf("label:%s,\ndata:%v\n", q.label, q.value)
}
func (q *QueryResult) GetLabel() string {
	return q.label
}
func (q *QueryResult) GetValue() model.Value {
	return q.value
}

// 抽取ip，label，value
func (q *QueryResult) CleanValue() [][]string {
	var midResult = [][]string{}
	var re = regexp.MustCompile(pattern)
	matched := re.FindAllStringSubmatch(q.value.String(), -1)
	for _, match := range matched {
		midResult = append(midResult, []string{match[1], match[2]})
	}
	return midResult
}

// 结构化查询结果
type StoreResult struct {
	ip             string
	cpuAvg         string
	cpuMax         string
	memAvg         string
	memMax         string
	diskUsage      string
	diskReadAvg    string
	diskReadMax    string
	diskWriteAvg   string
	diskWriteMax   string
	networkIn      string
	networkOut     string
	contextSwitchs string
	socketNums     string
}

// 用于灵活构建StoreResult
type Option func(*StoreResult)

func WithCpuAvg(cpuAvg string) Option {
	return func(sr *StoreResult) {
		sr.cpuAvg = cpuAvg
	}
}
func WithCpuMax(cpuMax string) Option {
	return func(sr *StoreResult) {
		sr.cpuMax = cpuMax
	}
}
func WithMemAvg(memAvg string) Option {
	return func(sr *StoreResult) {
		sr.memAvg = memAvg
	}
}
func WithMemMax(memMax string) Option {
	return func(sr *StoreResult) {
		sr.memMax = memMax
	}
}
func WithDiskUsage(diskUsage string) Option {
	return func(sr *StoreResult) {
		sr.diskUsage = diskUsage
	}
}
func WithDiskReadAvg(diskReadAvg string) Option {
	return func(sr *StoreResult) {
		sr.diskReadAvg = diskReadAvg
	}
}
func WithDiskWriteAvg(diskWriteAvg string) Option {
	return func(sr *StoreResult) {
		sr.diskWriteAvg = diskWriteAvg
	}
}
func WithDiskReadMax(diskReadMax string) Option {
	return func(sr *StoreResult) {
		sr.diskReadMax = diskReadMax
	}
}
func WithDiskWriteMax(diskWriteMax string) Option {
	return func(sr *StoreResult) {
		sr.diskWriteMax = diskWriteMax
	}
}
func WithNetworkIn(networkIn string) Option {
	return func(sr *StoreResult) {
		sr.networkIn = networkIn
	}
}
func WithNetworkOut(networkOut string) Option {
	return func(sr *StoreResult) {
		sr.networkOut = networkOut
	}
}
func WithContextSwitchs(contextSwitchs string) Option {
	return func(sr *StoreResult) {
		sr.contextSwitchs = contextSwitchs
	}
}
func WithSocketNums(socketNums string) Option {
	return func(sr *StoreResult) {
		sr.socketNums = socketNums
	}
}

// func WithCpuAge(cpuAvg string) Option {
// 	return func(sr *StoreResult) {
// 		sr.cpuAvg = cpuAvg
// 	}
// }
func NewStoreResult(ip string, options ...Option) *StoreResult {
	sr := &StoreResult{ip: ip}
	for _, option := range options {
		option(sr)
	}
	return sr
}

func (sr *StoreResult) GetIp() string {
	return sr.ip
}
func (sr *StoreResult) ModifyStoreResult(options ...Option) {
	for _, option := range options {
		option(sr)
	}
}

type StoreResults []*StoreResult

func NewStoreResults() StoreResults {
	return []*StoreResult{}
}
func (srs StoreResults) FindIp(ip string) (bool, int) {
	for i, sr := range srs {
		if sr.GetIp() == ip {
			return true, i
		}
	}
	return false, -1
}

var metricsChan = make(chan *QueryResult)
var notifyChan = make(chan struct{})
var wgReceiver sync.WaitGroup
var storeResults = NewStoreResults()

func ExampleAPI_query(label string, promql string) {
	client, err := api.NewClient(api.Config{
		Address: "http://prometheus:9090",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, promql, time.Now())

	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	metricsChan <- NewQueryResult()(label, result)
	fmt.Println(label)
	<-notifyChan
}
func ShuffleResult(series int) {
	defer wgReceiver.Done()
	// storeResults := NewStoreResults()
	for i := 0; i < series; i++ {
		queryResult, ok := <-metricsChan
		fmt.Println("Receiver:", ok)
		// if !ok {
		// 	// 发送关闭通知到各发送者goroutine
		// 	close(notifyChan)
		// 	return
		// }
		switch queryResult.GetLabel() {
		case "cpu_usage_avg_percents":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithCpuAvg(result[1] + "%"))
				} else {
					sr := NewStoreResult(result[0], WithCpuAvg(result[1]+"%"))
					storeResults = append(storeResults, sr)
				}
			}
		case "cpu_usage_max_percents":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithCpuMax(result[1] + "%"))
				} else {
					sr := NewStoreResult(result[0], WithCpuMax(result[1]+"%"))
					storeResults = append(storeResults, sr)
				}
			}
		case "mem_usage_avg_percents":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithMemAvg(result[1] + "%"))
				} else {
					sr := NewStoreResult(result[0], WithMemAvg(result[1]+"%"))
					storeResults = append(storeResults, sr)
				}
			}
		case "mem_usage_max_percents":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithMemMax(result[1] + "%"))
				} else {
					sr := NewStoreResult(result[0], WithMemMax(result[1]+"%"))
					storeResults = append(storeResults, sr)
				}
			}

		case "rootdir_disk_usage_percents":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithDiskUsage(result[1] + "%"))
				} else {
					sr := NewStoreResult(result[0], WithDiskUsage(result[1]+"%"))
					storeResults = append(storeResults, sr)
				}
			}

		case "disk_read_speed_avg_KB_persec":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithDiskReadAvg(result[1] + "KB/s"))
				} else {
					sr := NewStoreResult(result[0], WithDiskReadAvg(result[1]+"KB/s"))
					storeResults = append(storeResults, sr)
				}
			}
		case "disk_read_speed_max_KB_persec":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithDiskReadMax(result[1] + "KB/s"))
				} else {
					sr := NewStoreResult(result[0], WithDiskReadMax(result[1]+"KB/s"))
					storeResults = append(storeResults, sr)
				}
			}

		case "disk_write_speed_avg_KB_persec":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithDiskWriteAvg(result[1] + "KB/s"))
				} else {
					sr := NewStoreResult(result[0], WithDiskWriteAvg(result[1]+"KB/s"))
					storeResults = append(storeResults, sr)
				}
			}
		case "disk_write_speed_max_KB_persec":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithDiskWriteMax(result[1] + "KB/s"))
				} else {
					sr := NewStoreResult(result[0], WithDiskWriteMax(result[1]+"KB/s"))
					storeResults = append(storeResults, sr)
				}
			}
		case "network_in_speed_MB_persec":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithNetworkIn(result[1] + "MB/s"))
				} else {
					sr := NewStoreResult(result[0], WithNetworkIn(result[1]+"MB/s"))
					storeResults = append(storeResults, sr)
				}
			}
		case "network_out_speed_MB_persec":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithNetworkOut(result[1] + "MB/s"))
				} else {
					sr := NewStoreResult(result[0], WithNetworkOut(result[1]+"MB/s"))
					storeResults = append(storeResults, sr)
				}
			}
		case "context_switches_persec":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithContextSwitchs(result[1]))
				} else {
					sr := NewStoreResult(result[0], WithContextSwitchs(result[1]))
					storeResults = append(storeResults, sr)
				}
			}
		case "socket_nums_K":
			results := queryResult.CleanValue()
			for _, result := range results {
				ok, index := storeResults.FindIp(result[0])
				if ok {
					storeResults[index].ModifyStoreResult(WithSocketNums(result[1] + "K"))
				} else {
					sr := NewStoreResult(result[0], WithSocketNums(result[1]+"K"))
					storeResults = append(storeResults, sr)
				}
			}
		default:
			fmt.Printf("Default")
		}
	}
	close(notifyChan)
}

func main() {
	promqls := map[string]string{
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
		"context_switches_persec":        "rate(node_context_switches_total[5m])",
		"socket_nums_K":                  "node_sockstat_sockets_used/1000",
	}
	for label, sql := range promqls {
		go func(label, sql string) {
			ExampleAPI_query(label, sql)
		}(label, sql)
	}
	wgReceiver.Add(1)
	go ShuffleResult(len(promqls))
	wgReceiver.Wait()
	for _, n := range storeResults {
		fmt.Println(*n)
	}
}
