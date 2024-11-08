package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/felix/wslog"
)

func Wslog(hub *wslog.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if handleError(c, err) {
			return
		}
		uid, err := mWuserId(c)
		if wshandleError(wsConn, err) {
			return
		}
		//defer wsConn.Close()
		wslog.ClientRunWith(hub, wsConn, uid)

	}
}

func WslogChannel(hub *wslog.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if handleError(c, err) {
			return
		}
		uid, err := mWuserId(c)
		if wshandleError(wsConn, err) {
			return
		}
		//defer wsConn.Close()
		wslog.ClientRunWith(hub, wsConn, uid)

	}
}
