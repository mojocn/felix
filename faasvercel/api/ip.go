package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func MyInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["ip"] = r.RemoteAddr
	resp["user-agent"] = r.UserAgent()
	resp["accept-language"] = r.Header.Get("Accept-Language")
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("Error happened in JSON marshal. Err: %s", err)
	} else {
		w.Write(jsonResp)
	}
}
