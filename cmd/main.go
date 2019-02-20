package main

import (
	"github.com/bingtianbaihua/goproxy/proxy"

	"github.com/bingtianbaihua/goproxy/log"
)

func main() {
	ps := proxy.NewProxyServer()

	log.Info("begin proxy")
	log.Error("proxy exit %v", ps.ListenAndServe())
}
