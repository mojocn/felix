package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/libragen/felix/model"
)

func WslogMsgDelete(c *gin.Context) {

	ids := []int{}
	err := c.ShouldBind(&ids)
	if handleError(c, err) {
		return
	}
	var mdl model.WslogMsg
	err = mdl.Delete(ids)
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
