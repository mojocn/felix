package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/libragen/felix/model"
	"github.com/libragen/felix/wslog"
)

func WslogHookCreate(c *gin.Context) {
	var mdl model.WslogHook
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	uid, err := mWuserId(c)
	if handleError(c, err) {
		return
	}
	mdl.UserId = uid
	err = mdl.Create()
	if handleError(c, err) {
		return
	}
	jsonData(c, mdl)
}

func WslogHookAll(c *gin.Context) {
	mdl := model.WslogHook{}
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

func WslogHookDelete(c *gin.Context) {
	var mdl model.WslogHook
	id, err := parseParamID(c)
	if handleError(c, err) {
		return
	}
	mdl.Id = id
	err = mdl.Delete()
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}
func WslogHookUpdate(c *gin.Context) {
	var mdl model.WslogHook
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	err = mdl.Update()
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}

func WsLogHookApi(hub *wslog.Hub) gin.HandlerFunc {

	return func(c *gin.Context) {
		mdl := model.WslogMsg{}
		err := c.ShouldBind(&mdl.SlackMsg)
		if handleError(c, err) {
			return
		}
		hookId, err := mWhookApiHookId(c)
		if handleError(c, err) {
			return
		}
		mdl.HookId = hookId
		mdl.UserId = 0
		if handleError(c, err) {
			return
		}
		err = mdl.Create()
		if handleError(c, err) {
			return
		}
		hub.AddMsg(mdl)
		jsonSuccess(c)
	}
}
