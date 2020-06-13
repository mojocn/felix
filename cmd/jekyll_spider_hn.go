package cmd

import (
	"time"

	"github.com/libragen/felix/spiderhn"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// spiderHNCmd represents the spiderHN command
var spiderHNCmd = &cobra.Command{
	Use:   "spiderHN",
	Short: "tech.mojotv.cn: spider hacknews",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		techMojoSpiderHN()
	},
}

func init() {
	rootCmd.AddCommand(spiderHNCmd)

}

var gitCount = 1

func createCmds() []spiderhn.Cmd {
	gitCount++
	gifConfig1 := []spiderhn.Cmd{
		{"git", []string{"config", "--global", "user.email", "'dejavuzhou@qq.com'"}},
	}
	gifConfig2 := []spiderhn.Cmd{
		{"git", []string{"config", "--global", "user.email", "'1413507308@qq.com'"}},
	}
	cmds := []spiderhn.Cmd{
		{"git", []string{"config", "--global", "user.name", "'EricZhou'"}},
		{"git", []string{"stash"}},
		{"git", []string{"pull", "origin", "master"}},
		{"git", []string{"stash", "apply"}},
		{"git", []string{"add", "."}},
		{"git", []string{"status"}},
		{"git", []string{"commit", "-am", "hacknews-update" + time.Now().Format(time.RFC3339)}},
		{"git", []string{"status"}},
		{"git", []string{"push", "origin", "master"}},
		//{"netstat", []string{"-lntp"}},
		//{"free", []string{"-m"}},
		//{"ps", []string{"aux"}},
	}
	if gitCount%2 == 0 {
		cmds = append(gifConfig2, cmds...)
	} else {
		cmds = append(gifConfig1, cmds...)
	}
	return cmds
}

func techMojoSpiderHN() {
	if err := spiderhn.SpiderHackNews(); err != nil {
		logrus.Error(err)
	}
	if err := spiderhn.SpiderHackShows(); err != nil {
		logrus.Error(err)
	}
	if err := spiderhn.ParsemarkdownHacknews(); err != nil {
		logrus.Error(err)
	}
}
