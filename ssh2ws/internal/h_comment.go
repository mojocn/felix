package internal

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/libragen/felix/model"
)

func CommentAll(c *gin.Context) {
	q := &model.CommentQ{}
	err := c.ShouldBindQuery(q)
	if handleError(c, err) {
		return
	}
	data, err := q.SearchAll()
	if handleError(c, err) {
		return
	}
	json200(c, data)
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
