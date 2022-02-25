// 定义全局变量
package global

import (
	"node_metrics_go/pkg/setting"

	"go.uber.org/zap"
)

var (
	MonitorSetting *setting.MonitorConfig
	Logger         *zap.Logger
	ConfigPath     string
)
