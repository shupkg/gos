package gos

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

func (m *Mux) SetBasicAuth(auth M) {
	m.getAuth = auth
}

func basicAuth(users M) func(w http.ResponseWriter, req *http.Request) bool {
	return func(w http.ResponseWriter, req *http.Request) bool {
		if len(users) == 0 {
			return false
		}

		auth := req.Header.Get("Authorization")
		if auth == "" {
			return requestBasicAuth(w, "")
		}

		auths := strings.SplitN(auth, " ", 2)
		if len(auths) != 2 {
			return requestBasicAuth(w, "错误的 Authorization 头: %s", auth)
		}

		authMethod := auths[0]
		authB64 := auths[1]
		switch authMethod {
		case "Basic":
			authBytes, err := base64.StdEncoding.DecodeString(authB64)
			if err != nil {
				return requestBasicAuth(w, "DecodeString: %s", authB64)
			}
			userPwd := strings.SplitN(string(authBytes), ":", 2)
			if len(userPwd) != 2 {
				w.WriteHeader(http.StatusForbidden)
				return requestBasicAuth(w, "SplitN: %s", string(authBytes))
			}
			username, password := userPwd[0], userPwd[1]
			if pwd, find := users[username]; find {
				if pwd == password {
					return false
				}
			}
			return requestBasicAuth(w, "password error: %s,%s", username, password)
		default:
			return requestBasicAuth(w, "错误的验证方法: %s", authMethod)
		}
	}
}

//func forbidden(w http.ResponseWriter, format string, args ...interface{}) bool {
//	w.WriteHeader(http.StatusForbidden)
//	fmt.Fprintf(w, format, args...)
//	return true
//}

func requestBasicAuth(w http.ResponseWriter, format string, args ...interface{}) bool {
	w.Header().Set("WWW-Authenticate", `Basic realm="Login"`)
	w.WriteHeader(http.StatusUnauthorized)
	if len(format) > 0 {
		log.Printf(format+"\n", args...)
	}
	return true
}
