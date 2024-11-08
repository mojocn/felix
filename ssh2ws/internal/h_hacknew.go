package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/felix/model"
)

func HackNewAll(c *gin.Context) {
	q := &model.HackNewQ{}
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

func HackNewUpdate(c *gin.Context) {
	var mdl model.HackNew
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
func HackNewRm(c *gin.Context) {
	ids := []uint{}
	var mdl model.HackNew
	err := c.ShouldBind(&ids)
	if handleError(c, err) {
		return
	}
	err = mdl.Delete(ids)
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}
