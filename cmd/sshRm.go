package cmd

import (
	"log"
	"strconv"

	"github.com/libragen/felix/model"
	"github.com/spf13/cobra"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "sshrm",
	Short: "delete a ssh connection",
	Long:  `usage: felix sshdel 2`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal("ID must be an integer")
		}
		if err := model.MachineDelete(uint(id)); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}
