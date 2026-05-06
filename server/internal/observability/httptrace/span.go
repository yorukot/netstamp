package httptrace

import "net/http"

func RequestSpanName(_ string, r *http.Request) string {
	return r.Method + " " + r.URL.Path
}
