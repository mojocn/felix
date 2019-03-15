package cmd

import (
	"log"
	"strconv"

	"github.com/dejavuzhou/felix/flx"
	"github.com/dejavuzhou/felix/model"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "sshproxy",
	Short: "ssh port proxy",
	Long:  `usage : felix sshproxy 2 -l 127.0.0.1:5555 -r 127.0.0.1:3306`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			flx.AllMachines("")
			return
		}
		dbId, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal("ID must be an integer:", err)
		}
		h, err := model.MachineFind(uint(dbId))
		if err != nil {
			log.Fatal("wrong ID", err)
		}
		color.Cyan("porxy ssh's (%s) to local: (%s)...", remoteAddr, localAddr)
		if err := flx.RunProxy(h, localAddr, remoteAddr); err != nil {
			log.Fatal(err)
		}
	},
}
var localAddr, remoteAddr string

func init() {
	rootCmd.AddCommand(proxyCmd)
	proxyCmd.Flags().StringVarP(&localAddr, "local", "l", "127.0.0.1:3306", "local addr")
	proxyCmd.Flags().StringVarP(&remoteAddr, "remote", "r", "127.0.0.1:3306", "remote addr")
}
