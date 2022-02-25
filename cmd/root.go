package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{}

// 供主程序调用
func Execute() error {
	return rootCmd.Execute()
}

// 初始化，将子命令嵌入到根命令中
func init() {
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(versionCmd)
}
