package internal

import (
	"strings"

	"github.com/dejavuzhou/felix/model"
	"github.com/gin-gonic/gin"
)

type commentQ struct {
	model.PaginationQ
	model.Comment
}

func CommentAll(c *gin.Context) {
	q := &commentQ{}
	err := c.ShouldBindQuery(q)
	if handleError(c, err) {
		return
	}
	mdl := q.Comment
	query := &q.PaginationQ
	list, total, err := mdl.All(query)
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, query)
}

func CommentCreate(c *gin.Context) {
	uid, err := mWuserId(c)
	if handleError(c, err) {
		return
	}
	var mdl model.Comment
	err = c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	mdl.Content = strings.TrimSpace(mdl.Content)
	if len(mdl.Content) == 0 {
		jsonError(c, "content can't be empty string")
		return
	}
	mdl.UserId = uid
	err = mdl.Create()
	if handleError(c, err) {
		return
	}
	jsonData(c, mdl)
}

func CommentAction(c *gin.Context) {
	id, err := parseParamID(c)
	if handleError(c, err) {
		return
	}
	uid, err := mWuserId(c)
	if handleError(c, err) {
		return
	}
	action := c.Param("action")

	err = (&model.Comment{}).Action(id, uid, action)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}

func CommentDelete(c *gin.Context) {
	var mdl model.Comment
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
