package cmd

import (
	"github.com/mojocn/felix/ssh2ws"
	"log"

	"github.com/libragen/felix/model"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "purge all felix configuration",
	Long:  `purge all felix info by destroying SQLite database file`,
	Run: func(cmd *cobra.Command, args []string) {
		ssh2ws.RunSsh2ws("", "", "", "", 0, false)
		if err := model.FlushSqliteDb(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
