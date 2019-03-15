package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"mime/multipart"
	"path"
)

func SftpUp(c *gin.Context) {
	sftpClient, err := getSftpClient(c)
	if handleError(c, err) {
		return
	}
	file, err := c.FormFile("file")
	if handleError(c, err) {
		return
	}
	fullPath := c.Query("path")

	err = uploadFile(fullPath, sftpClient, file)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)

}

func uploadFile(desDir string, client *sftp.Client, header *multipart.FileHeader) error {
	if desDir == "$HOME" {
		wd, err := client.Getwd()
		if err != nil {
			return err
		}
		desDir = wd
	}
	srcFile, err := header.Open()
	if err != nil {
		return err
	}
	dstFile, err := client.Create(path.Join(desDir, header.Filename))
	if err != nil {
		return err
	}
	defer srcFile.Close()
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		return err
	}
	return nil
}
