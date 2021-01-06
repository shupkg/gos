package gos

import (
	"net/http"
	"sync"
)

var ctxPool = sync.Pool{New: func() interface{} {
	return &Context{}
}}

func getContext(mux *Mux, req *http.Request) *Context {
	ctx := ctxPool.Get().(*Context)
	ctx.Request = req
	ctx.mux = mux
	ctx.Printer = mux.Printer
	return ctx
}

func (c *Context) drop() {
	c.Request = nil
	c.params = nil
	c.mux = nil
	c.Printer = nil
	ctxPool.Put(c)
}
