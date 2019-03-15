package internal

import (
	"github.com/gin-gonic/gin"
)

func SftpRm(c *gin.Context) {
	sftpClient, err := getSftpClient(c)
	if handleError(c, err) {
		return
	}
	fullPath := c.Query("path")
	if fullPath == "/" || fullPath == "$HOME" {
		jsonError(c, "can't delete / or $HOME dir")
		return
	}

	err = sftpClient.Remove(fullPath)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)

}
