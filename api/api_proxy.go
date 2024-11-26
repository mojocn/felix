package api

import (
	"encoding/json"
	"fmt"
	"github.com/mojocn/felix/model"
	"net/http"
)

func apiProxyList(w http.ResponseWriter, req *http.Request) {
	rows := []model.Proxy{}
	err := model.DB().Find(&rows).Error
	if checkErr(w, err) {
		return
	}
	responseJson(w, http.StatusOK, rows)
}

func apiProxyUpdate(w http.ResponseWriter, req *http.Request) {
	ins := new(model.Proxy)
	err := json.NewDecoder(req.Body).Decode(ins)
	if checkErr(w, err) {
		return
	}
	err = model.DB().Save(ins).Error
	if checkErr(w, err) {
		return
	}
	responseJson(w, http.StatusOK, ins)
}

func apiProxyCreate(w http.ResponseWriter, req *http.Request) {
	ins := new(model.Proxy)
	err := json.NewDecoder(req.Body).Decode(ins)
	if checkErr(w, err) {
		return
	}
	ins.ID = 0
	err = model.DB().Save(ins).Error
	if checkErr(w, err) {
		return
	}
	responseJson(w, http.StatusOK, ins)
}

func apiProxyDelete(w http.ResponseWriter, req *http.Request) {
	ins := new(model.Proxy)
	err := json.NewDecoder(req.Body).Decode(ins)
	if checkErr(w, err) {
		return
	}
	if ins.ID == 0 {
		err = fmt.Errorf("id can not be 0")
	}
	if checkErr(w, err) {
		return
	}
	err = model.DB().Delete(ins).Error
	if checkErr(w, err) {
		return
	}
	responseJson(w, http.StatusOK, ins)
}
