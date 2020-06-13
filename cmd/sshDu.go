package cmd

import (
	"strconv"

	"github.com/fatih/color"
	"github.com/libragen/felix/model"
	"github.com/spf13/cobra"
)

// hostCpCmd represents the hostCp command
var hostCpCmd = &cobra.Command{
	Use:   "sshdu",
	Short: "duplicate a ssh connection",
	Long:  `duplicate a ssho info for quick create a ssh connection by using sshedit cmd,usge: felix sshdu 1`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		argId, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			cmd.Help()
			color.Yellow("ID must be a int:", err)
		}
		if err := model.MachineDuplicate(uint(argId)); err != nil {
			color.Red("%s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(hostCpCmd)
}
