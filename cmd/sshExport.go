package cmd

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/mojocn/felix/model"
	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
)

// sshExportCmd represents the sshexport command
var sshExportCmd = &cobra.Command{
	Use:   "sshexport",
	Short: "export all ssh connection configuration to a csv file",
	Long: `export all ssh connection info to a csv file,
for massively editing ssh connection,
after that use felix "sshimport -f 'path' -F" to update ssh connection
usage: felix sshexport
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := exportAllSshInfoToHomeDirCsvFile(); err != nil {
			log.Fatal(err)
		}
		color.Cyan("you can use export all ssh rows")
		color.Cyan("then edit them in excel")
		color.Cyan("use the file to felix sshimport with F flag")
		color.Yellow("update all ssh info massively")
	},
}

func init() {
	rootCmd.AddCommand(sshExportCmd)
}

func exportAllSshInfoToHomeDirCsvFile() error {
	mcs, err := model.MachineAll("")
	if err != nil {
		return err
	}
	filePath, _ := homedir.Expand("~/allSshInfo.csv")
	csvFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	rows := [][]string{
		{"ssh_user(optional can be blank)", "ssh_password(optional can be blank)", "ssh_name", "ssh_host", "ssh_port (optional can be blank)"},
	}
	for _, mc := range mcs {
		one := []string{mc.User, mc.Password, mc.Name, mc.Host, strconv.Itoa(int(mc.Port))}
		rows = append(rows, one)
	}

	err = csvWriter.WriteAll(rows)
	if err != nil {
		return err
	}
	color.Cyan("ssh import csv template has exported into %s", filePath)
	color.Yellow("use Excel to add ssh info into a row")
	return nil
}
