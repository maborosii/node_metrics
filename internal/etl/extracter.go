package etl

import (
	"context"
	"regexp"
	"time"

	. "node_metrics_go/pkg/log"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type QueryResult struct {
	label string
	value model.Value
}

func NewQueryResult() func(label string, value model.Value) *QueryResult {
	return func(label string, value model.Value) *QueryResult {
		return &QueryResult{label: label, value: value}
	}
}

// func (q *QueryResult) Print() {
// 	fmt.Printf("label:%s,\ndata:%v\n", q.label, q.value)
// }
func (q *QueryResult) GetLabel() string {
	return q.label
}
func (q *QueryResult) GetValue() model.Value {
	return q.value
}

// patter: 正则表达式
// 抽取ip，label，metrics
func (q *QueryResult) CleanValue(pattern string) [][]string {
	var midResult = [][]string{}
	var re = regexp.MustCompile(pattern)
	matched := re.FindAllStringSubmatch(q.value.String(), -1)
	for _, match := range matched {
		midResult = append(midResult, []string{match[1], match[2]})
	}
	return midResult
}

func ClientForProm(address string) v1.API {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		Log.Fatal("Error creating client: %v\n", err)
	}
	v1api := v1.NewAPI(client)
	return v1api
}

func QueryFromProm(label string, promql string, api v1.API) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	result, warnings, err := api.Query(ctx, promql, time.Now())

	if err != nil {
		Log.Fatal("Error querying Prometheus: %v\n", err)
	}
	if len(warnings) > 0 {
		Log.Warning("Warnings: %v\n", warnings)
	}
	metricsChan <- NewQueryResult()(label, result)
	Log.Info("label: ", label, " has been gotten")
	<-notifyChan
}
