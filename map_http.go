package gos

import (
	"errors"
	"net/http"
)

func NewMap(maps ...M) *Map {
	return (&Map{M: M{}}).Merge(maps...)
}

func MapOK() M {
	return M{keyFieldCode: "OK"}
}

func MapStatus(status int) *Map {
	return NewMap().SetStatus(status)
}

//输入错误
func MapBad(code, format string, args ...interface{}) *Map {
	return MapStatus(http.StatusBadRequest).Basic(code, format, args...)
}

//已验证，无权操作
func MapForbidden(code, format string, args ...interface{}) *Map {
	return MapStatus(http.StatusForbidden).Basic(code, format, args...)
}

//未验证，需要登录
func MapUnauthorized(code, format string, args ...interface{}) *Map {
	return MapStatus(http.StatusUnauthorized).Basic(code, format, args...)
}

//未处理的错误
func MapUnhandled(code, format string, args ...interface{}) *Map {
	return MapStatus(http.StatusInternalServerError).Basic(code, format, args...)
}

func Redirect(status int, redirectTo string) *Map {
	return MapStatus(status).Set(keyFieldHttpRedirect, redirectTo)
}

func IsMapErrorHandled(w http.ResponseWriter, req *http.Request, err error, status *int) (handled bool) {
	var m *Map
	if errors.As(err, &m) {
		var redirectTo string
		if m.Redirect(status, &redirectTo) {
			http.Redirect(w, req, redirectTo, *status)
			return true
		}
		Render(w, m, status)
		return true
	}
	return false
}
