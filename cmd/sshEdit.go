package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/mojocn/felix/model"

	"github.com/spf13/cobra"
)

// hostUpdateCmd represents the hostEdit command
var hostUpdateCmd = &cobra.Command{
	Use:   "sshedit",
	Short: "update a ssh connection",
	Long:  `update a ssh connection configuration, usually be called after sshud command, usage felix sshedit 1 -n=Awesome`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		argId, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			cmd.Help()
			fmt.Println("ID must be an int:", err)
		}
		if err := model.MachineUpdate(updateName, updateAddr, updateUser, updatePassword, updateKey, updateAuth, uint(argId), updatePort); err != nil {
			log.Fatal(err)
		}
	},
}
var updateKey, updateName, updateAddr, updateUser, updatePassword, updateAuth string
var updatePort uint

func init() {
	rootCmd.AddCommand(hostUpdateCmd)

	hostUpdateCmd.Flags().StringVarP(&updatePassword, "password", "p", "", "ssh pasword")
	hostUpdateCmd.Flags().StringVarP(&updateKey, "key", "k", "", "ssh auth key path")
	hostUpdateCmd.Flags().StringVarP(&updateName, "name", "n", "", "ssh name")
	hostUpdateCmd.Flags().StringVarP(&updateAddr, "addr", "a", "", "ssh domain or ip")
	hostUpdateCmd.Flags().StringVarP(&updateUser, "user", "u", "", "ssh login user name")
	//hostUpdateCmd.Flags().UintVarP(&id, "id", "i", 0, "SSH服务器ID(部署参数)")
	hostUpdateCmd.Flags().StringVarP(&updateAuth, "auth", "", "", "auth type, must be password or key")
	hostUpdateCmd.Flags().UintVar(&updatePort, "port", 0, "ssh port")

}
