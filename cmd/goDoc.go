package cmd

import (
	"github.com/fatih/color"
	"github.com/mojocn/felix/util"
	"github.com/spf13/cobra"
	"time"
)

// godocCmd represents the godoc command
var godocCmd = &cobra.Command{
	Use:   "godoc",
	Short: "golang.google.cn/pkg",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		color.Blue("访问Go中国官方网站 https://golang.google.cn/pkg/")
		time.Sleep(time.Second * 1)
		util.BrowserOpen("https://golang.google.cn/pkg/")
	},
}

func init() {
	rootCmd.AddCommand(godocCmd)
}
