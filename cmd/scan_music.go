package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	"github.com/spf13/cobra"
)

var scanMusicCmd = &cobra.Command{
	Use:   "scanMusic",
	Short: "scan music files",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		//https://github.com/dhowden/tag music tag
		baseDir := "/data"
		baseURL := "https://s3.mojotv.cn"
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
				Url:    baseURL + strings.ReplaceAll(path, baseDir, ""),
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
		err = os.WriteFile(baseDir+"/music.json", bs, 0666)
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
	pic := tag.Picture()
	if pic == nil {
		return ""
	}

	// Resize pic.Data to 100x100 using Go stdlib
	img, _, err := image.Decode(bytes.NewReader(pic.Data))
	if err != nil {
		return ""
	}
	size := 100
	// Create a new 100x100 RGBA image
	resized := image.NewRGBA(image.Rect(0, 0, size, size))
	// Scale the image to 100x100 using NearestNeighbor manually (since stdlib has no draw.NearestNeighbor)
	// Simple nearest-neighbor resize implementation
	bounds := img.Bounds()
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			srcX := bounds.Min.X + (x * bounds.Dx() / size)
			srcY := bounds.Min.Y + (y * bounds.Dy() / size)
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}

	// Encode the resized image back to bytes as JPEG, regardless of original format
	var buf bytes.Buffer
	jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 25}) // Always encode as JPEG, lower quality

	return fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(buf.Bytes()))
}
