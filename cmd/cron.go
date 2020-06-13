package cmd

import (
	"github.com/libragen/felix/cronjob"
	"github.com/spf13/cobra"
)

// cronCmd represents the cron command
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "每3小时spider Hacknews jekyll build",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := cronjob.NewScheduler()
		s.Every(1).Hours().Do(techMojoSpiderHN)
		<-s.Start()
	},
}

func init() {
	rootCmd.AddCommand(cronCmd)
}
