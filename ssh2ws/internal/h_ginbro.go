package internal

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dejavuzhou/felix/ginbro"
	"github.com/dejavuzhou/felix/model"
	"github.com/gin-gonic/gin"
)

func GinbroAll(c *gin.Context) {
	mdl := model.Ginbro{}
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if handleError(c, err) {
		return
	}
	list, total, err := mdl.All(query)
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, query)
}
func GinbroGen(c *gin.Context) {

	var m model.Ginbro
	err := c.ShouldBind(&m)
	if handleError(c, err) {
		return
	}
	err = m.Create()
	if handleError(c, err) {
		return
	}
	app, err := ginbro.Run(m)
	if handleError(c, err) {
		return
	}
	err = app.ListAppFileTree()
	if handleError(c, err) {
		return
	}
	jsonData(c, app)
}
func GinbroDb(c *gin.Context) {
	var gb model.Ginbro
	err := c.ShouldBind(&gb)
	if handleError(c, err) {
		return
	}
	data, err := ginbro.FetchDbColumn(gb)
	if handleError(c, err) {
		return
	}
	jsonData(c, data)
}

func GinbroDownload(c *gin.Context) {
	srcPath := c.Query("p")
	if srcPath == "" {
		jsonError(c, "query argument is required")
		return
	}
	srcPath = filepath.Clean(srcPath)

	buf := new(bytes.Buffer)

	// CreateUserOfRole a new zip archive.
	w := zip.NewWriter(buf)

	err := filepath.Walk(srcPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignore directories and hidden files.
		// No entry is needed for directories in a zip file.
		// Each file is represented with a path, no directory
		// entities are required to build the hierarchy.
		if fi.IsDir() || strings.HasPrefix(fi.Name(), ".") {
			return nil
		}
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fHeader, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		fHeader.Name = filepath.ToSlash(relPath)
		fHeader.Method = zip.Deflate
		fHeader.Comment = "https://github.com/mojocn"

		writer, err := w.CreateHeader(fHeader)
		if err != nil {
			return err
		}
		_, err = writer.Write(b)
		return err
	})
	if handleError(c, err) {
		return
	}
	err = w.Close()
	if handleError(c, err) {
		return
	}
	zipName := time.Now().Format("FelixGinbro_2006_01_02T15_04_05.zip")
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, zipName),
	}

	c.DataFromReader(http.StatusOK, int64(buf.Len()), "application/zip", buf, extraHeaders)

}

func GinbroModel(c *gin.Context) {
	var gb model.Ginbro
	err := c.ShouldBind(&gb)
	if handleError(c, err) {
		return
	}
}
