package gos

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	tzh "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate = validator.New()
	zhTrans  = ut.New(zh.New()).GetFallback()
)

func init() {
	tzh.RegisterDefaultTranslations(validate, zhTrans)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
			if name == "-" {
				return ""
			}
		}

		return name
	})
}

func isValidateErrorHandled(w http.ResponseWriter, err error, status *int) (handled bool) {
	var vErrs validator.ValidationErrors
	if errors.As(err, &vErrs) {
		var fields = M{}
		for _, ex := range vErrs {
			fields = fields.Set(ex.Field(), ex.Translate(zhTrans))
		}

		*status = 400
		Render(w,
			M{
				"code": "PARAM_ERROR",
				"msg":  "输入参数错误",
				"data": fields,
			},
			status,
		)
		return true
	}

	return false
}
