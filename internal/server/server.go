package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

var server *http.Server

func Run(host, port string, handler http.Handler) error {
	log.Printf("Starting server on %s:%s", host, port)

	server = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", host, port),
		Handler:           handler,
		ReadHeaderTimeout: 200 * time.Millisecond,
		ReadTimeout:       500 * time.Millisecond,
	}

	return server.ListenAndServe()
}

func Shutdown() error {
	return server.Shutdown(context.Background())
}
