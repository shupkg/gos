package gos

import (
	"strconv"
)

func (c *Context) GetFormInt(name string) int64 {
	s := c.FormValue(name)
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func (c *Context) GetFormBool(name string) bool {
	s := c.FormValue(name)
	if s == "" {
		return false
	}
	i, _ := strconv.ParseBool(s)
	return i
}

func (c *Context) GetFormFloat(name string) float64 {
	s := c.FormValue(name)
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseFloat(s, 64)
	return i
}
