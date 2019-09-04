package cmd

import (
	"github.com/dejavuzhou/felix/fssh"
	"github.com/spf13/cobra"
)

// cronCmd represents the cron command
var telCmd = &cobra.Command{
	Use:   "tel",
	Short: "fake ssh server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fssh.ThisRun()
	},
}

func init() {
	rootCmd.AddCommand(telCmd)
}
