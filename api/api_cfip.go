package api

import (
	"encoding/json"
	"fmt"
	"github.com/mojocn/felix/model"
	"github.com/mojocn/felix/util"
	"net/http"
)

func apiCfIpInit(w http.ResponseWriter, req *http.Request) {
	client, err := util.NewCfIP()
	if checkErr(w, err) {
		return
	}
	var rows []model.CfIp
	client.AllIps(func(ip, cidr string) {
		ins := model.CfIp{
			IP:   ip,
			Cidr: cidr,
		}
		rows = append(rows, ins)
	})
	err = model.DB().CreateInBatches(rows, 200).Error
	if checkErr(w, err) {
		return
	}
	responseJson(w, http.StatusOK, "ok")
}

func apiCfIpList(w http.ResponseWriter, req *http.Request) {
	var rows []model.CfIp
	err := model.DB().Limit(100).Find(&rows).Error
	if checkErr(w, err) {
		return
	}
	responseJson(w, http.StatusOK, rows)
}

func apiCfIpUpdate(w http.ResponseWriter, req *http.Request) {
	ins := new(model.CfIp)
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

func apiCfIpCreate(w http.ResponseWriter, req *http.Request) {
	ins := new(model.CfIp)
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

func apiCfIpDelete(w http.ResponseWriter, req *http.Request) {
	ins := new(model.CfIp)
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
