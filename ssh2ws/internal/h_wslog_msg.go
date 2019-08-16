package internal

import (
	"github.com/dejavuzhou/felix/model"
	"github.com/gin-gonic/gin"
)

func WslogMsgDelete(c *gin.Context) {

	var mdl model.WslogMsg
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

func WslogMsgAll(c *gin.Context) {
	mdl := model.WslogMsg{}
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
