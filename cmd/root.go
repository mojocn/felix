package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/mojocn/felix/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "felix",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if isShowVersion {
			fmt.Println("Golang Env: %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
			fmt.Println("UTC build time:%s", buildTime)
			fmt.Println("Build from Github repo version: https://github.com/mojocn/felix/commit/%s", gitHash)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(bTime, gHash string) {
	buildTime = bTime
	gitHash = gHash
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var buildTime, gitHash string
var verbose, isShowVersion bool

func init() {
	cobra.OnInitialize(initFunc)
	rootCmd.Flags().BoolVarP(&isShowVersion, "version", "v", false, "show binary build information")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "verbose")
}

func initFunc() {
	model.CreateSQLiteDb(verbose)
	setupLog()
}
func setupLog() {
	lvl := logrus.InfoLevel
	logrus.SetLevel(lvl)
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)

}
