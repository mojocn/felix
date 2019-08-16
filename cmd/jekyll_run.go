package cmd

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tmSpiderSegmentFaultCmd represents the taskrm command
var jekyllServeCmd = &cobra.Command{
	Use:   "jekyllb",
	Short: "bundle execu jekyll $1",
	Long:  `alias bundle exec jekyll in golang`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		techMojoJekyllRun("build")
	},
}

func init() {
	rootCmd.AddCommand(jekyllServeCmd)
}

func techMojoJekyllRun(buildOrServe string) {
	jekyllDir := viper.GetString("tech_mojotv_cn.srcDir")

	thisCmd := exec.Command("bundle", "exec", "jekyll", buildOrServe)
	if runtime.GOOS == "linux" {
		thisCmd.Env = os.Environ()
		//解决supervisor 找不到bundle jekyll 的问题
		thisCmd.Env = append(thisCmd.Env, "PATH=/usr/local/rvm/gems/ruby-2.5.5/bin:/usr/local/rvm/gems/ruby-2.5.5@global/bin:/usr/local/rvm/rubies/ruby-2.5.5/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/usr/local/rvm/bin:/usr/local/go/bin:/gopath/bin:/root/.gem/ruby/bin:/root/bin")
	}
	thisCmd.Dir = jekyllDir
	o, err := thisCmd.CombinedOutput()
	logrus.Info(string(o))
	if err != nil {
		logrus.Error(err)
	}
}
