package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mojocn/felix/util"
	"strconv"

	"github.com/spf13/cobra"
)

// randomstringCmd represents the randomstring command
var randomstringCmd = &cobra.Command{
	Use:   "rands",
	Short: "生成随机字符床",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slen, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			color.Red("ID must be an integer")
			return
		}
		if slen > 0 {

			fmt.Println(util.RandomString(int(slen)))

		}

	},
}

func init() {
	rootCmd.AddCommand(randomstringCmd)

}
