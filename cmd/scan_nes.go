package cmd

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type nesItem struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// scanNesCmd represents the felixb command
var scanNesCmd = &cobra.Command{
	Use:   "scanNes",
	Short: "扫描.nes文件生成静态",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		baseDir := `/data`
		baseURL := "https://s3.mojotv.cn"
		var list []nesItem
		filepath.Walk(baseDir+"/nes", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.Contains(path, "Hack") || strings.Contains(path, "hack") {
				os.Remove(path)
			}
			if strings.HasSuffix(path, ".nes") {
				fileName := strings.ToLower(info.Name())
				item := nesItem{
					Name: strings.TrimSuffix(fileName, ".nes"),
					Url:  baseURL + strings.ReplaceAll(path, baseDir, ""),
				}
				list = append(list, item)
			}
			return nil
		})
		bs, err := json.Marshal(list)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(bs))
		os.WriteFile(baseDir+"/nes.json", bs, 0644)
	},
}

func init() {
	rootCmd.AddCommand(scanNesCmd)
}
