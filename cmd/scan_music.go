package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var scanMusicCmd = &cobra.Command{
	Use:   "scan-music",
	Short: "scan music files",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		//https://github.com/dhowden/tag music tag
		baseDir := "/code/tech.mojotv.cn"
		list := []music{}
		filepath.Walk(baseDir+"/music", func(path string, info os.FileInfo, err error) error {
			if !strings.HasSuffix(path, ".mp3") {
				return nil
			}
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			mFile, err := os.Open(path)
			if err != nil {
				return err
			}
			m, err := tag.ReadFrom(mFile)
			if err != nil {
				return err
			}

			one := music{
				Name:   m.Title(),
				Artist: m.Artist(),
				Url:    strings.ReplaceAll(path, baseDir, ""),
				Cover:  b64UriImage(m),
			}

			list = append(list, one)

			return nil
		})
		bs, err := json.Marshal(list)
		if err != nil {
			log.Println(err)
			return
		}
		err = os.WriteFile(baseDir+"/music.json", bs, 0644)
		if err != nil {
			log.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(scanMusicCmd)
}

type music struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Url    string `json:"url"`
	Cover  string `json:"cover"`
}

func b64UriImage(tag tag.Metadata) string {
	return "/assets/image/logo00.png"
	pic := tag.Picture()
	if pic == nil {
		return ""
	}

	return fmt.Sprintf("data:%s;%s", pic.MIMEType, base64.StdEncoding.EncodeToString(pic.Data))
}
