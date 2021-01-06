package gos

import (
	"fmt"
	"net/http"

	"github.com/tomasen/realip"
)

type Context struct {
	*http.Request

	mux    *Mux
	params map[string]string

	//data   interface{}
	//headerSets http.Header
	//httpStatus int

	Printer
}

func (c *Context) RealIp() string {
	return realip.FromRequest(c.Request)
}

func (c *Context) IsTLS() bool {
	return c.TLS != nil
}

func (c *Context) Scheme() string {
	if c.IsTLS() {
		return "https"
	}
	if scheme := c.Header.Get(HeaderXForwardedProto); scheme != "" {
		return scheme
	}
	if scheme := c.Header.Get(HeaderXForwardedProtocol); scheme != "" {
		return scheme
	}
	if ssl := c.Header.Get(HeaderXForwardedSsl); ssl == "on" {
		return "https"
	}
	if scheme := c.Header.Get(HeaderXUrlScheme); scheme != "" {
		return scheme
	}
	return "http"
}

func (c *Context) AbsURL() string {
	return fmt.Sprintf("%s://%s%s", c.Scheme(), c.Host, c.URL.String())
}

func (c *Context) MapBad(code, format string, args ...interface{}) *Map {
	return MapStatus(http.StatusBadRequest).Basic(code, format, args...)
}

func (c *Context) MapUnhandled(code, format string, args ...interface{}) *Map {
	return MapStatus(http.StatusInternalServerError).Basic(code, format, args...)
}
