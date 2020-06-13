package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/libragen/felix/model"
)

func SshLogAll(c *gin.Context) {
	query := &model.SshLogQ{}
	err := c.ShouldBindQuery(query)
	if handleError(c, err) {
		return
	}
	list, total, err := query.Search()
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, &query.PaginationQ)
}

func SshLogUpdate(c *gin.Context) {
	var mdl model.SshLog
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

func SshLogDelete(c *gin.Context) {
	var mdl model.SshLog
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
