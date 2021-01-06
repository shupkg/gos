package gos

import (
	"encoding/json"
	"net/http"
)

func Render(w http.ResponseWriter, data interface{}, status *int) {
	if *status == 0 {
		*status = 200
	}

	if data == nil {
		data = MapOK()
	}

	var mime = MimeHTML
	var body []byte

	switch o := data.(type) {
	case string:
		body = []byte(o)
	case []byte:
		body = o
	default:
		mime = MimeJSON
		var err error
		body, err = json.Marshal(data)
		if err != nil {
			*status = 500
			body, _ = NewMap().Basic("JSON_ENCODE", err.Error()).MarshalBinary()
		}
	}

	w.WriteHeader(*status)
	w.Header().Set(HeaderContentType, mime)
	w.Write(body)
}
