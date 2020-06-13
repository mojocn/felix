package cmd

import (
	"log"
	"strconv"

	"github.com/libragen/felix/flx"
	"github.com/libragen/felix/model"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "scpd",
	Short: "scp download file or folder",
	Long:  `download file or folder, usage: felix sshdl 2 -r="/home/root/awesome" -l="D;/awesome"`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbId, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal("ID must be an int", err)
		}
		h, err := model.MachineFind(uint(dbId))
		if err != nil {
			log.Fatal("wrong ID:", err)
		}
		err = flx.ScpRL(h, remotePath, localPath)
		if err != nil {
			log.Fatal(err)
		}
	},
}
var localPath, remotePath string

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&remotePath, "remote", "r", "", "ssh server file/folder remote path")
	downloadCmd.Flags().StringVarP(&localPath, "local", "l", "", "local path/folder path")
	downloadCmd.MarkFlagRequired("remote")
	downloadCmd.MarkFlagRequired("local")
}
