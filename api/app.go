package api

import (
	"log"
	"net/http"
)

type apiHandler struct{}

func (apiHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func AdminServer(addr string) *http.Server {
	log.Println("http api server starting on", addr)
	mux := http.NewServeMux()
	mux.Handle("/api/foo", apiHandler{})
	mux.HandleFunc("GET /api/meta", apiMeta)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return server
}
