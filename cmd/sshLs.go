package cmd

import (
	"github.com/dejavuzhou/felix/flx"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "sshls",
	Short: "list all ssh connection configuration or search by hostname",
	Long:  `usage: felix sshls -s ".cn",search ssh by hostname`,
	Run: func(cmd *cobra.Command, args []string) {
		flx.AllMachines(searchKey)
	},
}
var searchKey string

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&searchKey, "search", "s", "", "模糊搜索ssh服务器名称")
}
