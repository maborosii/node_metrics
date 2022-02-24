package setting

import (
	"path"
	"strings"

	"github.com/spf13/viper"
)

var Config = NewMonitorConfig()

type MonitorConfig struct {
	Address   string         `toml:"address"`
	TimeOut   int            `toml:"timeout"`
	Items     MonitorItems   `toml:"items"`
	Output    *MonitorOutput `toml:"output"`
	LogConfig *LogConf       `toml:"logconfig"`
}
type LogConf struct {
	Level      string `toml:"level"`
	LogFile    string `toml:"logfile"`
	MaxSize    int    `toml:"maxsize"`
	MaxAge     int    `toml:"maxage"`
	MaxBackups int    `toml:"maxbackups"`
}

func (conf *MonitorConfig) GetTimeOut() int {
	return conf.TimeOut
}
func (conf *MonitorConfig) GetAddress() string {
	return conf.Address
}
func (conf *MonitorConfig) GetMonitorItems() map[string]string {
	return conf.Items.ConvertToMap()
}
func (conf *MonitorConfig) GetOutputFileAndSheetName() (string, string) {
	return conf.Output.Project + "_" + conf.Output.FileName, conf.Output.SheetName
}
func (conf *MonitorConfig) GetOutputTitle() []string {
	return conf.Output.Title
}

// func (conf *MonitorConfig) GetLogConfig() *LogConf {
// 	return conf.LogConfig
// }

func NewMonitorConfig() *MonitorConfig {
	return &MonitorConfig{}
}

type MonitorItems []*MonitorItem

type MonitorItem struct {
	Metrics string `toml:"metrics"`
	Promql  string `toml:"promql"`
}

func (i MonitorItems) ConvertToMap() map[string]string {
	promsqls := make(map[string]string)
	for _, item := range i {
		promsqls[item.Metrics] = item.Promql
	}
	return promsqls
}

type MonitorOutput struct {
	Project   string   `toml:"project"`
	FileName  string   `toml:"filename"`
	SheetName string   `toml:"sheetname"`
	Title     []string `toml:"title"`
}

// 将配置文件映射到全局变量Config中
func InitConfig(filePath string, fileName string) error {
	viper.AddConfigPath(filePath)
	viper.SetConfigName(fileName)
	viper.SetConfigType(strings.TrimLeft(path.Ext(fileName), "."))

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(Config)
	if err != nil {
		return err
	} // // 监听配置文件更新
	// viper.WatchConfig()
	// viper.OnConfigChange(func(e fsnotify.Event) {
	// 	viper.Unmarshal(conf)
	// })

	return nil
}
