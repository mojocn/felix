package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// nesScanCmd represents the felixb command
var nesScanCmd = &cobra.Command{
	Use:   "nesScan",
	Short: "扫描.nes文件生成静态",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		nesDir := `D:\code\tech.mojotv.cn\NESroms`
		fw, err := os.Create(`D:\code\tech.mojotv.cn\_data\nes.yml`)
		if err != nil {
			log.Fatal(err)
		}
		filepath.Walk(nesDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && info.Name() != "NESroms" {
				fmt.Fprintf(fw, "%s:\n", info.Name())
			} else {

				if strings.HasSuffix(path, ".nes") {
					thisDir := filepath.Dir(path)
					fileName := strings.Replace(info.Name(), " ", "_", -1)
					fileName = strings.Replace(fileName, "!", "x", -1)
					fileName = strings.Replace(fileName, "&", "and", -1)
					fileName = strings.Replace(fileName, "-", "_", -1)
					fileName = strings.Replace(fileName, "(", "_", -1)
					fileName = strings.Replace(fileName, ")", "_", -1)
					fileName = strings.Replace(fileName, "[", "_", -1)
					fileName = strings.Replace(fileName, "]", "_", -1)
					fileName = strings.Replace(fileName, "___", "_", -1)
					fileName = strings.Replace(fileName, "__", "_", -1)

					nPath := filepath.Join(thisDir, fileName)
					err := os.Rename(path, nPath)
					if err != nil {
						log.Println("重命名失败", path, nPath, err)
					}
					fn := fmt.Sprintf(`    - "NESroms%s"`, strings.ReplaceAll(nPath, nesDir, ""))
					fn = strings.ReplaceAll(fn, `\`, "/") + "\n"
					fmt.Fprintf(fw, fn)
				}

			}
			return nil
		})

	},
}

func init() {
	rootCmd.AddCommand(nesScanCmd)
}
