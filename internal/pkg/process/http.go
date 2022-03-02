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

//
//func (r *httpServerRunner) Actor() (func() error, func(error)) {
//	httpServer := &http.Server{Handler: r.router}
//	return func() error {
//			logger.Warnf("http "+r.name+" server: started on %d port", r.port)
//			return errors.Wrap(httpServer.Serve(r.listener), "http server")
//		}, func(err error) {
//			ctx, cancel := context.WithTimeout(context.Background(), r.gracefulTimeout)
//			defer cancel()
//
//			httpServer.SetKeepAlivesEnabled(false)
//			if err := errors.Wrap(httpServer.Shutdown(ctx), "http "+r.name+" server: error during shutdown"); err != nil {
//				logger.Error("http "+r.name+" server: stop failure", err)
//			} else {
//				logger.Warn("http " + r.name + " server: gracefully stopped")
//			}
//		}
//}

//
//
//type patternType int
//
//var patternKey patternType
//
//func setURLPatternMiddleware(handler http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		newCtx := context.WithValue(r.Context(), patternKey, patternFunc(r))
//		handler.ServeHTTP(w, r.WithContext(newCtx))
//	})
//}
//
//func setRequestURLMiddleware(handler http.Handler) http.Handler {
//	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
//		if request.TLS == nil {
//			request.URL.Scheme = "http"
//		} else {
//			request.URL.Scheme = "https"
//		}
//		request.URL.Host = request.Host
//		handler.ServeHTTP(writer, request)
//	})
//}
//
//func patternFunc(r *http.Request) string {
//	rCtx := chi.RouteContext(r.Context())
//	path := rCtx.RoutePath
//	if path == "" {
//		if r.URL.RawPath != "" {
//			path = r.URL.RawPath
//		} else {
//			path = r.URL.Path
//		}
//	}
//	newCtx := chi.NewRouteContext()
//	if rCtx.Routes.Match(newCtx, r.Method, path) {
//		return newCtx.RoutePattern()
//	}
//	return ""
//}
//
//func opNameFunc(r *http.Request) string {
//	result := "HTTP " + r.Method
//	opName := r.Context().Value(patternKey).(string)
//	if opName != "" {
//		result += ": " + opName
//	}
//	return result
//}
//
//func setupSwaggerUI(swaggerJSON []byte, debugRunner *httpServerRunner, httpRouter chi.Router) {
//	if swaggerJSON == nil {
//		return
//	}
//	// add handlers to debug port
//	debugRunner.preRun = append(debugRunner.preRun, func() error {
//		swagger.AddJSONHandler(debugRunner.router, swaggerJSON)
//		swagger.AddSwaggerUIHandler(debugRunner.router)
//		// ducttape: pass requests from swagger-ui to public port
//		debugRunner.router.Mount("/", httpRouter)
//		return nil
//	})
//}
