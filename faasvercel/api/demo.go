package api

import (
	"net/http"
)

func Demo(w http.ResponseWriter, r *http.Request) {
	// hello world

	// response json hello world
	w.Write([]byte("hello world1"))

}
