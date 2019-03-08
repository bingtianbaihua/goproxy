package server

import (
	"net"
	"net/http"

	"github.com/bingtianbaihua/goproxy/config"
	"github.com/bingtianbaihua/goproxy/middleware"
	"github.com/bingtianbaihua/goproxy/model"
)

type ProxyServer struct {
	serv *http.Server
}

func NewProxyServer(cfg *config.Config) (*ProxyServer, error) {
	recover, _ := middleware.NewRecoverAdapter()
	auth, _ := middleware.NewAuthAdapter(cfg)
	proxy, _ := middleware.NewProxyAdapter()

	// build handler chains
	chains := model.Build(proxy, auth.HandleTask, recover.HandleTask)

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: chains,
	}
	return &ProxyServer{
		serv: server,
	}, nil
}

func (s *ProxyServer) ListenAndServe() error {
	return s.serv.ListenAndServe()
}
