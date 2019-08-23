package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dejavuzhou/felix/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func jsonError(c *gin.Context, msg interface{}) {
	c.AbortWithStatusJSON(200, gin.H{"ok": false, "msg": msg})
}
func jsonAuthError(c *gin.Context, msg interface{}) {
	c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"ok": false, "msg": msg})
}

func jsonData(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(200, gin.H{"ok": true, "data": data})
}

//func jsonPagination(c *gin.Context, list interface{}, total uint, query *model.PaginationQ) {
//	c.AbortWithStatusJSON(200, gin.H{"ok": true, "data": list, "total": total, "offset": query.Offset, "limit": query.Size})
//}
func jsonSuccess(c *gin.Context) {
	c.AbortWithStatusJSON(200, gin.H{"ok": true, "msg": "success"})
}
func jsonPagination(c *gin.Context, list interface{}, total uint, query *model.PaginationQ) {
	c.AbortWithStatusJSON(200, gin.H{"ok": true, "data": list, "total": total, "page": query.Page, "size": query.Size})
}
func json200(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(200, data)
}
func handleError(c *gin.Context, err error) bool {
	if err != nil {
		//logrus.WithError(err).Error("gin context http handler error")
		jsonError(c, err.Error())
		return true
	}
	return false
}
func handlerAuthMiddlewareError(c *gin.Context, err error) bool {
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"msg": err.Error()})
		return true
	}
	return false
}

func wshandleError(ws *websocket.Conn, err error) bool {
	if err != nil {
		logrus.WithError(err).Error("handler ws ERROR:")
		dt := time.Now().Add(time.Second)
		if err := ws.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); err != nil {
			logrus.WithError(err).Error("websocket writes control message failed:")
		}
		return true
	}
	return false
}

func parseParamID(c *gin.Context) (uint, error) {
	id := c.Param("id")
	parseId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, errors.New("id must be an unsigned int")
	}
	return uint(parseId), nil
}

func getAuthUser(c *gin.Context) (*model.User, error) {
	user := model.User{}
	uid, err := mWuserId(c)
	if err != nil {
		return nil, err
	}
	user.Id = uid
	err = user.One()
	if err != nil {
		return nil, fmt.Errorf("context can not get user of %d,error:%s", user.Id, err)
	}
	return &user, nil
}
func mWuserId(c *gin.Context) (uint, error) {
	return mWcontextGetUintKey(c, contextKeyUid)
}
func mWhookApiHookId(c *gin.Context) (uint, error) {
	return mWcontextGetUintKey(c, contextKeyWslogHookId)
}

func mWcontextGetUintKey(c *gin.Context, key string) (uint, error) {
	v, exist := c.Get(key)
	if !exist {
	}
	uintV, ok := v.(uint)
	if ok {
		return uintV, nil
	}
	return 0, fmt.Errorf("key for %s in gin.Context value is %v is not a uint type", key, v)
}
