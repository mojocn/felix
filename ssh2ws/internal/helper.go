package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mojocn/felix/model"
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
	v, exist := c.Get(contextKeyUserObj)
	if !exist {
		return nil, errors.New(contextKeyUserObj + " not exist")
	}
	user, ok := v.(model.User)
	if !ok {
		return nil, fmt.Errorf("v:%v is not type of model.User", user)
	}
	return &user, nil
}
func mWuserId(c *gin.Context) (uint, error) {
	v, exist := c.Get(contextKeyUserObj)
	if !exist {
		return 0, errors.New(contextKeyUserObj + " not exist")
	}
	user, ok := v.(model.User)
	if ok {
		return user.Id, nil
	}
	return 0, errors.New("can't convert to user struct")
}
func mWhookApiHookId(c *gin.Context) (uint, error) {
	return mWcontextGetUintKey(c, contextKeyWslogHookId)
}

func mWcontextGetUintKey(c *gin.Context, key string) (uint, error) {
	v, exist := c.Get(key)
	if !exist {
		return 0, errors.New(key + " not exist")
	}
	uintV, ok := v.(uint)
	if ok {
		return uintV, nil
	}
	return 0, fmt.Errorf("key for %s in gin.Context value is %v is not a uint type", key, v)
}
