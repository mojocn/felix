package cmd

import (
	"os"
	"runtime"
	"strconv"

	"github.com/dejavuzhou/felix/model"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// tasklsCmd represents the taskls command
var tasklsCmd = &cobra.Command{
	Use:   "task",
	Short: "list all rows in TaskList",
	Long:  `usage:felix task`,
	Run: func(cmd *cobra.Command, args []string) {
		mcs, err := model.TaskAll(searchKey)
		if err != nil {
			color.Red("DB error [%s]", err)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Category", "Task", "Status", "DeadLine", "CreatedAt"})
		table.SetBorder(true) // Set Border to false
		//table.SetAutoMergeCells(true)
		//table.SetRowLine(true)

		setListTableColor(table)
		for _, mc := range mcs {
			row := []string{
				strconv.Itoa(int(mc.Id)),
				mc.Category,
				mc.Content,
				mc.Status,
				mc.Deadline.Format(model.TimeLayout),
				mc.CreatedAt.Format(model.TimeLayout),
			}
			table.Append(row)
		}
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(tasklsCmd)
	tasklsCmd.Flags().StringVarP(&searchKey, "search", "s", "", "模糊搜索Task名称")
}

func setListTableColor(table *tablewriter.Table) {
	if isatty.IsCygwinTerminal(os.Stdout.Fd()) || (runtime.GOOS != "windows") {
		table.SetHeaderColor(
			tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgHiGreenColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgHiGreenColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgHiGreenColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgHiGreenColor, tablewriter.Bold},
			tablewriter.Colors{tablewriter.FgHiGreenColor, tablewriter.Bold},
		)
		table.SetColumnColor(
			tablewriter.Colors{tablewriter.FgRedColor},
			tablewriter.Colors{tablewriter.FgCyanColor},
			tablewriter.Colors{tablewriter.FgCyanColor},
			tablewriter.Colors{tablewriter.FgHiYellowColor},
			tablewriter.Colors{tablewriter.FgHiMagentaColor},
			tablewriter.Colors{tablewriter.FgHiWhiteColor},
		)
	}
}
