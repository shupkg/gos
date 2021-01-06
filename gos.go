package gos

import (
	"net/http"
	"time"
)

type HandlerFunc func(*Context) (interface{}, error)

func NewServices(services ...interface{}) *Mux {
	mux := New()
	mux.HandleServices(services...)
	return mux
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/favicon.ico" {
		return
	}

	if cors(w, req) {
		return
	}

	c := getContext(m, req)
	defer c.drop()

	var (
		data    interface{}
		err     error
		status  int
		startAt = time.Now()
	)

	finished := func() bool {
		defer func() {
			if re := recover(); re != nil {
				m.Printf("RECOVER: %v\n", re)
				err = MapUnhandled("RECOVER", "%v", re)
			}
		}()

		res := m.trie.Match(c.URL.Path)
		c.params = res.Params
		if res.Node == nil {
			if res.TSR != "" || res.FPR != "" {
				c.URL.Path = res.TSR
				if res.FPR != "" {
					c.URL.Path = res.FPR
				}
				code := http.StatusMovedPermanently
				if c.Method != http.MethodGet {
					code = http.StatusTemporaryRedirect
				}
				http.Redirect(w, req, c.URL.String(), code)
				return true
			}
			http.NotFound(w, req)
			return true
		} else {
			if handler, ok := res.Node.GetHandler(c.Method).(HandlerFunc); !ok {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return true
			} else {
				if req.Method == http.MethodGet && basicAuth(m.getAuth)(w, req) {
					return true
				}
				data, err = handler(c)
				return false
			}
		}
	}()

	if finished {
		return
	}

	if err != nil {
		handledError(w, c, err, &status)
	} else {
		Render(w, data, &status)
	}

	t := time.Now().Sub(startAt)
	if t > time.Millisecond*100 || (status >= 400 && status != 404) {
		m.Printf("| %-3d | %-7s | %-10s | %s\n", status, c.Method, t.String(), c.AbsURL())
	}
}

func handledError(w http.ResponseWriter, c *Context, err error, status *int) {
	if isBindErrorHandled(w, err, status) {
		return
	}

	if isValidateErrorHandled(w, err, status) {
		return
	}

	if IsMapErrorHandled(w, c.Request, err, status) {
		return
	}

	*status = 500
	Render(w, M{"code": "ERROR", "msg": err.Error()}, status)
}

func redirect(status int, redirectTo string) HandlerFunc {
	return func(c *Context) (interface{}, error) {
		return nil, Redirect(status, redirectTo)
	}
}

func statusText(status int, code string, m ...M) HandlerFunc {
	return func(c *Context) (interface{}, error) {
		return nil, MapStatus(status).SetCode(code).Merge(m...)
	}
}
