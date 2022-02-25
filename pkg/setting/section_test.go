package setting

import "testing"

func TestMonitorConfig_GetTimeOut(t *testing.T) {
	type fields struct {
		Address   string
		TimeOut   int
		Items     MonitorItems
		Output    *MonitorOutput
		LogConfig *LogConf
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
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
				t.Errorf("MonitorConfig.GetTimeOut() = %v, want %v", got, tt.want)
			}
		})
	}
}
