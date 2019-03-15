package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"io/ioutil"
	"net/http"
	"os"
)

func scpFetchFile(c *gin.Context) (*sftp.File, os.FileInfo, error) {
	sftpClient, err := getSftpClient(c)
	if err != nil {
		return nil, nil, err
	}
	fullPath := c.Query("path")
	fileInfo, err := sftpClient.Stat(fullPath)
	if err != nil {
		return nil, nil, err
	}
	if fileInfo.IsDir() {
		return nil, nil, fmt.Errorf("%s is not a file", fullPath)
	}
	f, err := sftpClient.Open(fullPath)
	return f, fileInfo, err
}
func SftpDl(c *gin.Context) {
	file, fileInfo, err := scpFetchFile(c)
	defer file.Close()
	if handleError(c, err) {
		return
	}
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, fileInfo.Name()),
	}
	c.DataFromReader(http.StatusOK, fileInfo.Size(), "application/octet-stream", file, extraHeaders)
}
func SftpCat(c *gin.Context) {
	file, fileInfo, err := scpFetchFile(c)
	defer file.Close()
	if handleError(c, err) {
		return
	}
	b, err := ioutil.ReadAll(file)
	if handleError(c, err) {
		return
	}
	//c.String(200,"utf-8",file)
	c.AbortWithStatusJSON(200, gin.H{"ok": true, "data": string(b), "msg": fileInfo.Name()})
}
