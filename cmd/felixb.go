package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"runtime"
)

// felixbCmd represents the felixb command
var felixbCmd = &cobra.Command{
	Use:   "felixb",
	Short: "在centos7 上编译felix",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		command := `rm -rf ${GOBIN}/felix;
go install;
pkill -f 'felix sshw';
netst -lntp;
ps -ef | grep felix;
pkill -f 'felix sshw';
`
		ls, err := runCmd(command)
		fmt.Println(ls)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(felixbCmd)

}

func runCmd(cmd string) (string, error) {
	execName := "/bin/sh"
	if runtime.GOOS == "windows" {
		execName = "C:\\Program Files\\Git\\bin\\bash.exe"
	}
	if runtime.GOOS == "linux" {
		execName = "/bin/sh"
	}

	b, err := exec.Command(execName, cmd).Output()
	if err != nil {
		return "", err
	}
	return string(b), nil
}
