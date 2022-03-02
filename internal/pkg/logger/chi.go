package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
)

func NewStructuredLogger(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

type StructuredLogger struct {
	Logger *zap.SugaredLogger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: l.Logger}
	var logFields []interface{}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields = append(logFields, "req_id", reqID)
	}

	logFields = tracingFields(r, logFields)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields = append(logFields, "http_scheme", scheme)
	logFields = append(logFields, "http_proto", r.Proto)
	logFields = append(logFields, "http_method", r.Method)
	logFields = append(logFields, "uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI))

	entry.Logger = entry.Logger.With(logFields...)

	return entry
}

func tracingFields(r *http.Request, logFields []interface{}) []interface{} {
	if span := opentracing.SpanFromContext(r.Context()); span != nil {
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			return append(logFields,
				zap.Stringer("trace_id", sc.TraceID()),
				zap.Stringer("span_id", sc.SpanID()),
			)
		}
	}
	return logFields
}

type StructuredLoggerEntry struct {
	Logger *zap.SugaredLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.With(
		"resp_status", status,
		"resp_bytes_length", bytes,
		"resp_elapsed_ms", float64(elapsed.Nanoseconds())/1000000.0,
	)

	l.Logger.Info("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.With(
		"stack", string(stack),
		"panic", fmt.Sprintf("%+v", v),
	)
}
