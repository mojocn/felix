package cmd

import (
	"fmt"
	"log"

	"github.com/dejavuzhou/felix/ginbro"
	"github.com/dejavuzhou/felix/model"
	"github.com/spf13/cobra"
)

var gc model.Ginbro

// restCmd represents the rest command
var restCmd = &cobra.Command{
	Use:     "ginbro",
	Short:   "generate a RESTful codebase from SQL database",
	Long:    `generate a RESTful APIs app with gin and gorm for gophers`,
	Example: `felix ginbro -a dev.wordpress.com:3306 -P go_package_name -n db_name -u db_username -p 'my_db_password' -d '~/thisDir'`,
	Run: func(cmd *cobra.Command, args []string) {
		app, err := ginbro.Run(gc)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("cd %s then run go run main.go test your codebase", app.AppDir)
	},
}

func init() {
	rootCmd.AddCommand(restCmd)

	restCmd.Flags().StringVarP(&gc.AppAddr, "listen", "l", "127.0.0.1:5555", "app's listening addr")
	restCmd.Flags().StringVarP(&gc.AppDir, "dir", "d", ".", "code project output directory,default is current working dir")
	restCmd.Flags().StringVarP(&gc.AppPkg, "pkg", "P", "", "eg1: github.com/dejavuzhou/ginSon, eg2: ginbroSon")
	restCmd.Flags().StringVar(&gc.AuthTable, "authTable", "users", "login user table")
	restCmd.Flags().StringVar(&gc.AuthColumn, "authColumn", "password", "bcrypt password column")
	restCmd.Flags().StringVarP(&gc.DbUser, "dbUser", "u", "root", "database username")
	restCmd.Flags().StringVarP(&gc.DbPassword, "dbPassword", "p", "password", "database user password")
	restCmd.Flags().StringVarP(&gc.DbAddr, "dbAddr", "a", "127.0.0.1:3306", "database connection addr")
	restCmd.Flags().StringVarP(&gc.DbName, "dbName", "n", "", "database name")
	restCmd.Flags().StringVarP(&gc.DbChar, "dbChar", "c", "utf8", "database charset")
	restCmd.Flags().StringVarP(&gc.DbType, "dbType", "t", "mysql", "database type: mysql/postgres/mssql/sqlite")

	restCmd.MarkFlagRequired("pkg")
	restCmd.MarkFlagRequired("dbAddr")
}
