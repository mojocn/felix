package internal

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mojocn/felix/flx"
	"github.com/mojocn/felix/model"
	"github.com/mojocn/felix/util"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// handle webSocket connection.
// first,we establish a ssh connection to ssh server when a webSocket comes;
// then we deliver ssh data via ssh connection between browser and ssh server.
// That is, read webSocket data from browser (e.g. 'ls' command) and send data to ssh server via ssh connection;
// the other hand, read returned ssh data from ssh server and write back to browser via webSocket API.
func WsSsh(c *gin.Context) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if handleError(c, err) {
		return
	}
	defer wsConn.Close()

	cIp := c.ClientIP()

	userM, err := getAuthUser(c)
	if handleError(c, err) {
		return
	}
	cols, err := strconv.Atoi(c.DefaultQuery("cols", "80"))
	if wshandleError(wsConn, err) {
		return
	}
	rows, err := strconv.Atoi(c.DefaultQuery("rows", "40"))
	if wshandleError(wsConn, err) {
		return
	}
	idx, err := parseParamID(c)
	if wshandleError(wsConn, err) {
		return
	}
	mc, err := model.MachineFind(idx)
	if wshandleError(wsConn, err) {
		return
	}

	client, err := flx.NewSshClient(mc)
	if wshandleError(wsConn, err) {
		return
	}
	defer client.Close()
	startTime := time.Now()
	ssConn, err := util.NewSshConn(cols, rows, client)
	if wshandleError(wsConn, err) {
		return
	}
	defer ssConn.Close()

	sws, err := model.NewLogicSshWsSession(cols, rows, true, client, wsConn)
	if wshandleError(wsConn, err) {
		return
	}
	defer sws.Close()

	quitChan := make(chan bool, 3)
	sws.Start(quitChan)
	go sws.Wait(quitChan)

	<-quitChan
	//保存日志

	//write logs
	xtermLog := model.SshLog{
		StartedAt: startTime,
		UserId:    userM.Id,
		Log:       sws.LogString(),
		MachineId: idx,
		ClientIp:  cIp,
	}
	err = xtermLog.Create()
	if wshandleError(wsConn, err) {
		return
	}
}
