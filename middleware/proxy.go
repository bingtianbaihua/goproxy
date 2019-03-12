package middleware

import (
	"io"
	"net"
	"net/http"

	"github.com/bingtianbaihua/goproxy/util"

	"github.com/bingtianbaihua/goproxy/log"
)

var (
	proxyConnectSuccess = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")
)

type ProxyAdapter struct {
	Transport http.RoundTripper
}

func NewProxyAdapter() (*ProxyAdapter, error) {
	return &ProxyAdapter{
		Transport: http.DefaultTransport,
	}, nil
}

func (p *ProxyAdapter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "CONNECT" {
		p.HttpsHandler(rw, req)
	} else {
		p.HttpHandler(rw, req)
	}
}

func (p *ProxyAdapter) HttpHandler(rw http.ResponseWriter, req *http.Request) {
	util.RemoveProxyHeaders(req)

	resp, err := p.Transport.RoundTrip(req)
	if err != nil {
		log.Error("transport RoundTrip error: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	util.ClearHeaders(rw.Header())
	util.CopyHeaders(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)

	nr, err := io.Copy(rw, resp.Body)
	if err != nil && err != io.EOF {
		log.Error("got an error when copy remote response to client.%v", err)
		return
	}
	log.Info("copied %v bytes from remote host %v.", nr, req.URL.Host)
}

func (p *ProxyAdapter) HttpsHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info("[CONNECT] proxy tried to connect to remote host %v", req.URL.Host)

	hj, _ := rw.(http.Hijacker)
	client, _, err := hj.Hijack()
	if err != nil {
		log.Error("failed to get tcp connection of", req.RequestURI)
		http.Error(rw, "Failed", http.StatusBadRequest)
		return
	}

	remote, err := net.Dial("tcp", req.URL.Host)
	if err != nil {
		log.Error("proxy failed to connect %v", req.RequestURI)
		http.Error(rw, "Failed", http.StatusBadRequest)
		client.Close()
		return
	}

	client.Write(proxyConnectSuccess)

	go util.Pipe(remote, client)
	go util.Pipe(client, remote)
}
