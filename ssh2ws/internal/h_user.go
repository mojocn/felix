package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/felix/model"
)

func UserAll(c *gin.Context) {
	mdl := model.User{}
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

func RegisterCommenter(c *gin.Context) {
	var mdl model.User
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	err = mdl.CreateUserOfRole(8)
	if handleError(c, err) {
		return
	}
	jsonData(c, mdl)
}

func UserUpdate(c *gin.Context) {
	var mdl model.User
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

func UserDelete(c *gin.Context) {
	var mdl model.User
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
