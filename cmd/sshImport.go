package cmd

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/libragen/felix/model"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// sshImportCmd represents the testImport command
var sshImportCmd = &cobra.Command{
	Use:   "sshimport",
	Short: "import massive ssh connection configuration from a csv file",
	Long:  `usage: felix sshimport -f import.csv -u felix -p pewdiepie -F`,
	Run: func(cmd *cobra.Command, args []string) {
		if isFlushSsh {
			if err := model.MachineDeleteAll(); err != nil {
				log.Fatalln(err)
			}
		}

		if isExportTemplate {
			exportCsvTemplateToHomeDir()
		} else {
			importHost()
		}
	},
}
var imPassword, imFile, imUser, imKey, imAuth string
var isExportTemplate, isFlushSsh bool

func init() {
	rootCmd.AddCommand(sshImportCmd)
	sshImportCmd.Flags().StringVarP(&imFile, "file", "f", ``, "a csv file has a lot ssh server info")
	sshImportCmd.Flags().StringVarP(&imPassword, "password", "p", "", "ssh password")
	sshImportCmd.Flags().StringVarP(&imUser, "user", "u", "", "ssh username")
	sshImportCmd.Flags().StringVarP(&imKey, "key", "k", "~/.ssh/id_rsa", "default SSH Private Key path")
	sshImportCmd.Flags().StringVarP(&imAuth, "auth", "", "password", "auth type only allows password and key")
	sshImportCmd.Flags().BoolVarP(&isExportTemplate, "template", "t", false, "is export csv template into HOME dir")
	sshImportCmd.Flags().BoolVarP(&isFlushSsh, "flush", "F", false, "is Flush all ssh rows then import csv")
}

func importHost() {
	file, err := os.Open(imFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(record[0], "can be blank") {
			continue
		}
		sshUser := imUser
		sshPassword := imPassword
		if record[0] != "" {
			sshUser = record[0]
		}
		if record[1] != "" {
			sshPassword = record[1]
		}
		var sshPort uint = 22
		if i, err := strconv.ParseUint(record[4], 10, 64); err != nil && i != 0 {
			sshPort = uint(i)
		}
		if err := model.MachineAdd(record[2], record[2], record[3], sshUser, sshPassword, imKey, imAuth, sshPort); err != nil {
			color.Red("db: %s", err)
		}

	}

}

func exportCsvTemplateToHomeDir() {
	filePath, _ := homedir.Expand("~/sshImportCsvTemplate.csv")
	csvFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	rows := [][]string{
		{"ssh_user(optional can be blank)", "ssh_password(optional can be blank)", "ssh_name", "ssh_host", "ssh_port (optional can be blank)"},
	}
	err = csvWriter.WriteAll(rows)
	if err != nil {
		log.Fatalln(err)
	}
	color.Cyan("ssh import csv template has exported into %s", filePath)
	color.Yellow("use Excel to add ssh info into a row")
}
