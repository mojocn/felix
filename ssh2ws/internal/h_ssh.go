package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/felix/flx"
	"github.com/mojocn/felix/model"
)

func SshAll(c *gin.Context) {
	mdl := model.Machine{}
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if handleError(c, err) {
		return
	}
	list, total, err := mdl.All(query)
	if handleError(c, err) {
		return
	}
	var mcs []model.Machine
	for _, vv := range *list {
		vv.Password = ""
		vv.Key = ""
		mcs = append(mcs, vv)
	}
	jsonPagination(c, mcs, total, query)
}
func SshOne(c *gin.Context) {
	id, err := parseParamID(c)
	if handleError(c, err) {
		return
	}
	mac := model.Machine{}
	mac.Id = id
	err = mac.One()
	if handleError(c, err) {
		return
	}
	mac.Password = ""
	mac.Key = ""
	info, err := flx.FetchHardwareInfo(&mac)
	if handleError(c, err) {
		return
	}
	//data := gin.H{"mac":mac,"info":info}
	jsonData(c, info)
}
func SshCreate(c *gin.Context) {
	var mdl model.Machine
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	err = mdl.Create()
	if handleError(c, err) {
		return
	}
	jsonData(c, mdl)
}

func SshUpdate(c *gin.Context) {
	var mdl model.Machine
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

func SshDelete(c *gin.Context) {
	var mdl model.Machine
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
