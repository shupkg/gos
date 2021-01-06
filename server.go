package gos

import (
	"net/http"
)

type Server interface {
	Start(handler http.Handler) error
}

func (m *Mux) Start(addr string) error {
	m.Printf("监听地址:", addr)
	return http.ListenAndServe(addr, m)
}

func (m *Mux) StartWith(s Server) error {
	return s.Start(m)
}
