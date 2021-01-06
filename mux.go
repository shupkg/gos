package gos

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Params represents named parameter values
type Params map[string]string

// HandlerFunc is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (variables).
//type HandlerFunc func(http.ResponseWriter, *http.Request, Params)

// Mux is a tire base HTTP request router which can be used to
// dispatch requests to different handler functions.
type Mux struct {
	trie *trie
	Printer
	getAuth M
}

// New returns a Mux instance.
func New(opts ...Options) *Mux {
	opt := defaultOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return &Mux{
		Printer: log.New(os.Stderr, "", log.LstdFlags),
		trie: &trie{
			ignoreCase: opt.IgnoreCase,
			fpr:        opt.FixedPathRedirect,
			tsr:        opt.TrailingSlashRedirect,
			root: &Node{
				parent:   nil,
				children: make(map[string]*Node),
				handlers: make(map[string]interface{}),
			},
		},
	}
}

// Handle registers a new handler with method and path in the Mux.
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (m *Mux) Handle(method, pattern string, handler HandlerFunc) {
	if method == "" {
		panic(fmt.Errorf("invalid method"))
	}
	m.trie.Define(pattern).Handle(strings.ToUpper(method), handler)
}

func (m *Mux) findHandler(c *Context) HandlerFunc {
	path := c.URL.Path
	method := c.Method
	res := m.trie.Match(path)
	c.params = res.Params
	if res.Node == nil {
		// FixedPathRedirect or TrailingSlashRedirect
		if res.TSR != "" || res.FPR != "" {
			c.URL.Path = res.TSR
			if res.FPR != "" {
				c.URL.Path = res.FPR
			}
			code := http.StatusMovedPermanently
			if method != "GET" {
				code = http.StatusTemporaryRedirect
			}
			return redirect(code, c.URL.String())
		}
		return statusText(404, "NotFound", M{"url": c.AbsURL()})
	} else {
		if handler, ok := res.Node.GetHandler(method).(HandlerFunc); !ok {
			return statusText(405, "NotAllowed", M{"url": c.AbsURL(), "method": c.Method})
		} else {
			return handler
		}
	}
}
