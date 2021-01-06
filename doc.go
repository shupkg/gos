package gos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"text/template"

	"github.com/russross/blackfriday/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var mdTpl *template.Template

func init() {
	mdTpl = template.Must(
		template.New("").Funcs(
			template.FuncMap{
				"json": func(value interface{}) string {
					d, _ := json.MarshalIndent(value, "", "  ")
					return string(d)
				},
				"form": func(value interface{}) string {
					s := formEncode(value)
					return strings.ReplaceAll(s, "&", "\n&")
				},
			},
		).Parse(string(DocGomd.Bytes())),
	)
}

type Doc struct {
	Name        string //名称
	Description string //描述
	Path        string //路径

	Errors map[int]M
	Input  interface{}
	Output interface{}

	InputParams []DocParamItem

	parsed bool
	sync.Mutex
}

type DocParamItem struct {
	Name        string //名称
	Type        string //类型
	Description string //描述
	Required    bool   //是否必须
	Example     string //Example
}

type Empty struct {
	Code string `json:"code,omitempty"`
}

func (ad *Doc) parse() {
	ad.Lock()
	defer ad.Unlock()
	if ad.parsed {
		return
	}
	ad.parsed = true

	if ad.Input == nil || len(ad.InputParams) > 0 {
		return
	}

	ad.InputParams = parseParams(reflect.ValueOf(ad.Input))
}

func (ad *Doc) Markdown(path string) (md []byte, html string) {
	ad.parse()
	ad.Path = path
	var buf bytes.Buffer
	if err := mdTpl.Execute(&buf, ad); err != nil {
		return []byte(fmt.Sprintf("`%s`", err.Error())), err.Error()
	}

	md = buf.Bytes()
	html = strings.NewReplacer(
		"{{ .Title }}", ad.Name,
		"{{ .Body }}", string(blackfriday.Run(md)),
	).Replace(string(DocHtml.Bytes()))

	return
}

func parseParams(vRef reflect.Value) []DocParamItem {
	var params []DocParamItem
	tRef := vRef.Type()
	for i := 0; i < vRef.NumField(); i++ {
		fieldT := tRef.Field(i)
		fieldV := vRef.Field(i)

		if fieldT.Anonymous {
			params = append(params, parseParams(fieldV)...)
			continue
		}

		var item DocParamItem
		tag := fieldT.Tag.Get("doc")
		item.Description = tag

		jTag := fieldT.Tag.Get("json")
		if jTag == "-" {
			continue
		}
		if jTags := strings.SplitN(jTag, ",", 2); len(jTags) > 0 {
			item.Name = strings.TrimSpace(jTags[0])
		}

		if item.Name == "" {
			item.Name = fieldT.Name
		}

		vTag := fieldT.Tag.Get("validate")
		item.Required = strings.Contains(vTag, "required") || !strings.Contains(jTag, ",omitempty")

		v, _ := json.Marshal(fieldV.Interface())
		item.Example = string(v)

		item.Type = fieldV.Kind().String()

		params = append(params, item)
	}

	return params
}

func wrapDoc(v interface{}) *Doc {
	switch o := v.(type) {
	case func() *Doc:
		return o()
	case func(*Doc):
		docH := &Doc{}
		o(docH)
		return docH
	case Doc:
		return &o
	case *Doc:
		return o
	}
	return nil
}

var gb18130 = simplifiedchinese.GB18030.NewEncoder()

type docNs []*Doc

func (g docNs) Len() int {
	return len(g)
}

func (g docNs) Less(i, j int) bool {
	a, _ := gb18130.Bytes([]byte(g[i].Name))
	b, _ := gb18130.Bytes([]byte(g[j].Name))
	l := len(b)
	for idx, chr := range a {
		if idx > l-1 {
			return false
		}
		if chr != b[idx] {
			return chr < b[idx]
		}
	}
	return true
}

func (g docNs) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}
