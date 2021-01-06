package gos

import (
	"fmt"
	"net/http"
)

const (
	keyFieldHttpStatus   = "__http_status"
	keyFieldHttpRedirect = "__redirect"

	keyFieldCode = "code"
	keyFieldMsg  = "msg"
	keyFieldData = "data"
)

func (m *Map) SetData(data interface{}) *Map {
	return m.Set(keyFieldData, data)
}

func (m *Map) SetCode(code string) *Map {
	if code != "" {
		return m.Set(keyFieldCode, code)
	}
	return m
}

func (m *Map) SetMessage(format string, args ...interface{}) *Map {
	if format != "" {
		return m.Set(keyFieldMsg, fmt.Sprintf(format, args...))
	}
	return m
}

func (m *Map) SetEx(err error) *Map {
	if err != nil {
		m.SetMessage(err.Error())
	}
	return m
}

func (m *Map) SetStatus(status int) *Map {
	if status > 0 {
		m.Set(keyFieldHttpStatus, status)
	}
	return m
}

func (m *Map) PopStatus(status *int) bool {
	v, find := m.Get(keyFieldHttpStatus)
	if find {
		m.Del(keyFieldHttpStatus)
		*status = v.(int)
	}
	return find
}

func (m *Map) Redirect(status *int, redirectTo *string) bool {
	m.PopStatus(status)
	if *status != http.StatusMovedPermanently && *status != http.StatusTemporaryRedirect {
		return false
	}

	if to, find := m.Get(keyFieldHttpRedirect); find {
		if sto, ok := to.(string); ok && sto != "" {
			m.Del(keyFieldHttpRedirect)
			*redirectTo = sto
			return true
		}
	}

	*status = 500
	return false
}

func (m *Map) Basic(code, format string, args ...interface{}) *Map {
	return m.SetCode(code).SetCode(fmt.Sprintf(format, args...))
}

func (m *Map) Error() string {
	var (
		code, _ = m.Get(keyFieldCode)
		msg, _  = m.Get(keyFieldMsg)
	)
	return fmt.Sprintf("[%v] %v", code, msg)
}
