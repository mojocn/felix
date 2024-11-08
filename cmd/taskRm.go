package cmd

import (
	"log"
	"strconv"

	"github.com/mojocn/felix/model"
	"github.com/spf13/cobra"
)

// taskrmCmd represents the taskrm command
var taskrmCmd = &cobra.Command{
	Use:   "taskrm",
	Short: "remove a row in TaskList",
	Long:  `usage: felix taskrm 1`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal("ID must be an integer")
		}
		if err := model.TaskRm(uint(id)); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(taskrmCmd)
}
