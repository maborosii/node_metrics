package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/q191201771/naza/pkg/bininfo"
)

func GetVersion() {
	// 添加编译信息
	v := flag.Bool("v", false, "show bin info")
	flag.Parse()
	if *v {
		_, _ = fmt.Fprint(os.Stderr, bininfo.StringifyMultiLine())
		os.Exit(1)
	}
}
