package cmd

import (
	"github.com/spf13/cobra"
)

// tmSpiderSegmentFaultCmd represents the taskrm command
var tmBuildJekyllCmd = &cobra.Command{
	Use:   "tmb",
	Short: "tech.mojotv.cn: spider hacknews && jekyll build",
	Long:  `只能在我家里面的电脑使用`,
	Run: func(cmd *cobra.Command, args []string) {
		techMojoSpiderHN()
		techMojoJekyllRun("build")
	},
}

func init() {
	rootCmd.AddCommand(tmBuildJekyllCmd)
}
