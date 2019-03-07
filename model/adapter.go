package model

import "net/http"

type NextHandler func(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request))

func BuildFunc(h func(http.ResponseWriter, *http.Request), next ...NextHandler) func(http.ResponseWriter, *http.Request) {
	if len(next) == 0 {
		return h
	}
	return BuildFunc(func(w http.ResponseWriter, req *http.Request) {
		next[0](w, req, h)
	}, next[1:]...)
}

func Build(h http.Handler, next ...NextHandler) http.Handler {
	return http.HandlerFunc(BuildFunc(h.ServeHTTP, next...))
}
