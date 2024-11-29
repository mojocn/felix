package rsver

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	//json response hello world
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}
func ping(w http.ResponseWriter, r *http.Request) {
	//json response hello world
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ping": "pong"}`))
}
func Run() {
	fmt.Println("WebSocket server starting on 127.0.0.1:80")
	http.HandleFunc("/ws-s5", wsS5)
	http.HandleFunc("/ws-vless", wsVless)
	http.HandleFunc("/", home)
	http.HandleFunc("/ping", ping)
	log.Fatal(http.ListenAndServe(":80", nil))
}
