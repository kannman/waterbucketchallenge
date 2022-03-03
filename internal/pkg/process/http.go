package process

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/kannman/waterbucketchallenge/internal/pkg/logger"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type httpServerRunner struct {
	listener        net.Listener
	router          chi.Router
	port            int
	name            string
	gracefulTimeout time.Duration
}

func NewHTTPServerRunner(router chi.Router, name string, port int, gracefulTimeout time.Duration) (func() error, func(error), error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, nil, errors.Wrap(err, "http "+name+" listener init failure")
	}

	r := &httpServerRunner{
		listener:        listener,
		router:          router,
		port:            port,
		name:            name,
		gracefulTimeout: gracefulTimeout,
	}

	httpServer := &http.Server{Handler: r.router}
	return func() error {
			logger.Warnf("http "+r.name+" server: started on %d port", r.port)
			return errors.Wrap(httpServer.Serve(r.listener), "http server")
		}, func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), r.gracefulTimeout)
			defer cancel()

			httpServer.SetKeepAlivesEnabled(false)
			if err := errors.Wrap(httpServer.Shutdown(ctx), "http "+r.name+" server: error during shutdown"); err != nil {
				logger.Error("http "+r.name+" server: stop failure", err)
			} else {
				logger.Warn("http " + r.name + " server: gracefully stopped")
			}
		}, nil
	//return r, nil
}
