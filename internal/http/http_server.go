// Package http is responsible for creating an HTTP server
package http

import (
	"fmt"
	"net"
	"net/http"
)

// Server is an http.Server and its listener
type Server struct {
	httpServer *http.Server
	listener   net.Listener
}

// NewHTTPServer starts a logged http server on the given address
func NewHTTPServer(address string, port int, router http.Handler) (*Server, error) {
	// Create the logged handler
	loggedHandler := logHandler(router)

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

// Close will shutdown the listener and the http server
func (srv *Server) Close() error {
	srv.listener.Close()
	srv.httpServer.Close()
	return nil
}

// Port of the listener
func (srv *Server) Port() int {
	return srv.listener.Addr().(*net.TCPAddr).Port
}
