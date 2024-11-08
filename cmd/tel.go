package cmd

import (
	"github.com/mojocn/felix/fssh"
	"github.com/spf13/cobra"
)

// cronCmd represents the cron command
var telCmd = &cobra.Command{
	Use:   "tel",
	Short: "fake ssh server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fssh.LoadOrCreateKey()
		fssh.ThisRun()
	},
}

func init() {
	rootCmd.AddCommand(telCmd)
}
