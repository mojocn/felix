package cmd

import (
	"github.com/dejavuzhou/felix/util"
	"github.com/spf13/cobra"
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "open a tab in browser to convert json to golang struct",
	Long:  `open a tab in browser to convert json to golang struct powered by https://quicktype.io/`,
	Run: func(cmd *cobra.Command, args []string) {
		util.BrowserOpen("https://app.quicktype.io/")
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
}
