package cmd

import (
	"fmt"
	"github.com/libragen/felix/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

// slackCmd represents the slack command
var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		tc := time.NewTicker(time.Second * 2)
		defer tc.Stop()
		for {
			select {
			case <-tc.C:
				title := util.RandomString(34)
				logrus.WithField("time", time.Now()).WithField("fint", 1).WithField("fBool", false).WithField("fstring", "awesome").WithField("fFloat", 0.45).WithError(fmt.Errorf("error fmt format: %s", "felix is awesome")).Error("this mgs ", "error ", title)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(slackCmd)
}
