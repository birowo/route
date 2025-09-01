package main

import (
	"fmt"

	fh "github.com/birowo/route/fasthttp"
	"github.com/valyala/fasthttp"
)

func main() {
	r := new(fh.Handle)
	r.GET([]byte("/a/:/b/:"), func(ctx *fasthttp.RequestCtx) {
		params := ctx.UserValue(fh.Params).([]string)
		fmt.Fprintf(ctx, "params: %q", params)
	})
	r.POST([]byte("/a/:/b/:"), func(ctx *fasthttp.RequestCtx) {
		params := ctx.UserValue(fh.Params).([]string)
		fmt.Fprintf(ctx, "params: %q\nbody: %s", params, ctx.PostBody())
	})
	fasthttp.ListenAndServe(":8080", r.Handler)
}
