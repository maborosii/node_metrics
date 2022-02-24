package setting

import (
	"fmt"
	"reflect"
	"testing"
)

func TestConf_GetTimeOut(t *testing.T) {
	InitConfig("../conf", "monitor.toml")
	fmt.Println(*Config.LogConfig)
	tests := []struct {
		name   string
		fields *MonitorConfig
		want   int
	}{
		{name: "config", fields: Config, want: 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &MonitorConfig{
				Address:   tt.fields.Address,
				TimeOut:   tt.fields.TimeOut,
				Items:     tt.fields.Items,
				Output:    tt.fields.Output,
				LogConfig: tt.fields.LogConfig,
			}
			if got := conf.GetTimeOut(); got != tt.want {
				t.Errorf("Conf.GetTimeOut() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConf_GetMonitorItems(t *testing.T) {
	InitConfig("../conf", "monitor.toml")
	tests := []struct {
		name   string
		fields *MonitorConfig
		want   map[string]string
	}{
		{name: "config", fields: Config, want: map[string]string{"1": "1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &MonitorConfig{
				Address:   tt.fields.Address,
				TimeOut:   tt.fields.TimeOut,
				Items:     tt.fields.Items,
				Output:    tt.fields.Output,
				LogConfig: tt.fields.LogConfig,
			}
			if got := conf.GetMonitorItems(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Conf.GetMonitorItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConf_GetOutputTitle(t *testing.T) {
	InitConfig("../conf", "monitor.toml")
	tests := []struct {
		name   string
		fields *MonitorConfig
		want   []string
	}{
		{name: "config", fields: Config, want: []string{"1", "1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &MonitorConfig{
				Address:   tt.fields.Address,
				TimeOut:   tt.fields.TimeOut,
				Items:     tt.fields.Items,
				Output:    tt.fields.Output,
				LogConfig: tt.fields.LogConfig,
			}
			if got := conf.GetOutputTitle(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Conf.GetOutputTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMonitorConfig_GetOutputFileAndSheetName(t *testing.T) {
	InitConfig("../conf", "monitor.toml")
	fmt.Println(*Config)
	tests := []struct {
		name   string
		fields *MonitorConfig
		want   string
		want1  string
	}{
		{name: "config", fields: Config, want: "1", want1: "2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &MonitorConfig{
				Address:   tt.fields.Address,
				TimeOut:   tt.fields.TimeOut,
				Items:     tt.fields.Items,
				Output:    tt.fields.Output,
				LogConfig: tt.fields.LogConfig,
			}
			got, got1 := conf.GetOutputFileAndSheetName()
			if got != tt.want {
				t.Errorf("MonitorConfig.GetOutputFileAndSheetName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MonitorConfig.GetOutputFileAndSheetName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
