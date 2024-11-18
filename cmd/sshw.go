package cmd

import (
	"github.com/mojocn/felix/ssh2ws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"time"
)

// sshwCmd represents the sshw command
var sshwCmd = &cobra.Command{
	Use:   "sshw",
	Short: "open a web UI for Felix https://localhost:2222",
	Long:  `the demo website is https://felix.mojotv.cn`,
	Run: func(cmd *cobra.Command, args []string) {

		addr := viper.GetString("sshw.addr")
		user := viper.GetString("sshw.user")
		password := viper.GetString("sshw.password")
		secret := viper.GetString("sshw.secret")
		ex := time.Duration(viper.GetInt("sshw.expire")) * time.Hour * 24

		if l := len(secret); l != 32 {
			log.Fatalf("secret length is %d,but its length must be 32", l)
		}

		if err := ssh2ws.RunSsh2ws(addr, user, password, secret, ex, verbose); err != nil {
			log.Fatal(err)
		}
	},
}
var expire uint
var secret string

func init() {
	rootCmd.AddCommand(sshwCmd)
	sshwCmd.Flags().StringVarP(&secret, "secret", "s", "", "jwt secret string length must be 32")
	sshwCmd.Flags().StringVarP(&addr, "addr", "a", ":2222", "listening addr")
	sshwCmd.Flags().StringVarP(&user, "user", "u", "admin", "auth user")
	sshwCmd.Flags().StringVarP(&password, "password", "p", "admin", "auth password")
	sshwCmd.Flags().UintVarP(&expire, "expire", "x", 60*24*30, "token expire in * minute")
}
