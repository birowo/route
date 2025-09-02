package nethttp

import (
	"net/http"

	"github.com/birowo/route"
)

type Handle struct {
	get, post route.Nodes[http.HandlerFunc]
}
type RW struct {
	http.ResponseWriter
	Params [][]byte
}

func (h *Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		v       http.HandlerFunc
		prms    route.Params
		prmsIdx int
	)
	switch r.Method {
	case "GET":
		v, prms, prmsIdx = h.get.Get([]byte(r.URL.Path))
	case "POST":
		v, prms, prmsIdx = h.post.Get([]byte(r.URL.Path))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if v == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	v(RW{w, prms[:prmsIdx]}, r)
}
func (h *Handle) GET(path string, hf http.HandlerFunc) {
	h.get.Set([]byte(path), hf)
}
func (h *Handle) POST(path string, hf http.HandlerFunc) {
	h.post.Set([]byte(path), hf)
}
