package api

import (
	"net/http"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func ApiServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/foo/bar", apiHandler{})
	mux.HandleFunc("GET /api/meta", apiMeta)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return server
}
