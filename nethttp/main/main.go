package main

import (
	"fmt"
	"io"
	"net/http"

	nh "github.com/birowo/route/nethttp"
)

func main() {
	r := new(nh.Handle)
	r.GET("/a/:/b/:", func(w http.ResponseWriter, r *http.Request) {
		rw := w.(nh.RW)
		fmt.Fprintf(w, "params: %q", rw.Params)
	})
	r.POST("/a/:/b/:", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		rw := w.(nh.RW)
		fmt.Fprintf(w, "params: %q\nbody: %s", rw.Params, body)
	})
	http.ListenAndServe(":8080", r)
}
