package etl

import (
	// . "node_metrics_go/pkg/log"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/prometheus/common/model"
)

func TestCleanValue(t *testing.T) {
	type fields struct {
		label string
		value model.Value
	}
	type args struct {
		pattern string
	}

	mockCtl := gomock.NewController(t)
	mockExtracter := NewMockValue(mockCtl)
	mockExtracter.EXPECT().String().Return(`{ip="172.16.57.18",a="2222"} => 53.040350877165196 @[1645422102.96]
	{ip="172.16.57.21"} => 2.336293859621752 @[1645422102.96]
	{ip="172.16.57.22"} => 0 @[1645422102.96]
	{ip="172.16.57.18",a="111"} => 0 @[1645422102.96]`)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   [][]string
	}{
		{
			name:   "normal",
			fields: fields{label: "cpu_max", value: mockExtracter},
			args:   args{pattern: `(?m)ip="(.*?)".*\s=>\s*(\d*\.?\d{0,2}).*$`},
			want: [][]string{
				[]string{"172.16.57.18", "53.04"},
				[]string{"172.16.57.21", "2.33"},
				[]string{"172.16.57.22", "0"},
				[]string{"172.16.57.18", "0"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QueryResult{
				label: tt.fields.label,
				value: tt.fields.value,
			}
			if got := q.CleanValue(tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryResult.CleanValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
