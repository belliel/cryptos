package http

import (
	"context"
	"github.com/belliel/crypto-price-aggregator/internal/server"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	srv *http.Server
}

func NewServer(ctx context.Context, addr string, router chi.Router) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &Server{
		srv: srv,
	}
}

func (s *Server) Serve() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.srv.Shutdown(context.TODO())
}
