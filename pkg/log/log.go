package log

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var logPath string = "node_metrics.log"

func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	Log.SetReportCaller(true)
	return Log
}
func init() {
	Log := NewLogger()
	Log.SetOutput(os.Stdout)

	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writers := []io.Writer{
		file,
		os.Stdout}

	//同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		Log.SetOutput(fileAndStdoutWriter)
	} else {
		Log.Info("failed to log to file.")
	}
	Log.SetLevel(logrus.InfoLevel)
}
