package cmd

import (
	"os/exec"

	"github.com/dejavuzhou/felix/cronjob"
	"github.com/dejavuzhou/felix/spiderhn"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cronCmd represents the cron command
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "每3小时spider Hacknews jekyll build",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		s := cronjob.NewScheduler()
		s.Every(3).Hours().Do(spiderHacknews)
		<-s.Start()
	},
}

func init() {
	rootCmd.AddCommand(cronCmd)
}

func spiderHacknews() {
	if err := spiderhn.SpiderHackNews(); err != nil {
		logrus.Error(err)
	}
	if err := spiderhn.SpiderHackShows(); err != nil {
		logrus.Error(err)
	}
	if err := spiderhn.ParsemarkdownHacknews(); err != nil {
		logrus.Error(err)
	}

	jekyllCmd := exec.Command("bundle", "exec", "jekyll", "build")
	proDir := viper.GetString("tech_mojotv_cn.srcDir")
	jekyllCmd.Dir = proDir
	b, err := jekyllCmd.Output()
	logrus.Info(string(b))
	if err != nil {
		logrus.Error(err)
	}

}
