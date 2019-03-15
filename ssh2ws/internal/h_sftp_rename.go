package internal

import (
	"github.com/gin-gonic/gin"
)

func SftpRename(c *gin.Context) {
	sftpClient, err := getSftpClient(c)
	if handleError(c, err) {
		return
	}
	oPath := c.Query("opath")
	nPath := c.Query("npath")
	err = sftpClient.Rename(oPath, nPath)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)

}
