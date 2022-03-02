package tracing

import (
	"io"
	"net/http"
	"os"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// Option is a function that sets some option on the client.
type Option = config.Option

// InitTracer uses environment variables to set the tracer's Configuration.
func InitTracer(options ...Option) (io.Closer, error) {
	if _, ok := os.LookupEnv("JAEGER_SERVICE_NAME"); !ok {
		_ = os.Setenv("JAEGER_SERVICE_NAME", "undefined")
	}
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, err
	}
	tracer, closer, err := cfg.NewTracer(options...)
	if err != nil {
		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	return closer, nil
}

func Middleware(handler http.Handler) http.Handler {
	return nethttp.MiddlewareFunc(opentracing.GlobalTracer(), handler.ServeHTTP)
}
