package cmd

import (
	"context"
	"os/exec"
	"time"

	"github.com/libragen/felix/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tmSpiderSegmentFaultCmd represents the taskrm command
var tmSpiderSegmentFaultCmd = &cobra.Command{
	Use:   "tmSpider",
	Short: "tech.mojotv.cn:spider segmentFault article && jekyll serve",
	Long:  `只能在我家里面的电脑使用`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jekyllDir := viper.GetString("tech_mojotv_cn.srcDir")

		err := util.ParseUrlPage(args[0], "div.article__content", jekyllDir)
		if err != nil {
			logrus.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
		defer cancel()

		commd := exec.CommandContext(ctx, "bundle", "exec", "jekyll", "serve")
		commd.Dir = jekyllDir
		err = commd.Start()
		if err != nil {
			cancel()
		}
	},
}

func init() {
	rootCmd.AddCommand(tmSpiderSegmentFaultCmd)
}
