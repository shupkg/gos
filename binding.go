package gos

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
)

var (
	formDecoder *form.Decoder
	formEncoder *form.Encoder
)

func init() {
	formDecoder = form.NewDecoder()
	formDecoder.SetTagName("json")
	formEncoder = form.NewEncoder()
	formEncoder.SetTagName("json")
}

func (c *Context) Bind(value interface{}) error {
	contentType := strings.TrimSpace(c.Header.Get("Content-Type"))
	if idx := strings.Index(contentType, ";"); idx > 0 {
		contentType = strings.ToLower(strings.TrimSpace(contentType[:idx]))
	}

	var err error
	switch contentType {
	case MimeJSON:
		var v []byte
		if v, err = ioutil.ReadAll(c.Body); err == nil {
			err = json.Unmarshal(v, value)
		}
		if err != nil {
			return c.MapBad("BIND_ERROR", "解析参数错误: %v", err.Error()).SetData(M{})
		}
	case MimeFormUrlencoded:
		if err = c.ParseForm(); err == nil {
			err = formDecoder.Decode(value, c.Form)
		} else {
			return c.MapBad("PARSE_FORM", err.Error())
		}
	case MimeFormData:
		if err = c.ParseMultipartForm(1 << 20); err == nil {
			err = formDecoder.Decode(value, c.MultipartForm.Value)
		} else {
			return c.MapBad("PARSE_FILE", err.Error())
		}
	default:
		return c.MapBad("BIND_ERROR", "不支持的Content-Type: %q", contentType)
	}

	if err != nil {
		return err
	}
	if err := validate.Struct(value); err != nil {
		return err
	}
	return nil
}

func formEncode(value interface{}) string {
	uv, _ := formEncoder.Encode(value)
	if len(uv) > 0 {
		return uv.Encode()
	}
	return ""
}

func isBindErrorHandled(w http.ResponseWriter, err error, status *int) (handled bool) {
	var dErrs form.DecodeErrors
	if errors.As(err, &dErrs) {
		var fields = M{}
		for field, ex := range dErrs {
			fields = fields.Set(field, ex.Error())
		}
		*status = 400
		Render(w,
			M{
				"code": "BIND_ERROR",
				"msg":  "解析参数错误",
				"data": fields,
			},
			status,
		)
		return true
	}

	return false
}
