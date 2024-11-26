package api

import (
	"encoding/json"
	"github.com/mojocn/felix/model"
	"net/http"
)

func responseJson(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func checkErr(w http.ResponseWriter, err error) (shouldReturn bool) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		return true
	}
	return false
}

func apiMeta(w http.ResponseWriter, req *http.Request) {
	row := new(model.Meta)
	err := model.DB().First(row).Error
	if checkErr(w, err) {
		return
	}
	responseJson(w, http.StatusOK, row)
}
