package internal

import (
	"github.com/dejavuzhou/felix/wslog"
	"github.com/gin-gonic/gin"
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
