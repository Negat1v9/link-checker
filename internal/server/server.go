package server

import (
	"context"
	"net/http"

	"github.com/Negat1v9/link-checker/config"
)

type Server struct {
	cfg        *config.ServerCfg
	httpServer http.Server
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: &cfg.ServerCfg,
		httpServer: http.Server{
			Addr: cfg.ListedAddr,
		},
	}
}

func (s *Server) Run() error {

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
