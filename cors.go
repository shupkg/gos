package gos

import (
	"net/http"
	"strings"
)

func cors(w http.ResponseWriter, req *http.Request) bool {
	var (
		headers = w.Header()
		origin  = req.Header.Get("Origin")
	)

	if req.Method == http.MethodOptions {
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")

		headers.Set("Access-Control-Allow-Origin", origin)
		headers.Set("Access-Control-Allow-Credentials", "true")
		headers.Set("Access-Control-Max-Age", "86400")

		reqMethod := strings.ToUpper(req.Header.Get("Access-Control-Request-Method"))
		if reqMethod != "" {
			headers.Set("Access-Control-Allow-Methods", reqMethod)
		}

		if reqHeaders := parseHeaderList(req.Header.Get("Access-Control-Request-Headers")); len(reqHeaders) > 0 {
			headers.Set("Access-Control-Allow-Headers", strings.Join(reqHeaders, ", "))
		}

		w.WriteHeader(http.StatusNoContent)

		return true
	}

	headers.Add("Vary", "Origin")
	headers.Set("Access-Control-Allow-Origin", origin)
	headers.Set("Access-Control-Allow-Credentials", "true")
	return false
}

func parseHeaderList(headerList string) []string {
	l := len(headerList)
	h := make([]byte, 0, l)
	upper := true
	// Estimate the number headers in order to allocate the right splice size
	t := 0
	for i := 0; i < l; i++ {
		if headerList[i] == ',' {
			t++
		}
	}
	headers := make([]string, 0, t)
	for i := 0; i < l; i++ {
		b := headerList[i]
		switch {
		case b >= 'a' && b <= 'z':
			if upper {
				h = append(h, b-('a'-'A'))
			} else {
				h = append(h, b)
			}
		case b >= 'A' && b <= 'Z':
			if !upper {
				h = append(h, b+('a'-'A'))
			} else {
				h = append(h, b)
			}
		case b == '-' || b == '_' || b == '.' || (b >= '0' && b <= '9'):
			h = append(h, b)
		}

		if b == ' ' || b == ',' || i == l-1 {
			if len(h) > 0 {
				// Flush the found header
				headers = append(headers, string(h))
				h = h[:0]
				upper = true
			}
		} else {
			upper = b == '-' || b == '_'
		}
	}
	return headers
}
