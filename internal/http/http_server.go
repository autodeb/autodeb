// Package http is responsible for creating an HTTP server
package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/log"
)

// Server is an http.Server and its listener
type Server struct {
	httpServer *http.Server
	listener   net.Listener
}

// NewHTTPServer starts a logged http server on the given address
func NewHTTPServer(address string, port int, router http.Handler, logger log.Logger) (*Server, error) {
	// Create the logged handler
	loggedHandler := logHandler(router, logger)

	listenAddress := fmt.Sprintf("%s:%d", address, port)

	// Listen on the given address
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return nil, err
	}

	// Create an http server with the logged handler
	httpServer := &http.Server{
		Handler: loggedHandler,
	}

	server := Server{
		httpServer: httpServer,
		listener:   listener,
	}

	go httpServer.Serve(listener)

	return &server, nil
}

// Shutdown will close the listener and the shutdown http server
func (srv *Server) Shutdown(ctx context.Context) error {
	srv.listener.Close()
	return srv.httpServer.Shutdown(ctx)
}

// Port of the listener
func (srv *Server) Port() int {
	return srv.listener.Addr().(*net.TCPAddr).Port
}
