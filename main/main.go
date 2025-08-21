package main

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/birowo/route"
)

func main() {
	type Itm struct {
		pathSet []byte
		v       int
		pathGet []byte
		prms    string
	}

	tbl := [...]Itm{
		{[]byte("/a/:/b/:/c"), 1, []byte("/a/A1/b/B1/c"), `["A1" "B1"]`},
		{[]byte("/a/:/b/:/c/:"), 2, []byte("/a/A2/b/B2/c/C2"), `["A2" "B2" "C2"]`},
		{[]byte("/a/:/b/c/:"), 3, []byte("/a/A3/b/c/C3"), `["A3" "C3"]`},
		{[]byte("/a/:/b/:"), 4, []byte("/a/A4/b/B4"), `["A4" "B4"]`},
		{[]byte("/a/:/b/:/d"), 5, []byte("/a/A5/b/B5/d"), `["A5" "B5"]`},
	}
	tree := &route.Nodes[int]{}
	for _, itm := range tbl {
		tree.Set(itm.pathSet, itm.v)
	}
	for _, itm := range tbl {
		v, prms, prmsIdx := tree.Get(itm.pathGet)
		if v != itm.v || fmt.Sprintf("%q", prms[:prmsIdx]) != itm.prms {
			println("FAIL")
			return
		}
		fmt.Printf(
			"path set: %s\nvalue: %d\npath get: %s\nparams: %q\n",
			itm.pathSet, v, itm.pathGet, prms[:prmsIdx],
		)
	}
	bs, _ := json.MarshalIndent(tree, "", "  ")
	bs = regexp.MustCompile(
		`{\s+"Chr": 0,\s+"Str"\: "",\s+"V"\: 0,\s+"Nds"\: null\s+},*\s+`,
	).ReplaceAll(bs, nil)
	println("tree:", string(bs))
}
