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

	mux.HandleFunc("GET /api/proxies", apiProxyList)
	mux.HandleFunc("PATCH /api/proxies", apiProxyUpdate)
	mux.HandleFunc("POST /api/proxies", apiProxyCreate)
	mux.HandleFunc("DELETE /api/proxies", apiProxyDelete)

	mux.HandleFunc("GET /api/cfip-init", apiCfIpInit)
	mux.HandleFunc("GET /api/cfips", apiCfIpList)
	mux.HandleFunc("PATCH /api/cfips", apiCfIpUpdate)
	mux.HandleFunc("POST /api/cfips", apiCfIpCreate)
	mux.HandleFunc("DELETE /api/cfips", apiCfIpDelete)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return server
}
