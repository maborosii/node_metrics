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

	"node_metrics/excelops/excelwriting"
	. "node_metrics/locallog"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/xuri/excelize/v2"
)

// 正则表达式匹配模式 --> 筛选出ip和value
var pattern = `(?m)ip="(.*?)".*\s=>\s*(\d*\.?\d{0,2}).*$`

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
func (sr *StoreResult) Print() {
	fmt.Printf("\n#######\nip: %s\ncpuAvg: %s\ncpuMax: %s\nmemAvg: %s\nmemMax: %s\ndiskUsage: %s\ndiskReadAvg: %s\ndiskReadMax: %s\ndiskWriteAvg: %s\ndiskWriteMax: %s\nnetworkIn: %s\nnetworkOut: %s\ncontextSwitchs: %s\nsocketNums: %s\n", sr.ip, sr.cpuAvg, sr.cpuMax, sr.memAvg, sr.memMax, sr.diskUsage, sr.diskReadAvg, sr.diskReadMax, sr.diskWriteAvg, sr.diskWriteMax, sr.networkIn, sr.networkOut, sr.contextSwitchs, sr.socketNums)
}
func (sr *StoreResult) ModifyStoreResult(options ...Option) {
	for _, option := range options {
		option(sr)
	}
}
func (sr *StoreResult) ConvertToSlice() []string {
	// 这里可以使用反射，但为了保证生成的slice的有序性
	return []string{
		sr.ip,
		sr.cpuAvg,
		sr.cpuMax,
		sr.memAvg,
		sr.memMax,
		sr.diskUsage,
		sr.diskReadAvg,
		sr.diskReadMax,
		sr.diskWriteAvg,
		sr.diskWriteMax,
		sr.networkIn,
		sr.networkOut,
		sr.contextSwitchs,
		sr.socketNums,
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

func CreateOrModifyStorResults(ip string, srs StoreResults, option Option) {
	ok, index := storeResults.FindIp(ip)
	if ok {
		storeResults[index].ModifyStoreResult(option)
	} else {
		sr := NewStoreResult(ip, option)
		storeResults = append(storeResults, sr)
	}
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
		results := queryResult.CleanValue()
		for _, result := range results {
			switch queryResult.GetLabel() {
			case "cpu_usage_avg_percents":
				CreateOrModifyStorResults(result[0], storeResults, WithCpuAvg(result[1]+"%"))
			case "cpu_usage_max_percents":
				CreateOrModifyStorResults(result[0], storeResults, WithCpuMax(result[1]+"%"))
			case "mem_usage_avg_percents":
				CreateOrModifyStorResults(result[0], storeResults, WithMemAvg(result[1]+"%"))
			case "mem_usage_max_percents":
				CreateOrModifyStorResults(result[0], storeResults, WithMemMax(result[1]+"%"))
			case "rootdir_disk_usage_percents":
				CreateOrModifyStorResults(result[0], storeResults, WithDiskUsage(result[1]+"%"))
			case "disk_read_speed_avg_KB_persec":
				CreateOrModifyStorResults(result[0], storeResults, WithDiskReadAvg(result[1]+"KB/s"))
			case "disk_read_speed_max_KB_persec":
				CreateOrModifyStorResults(result[0], storeResults, WithDiskReadMax(result[1]+"KB/s"))
			case "disk_write_speed_avg_KB_persec":
				CreateOrModifyStorResults(result[0], storeResults, WithDiskWriteAvg(result[1]+"KB/s"))
			case "disk_write_speed_max_KB_persec":
				CreateOrModifyStorResults(result[0], storeResults, WithDiskWriteMax(result[1]+"KB/s"))
			case "network_in_speed_MB_persec":
				CreateOrModifyStorResults(result[0], storeResults, WithNetworkIn(result[1]+"MB/s"))
			case "network_out_speed_MB_persec":
				CreateOrModifyStorResults(result[0], storeResults, WithNetworkOut(result[1]+"MB/s"))
			case "context_switches_persec_K":
				CreateOrModifyStorResults(result[0], storeResults, WithContextSwitchs(result[1]+"K"))
			case "socket_nums_K":
				CreateOrModifyStorResults(result[0], storeResults, WithSocketNums(result[1]+"K"))
			default:
				fmt.Printf("Default")
			}
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
		"context_switches_persec_K":      "rate(node_context_switches_total[5m])/1000",
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
		n.Print()
	}
	sheetname := "Sheet1"
	savexlsx := "巡检报告.xlsx"

	f := excelize.NewFile()
	index := f.NewSheet(sheetname)
	writeResults := [][]string{}
	for _, sr := range storeResults {
		writeResults = append(writeResults, sr.ConvertToSlice())
	}
	excelwriting.WriteExcel(f, sheetname, writeResults)
	f.SetActiveSheet(index)
	if err := f.SaveAs(savexlsx); err != nil {
		Log.Fatal(err)
	}
}
