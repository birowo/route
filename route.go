package route

import (
	"bytes"
	"unsafe"
)

const chrStSz = 126

type N[T any] struct {
	Chr byte
	Str string
	V   T
	Nds *Nodes[T]
}
type Nodes[T any] [chrStSz]N[T]

const prmsIdxMax = 64

type Params [prmsIdxMax][]byte

func StrBs(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
func params[T any](str []byte, v T) (str_ string, ns *Nodes[T], v_ T) {
	j := len(str)
	for i := j - 1; i > -1; i-- {
		if str[i] == ':' {
			ns = &Nodes[T]{':': {':', StrBs(str[i+1 : j]), v, ns}}
			j = i
			v = v_
		}
	}
	str_ = StrBs(str[:j])
	v_ = v
	return
}
func (ns *Nodes[T]) Set(str []byte, v T) {
	strLen := len(str)
	prmsIdx := 0
	for i := range strLen {
		if str[i] == ':' {
			prmsIdx++
		}
	}
	if strLen == 0 || prmsIdx > prmsIdxMax {
		return
	}
	i := 0
	for i < strLen {
		chr := str[i]
		i++
		if chr != ns[chr].Chr {
			ns[chr].Chr = chr
			ns[chr].Str, ns[chr].Nds, ns[chr].V = params(str[i:], v)
			return
		}
		nsChrStrLen := len(ns[chr].Str)
		j := 0
		for i < strLen && j < nsChrStrLen && str[i] == ns[chr].Str[j] {
			i++
			j++
		}
		if i < strLen && j < nsChrStrLen {
			var ns_ Nodes[T]
			chr_ := str[i]
			ns_[chr_].Chr = chr_
			ns_[chr_].Str, ns_[chr_].Nds, ns_[chr_].V = params(str[i+1:], v)
			chr_ = ns[chr].Str[j]
			ns_[chr_] = N[T]{chr_, ns[chr].Str[j+1:], ns[chr].V, ns[chr].Nds}
			ns[chr].Str = ns[chr].Str[:j]
			var v_ T
			ns[chr].V = v_
			ns[chr].Nds = &ns_
			return
		}
		if j < nsChrStrLen {
			var ns_ Nodes[T]
			chr_ := ns[chr].Str[j]
			ns_[chr_] = N[T]{chr_, ns[chr].Str[j+1:], ns[chr].V, ns[chr].Nds}
			ns[chr].Str = ns[chr].Str[:j]
			ns[chr].V = v
			ns[chr].Nds = &ns_
			return
		}
		if i < strLen {
			if ns[chr].Nds == nil {
				ns[chr].Nds = &Nodes[T]{}
			}
		}
		ns = ns[chr].Nds
	}
}
func (ns *Nodes[T]) Get(str []byte) (v T, prms Params, prmsIdx int) {
	strLen := len(str)
	if strLen == 0 {
		return
	}
	i := 0
	for i < strLen && ns != nil {
		chr := str[i]
		if ns[chr].Chr == chr {
			i++
			nsChrStr := []byte(ns[chr].Str)
			end := i + len(nsChrStr)
			if end > strLen {
				return
			}
			if bytes.Equal(str[i:end], nsChrStr) {
				if end == strLen {
					v = ns[chr].V
					return
				}
				ns = ns[chr].Nds
			}
			i = end
		} else if ns[':'].Chr == ':' {
			j := i
			for i < strLen && str[i] != '/' {
				i++
			}
			prms[prmsIdx] = str[j:i]
			prmsIdx++
			nsChrStr := []byte(ns[':'].Str)
			end := i + len(nsChrStr)
			if end > strLen {
				return
			}
			if bytes.Equal(str[i:end], nsChrStr) {
				if end == strLen {
					v = ns[':'].V
					return
				}
				ns = ns[':'].Nds
			}
			i = end
		}
	}
	return
}
