package main

import (
	"fmt"

	"github.com/birowo/route"
	"github.com/valyala/fasthttp"
)

type Handle struct {
	get, post route.Nodes[fasthttp.RequestHandler]
}

const Params = 0

func (h *Handle) Handler(ctx *fasthttp.RequestCtx) {
	var (
		v       fasthttp.RequestHandler
		params  route.Params
		prmsIdx int
	)
	switch string(ctx.Method()) {
	case "GET":
		v, params, prmsIdx = h.get.Get(ctx.Path())
	case "POST":
		v, params, prmsIdx = h.post.Get(ctx.Path())
	default:
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		return
	}
	if v == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	ctx.SetUserValue(Params, params[:prmsIdx])
	v(ctx)
}
func (h *Handle) GET(path []byte, rh fasthttp.RequestHandler) {
	h.get.Set(path, rh)
}
func (h *Handle) POST(path []byte, rh fasthttp.RequestHandler) {
	h.post.Set(path, rh)
}

func main() {
	r := new(Handle)
	r.GET([]byte("/a/:/b/:"), func(ctx *fasthttp.RequestCtx) {
		params := ctx.UserValue(Params).([]string)
		fmt.Fprintf(ctx, "params: %q", params)
	})
	r.POST([]byte("/a/:/b/:"), func(ctx *fasthttp.RequestCtx) {
		params := ctx.UserValue(Params).([]string)
		fmt.Fprintf(ctx, "params: %q\nbody: %s", params, ctx.PostBody())
	})
	fasthttp.ListenAndServe(":8080", r.Handler)
}
