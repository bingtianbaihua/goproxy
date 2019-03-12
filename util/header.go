package util

import "net/http"

var (
	hopHeaders = []string{
		"Connection",
		"Proxy-Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailer",
		"Transfer-Encoding",
		"Upgrade",
	}
)

func CopyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func ClearHeaders(headers http.Header) {
	for key, _ := range headers {
		headers.Del(key)
	}
}

func RemoveProxyHeaders(req *http.Request) {
	req.RequestURI = ""
	for _, h := range hopHeaders {
		req.Header.Del(h)
	}
}
