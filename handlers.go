package gos

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var (
	regexpNoChar = regexp.MustCompile(`([\W]+)`)
	regexpUpChar = regexp.MustCompile(`([A-Z]+)`)
)

func (m *Mux) HandleService(pathPrefix string, svc interface{}) {
	v := reflect.ValueOf(svc)

	var prefix string
	if t := v.Type(); t.Kind() == reflect.Ptr {
		prefix = lowerPath(t.Elem().Name())
	} else {
		prefix = lowerPath(t.Name())
	}

	if nameMethod, ok := v.Type().MethodByName("Name"); ok {
		if nFunc, ok := v.Method(nameMethod.Index).Interface().(func() string); ok {
			prefix = nFunc()
		}
	}

	var docs []*Doc
	for i := 0; i < v.NumMethod(); i++ {
		n := v.Type().Method(i).Name
		if !unicode.IsUpper([]rune(n)[0]) {
			continue
		}

		method := v.Method(i).Interface()

		if handlerFunc := wrapHandler(method); handlerFunc != nil {
			m.Handle(http.MethodPost, joinPath(pathPrefix, prefix, lowerPath(n)), handlerFunc)
			continue
		}

		if doc := wrapDoc(method); doc != nil {
			doc.Path = n
			docs = append(docs, doc)
		}
	}

	if len(docs) > 0 {
		indexPath := joinPath(pathPrefix, prefix)

		var links []string
		sort.Sort(docNs(docs))
		for _, doc := range docs {
			n := lowerPath(doc.Path, "Doc")
			fullPath := joinPath(pathPrefix, prefix, n)
			relPath := filepath.Base(indexPath) + "/" + n
			links = append(links,
				fmt.Sprintf("<dt>· %s</dt>\n\t\t<dd><a href=%q>%s</a></dd>", doc.Name, relPath, fullPath),
			)

			md, html := doc.Markdown(fullPath)
			m.Handle(http.MethodGet, fullPath, func(c *Context) (interface{}, error) {
				output := c.FormValue("output")
				switch output {
				case "md":
					return md, nil
				default:
					return html, nil
				}
			})
		}

		s := []byte(
			strings.NewReplacer(
				"{{ .Title }}", "文档",
				"{{ .Links }}", strings.Join(links, "\n\t\t"),
			).Replace(string(DocIndexGohtml.Bytes())),
		)
		m.Handle(http.MethodGet, indexPath, func(c *Context) (interface{}, error) {
			return s, nil
		})
	}
}

func (m *Mux) HandleServices(ss ...interface{}) {
	for _, svc := range ss {
		m.HandleService("", svc)
	}
}

func joinPath(s ...string) string {
	var n []string
	for _, ss := range s {
		if ss = strings.Trim(strings.TrimSpace(ss), "/"); ss != "" {
			n = append(n, ss)
		}
	}
	return "/" + strings.Join(n, "/")
}

func lowerPath(name string, trims ...string) string {
	if len(name) == 0 {
		return ""
	}
	for _, trim := range trims {
		name = strings.TrimSuffix(strings.TrimPrefix(name, trim), trim)
	}
	name = regexpNoChar.ReplaceAllString(name, "_")
	name = regexpUpChar.ReplaceAllString(name, "_$1")
	name = strings.Trim(strings.ToLower(name), "_")
	if name[0] >= '0' && name[0] <= '9' {
		name = "_" + name
	}
	return name
}

func wrapHandler(hFunc interface{}) HandlerFunc {
	switch hFunc := hFunc.(type) {
	case func(*Context) (interface{}, error):
		return hFunc
	case func(*Context) error:
		return func(c *Context) (interface{}, error) {
			return nil, hFunc(c)
		}
	default:
		return nil
	}
}
