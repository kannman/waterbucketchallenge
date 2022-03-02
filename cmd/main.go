package main

import (
	"net/http"
	"time"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/kannman/waterbucket/internal/app/controller"
	"github.com/kannman/waterbucket/internal/pkg/logger"
	"github.com/kannman/waterbucket/internal/pkg/process"
	"github.com/kannman/waterbucket/internal/pkg/tracing"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpswagger "github.com/swaggo/http-swagger"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap/zapcore"

	// import docs
	_ "github.com/kannman/waterbucket/api"
)

type Config struct {
	LogLevel                zapcore.Level `envconfig:"LOG_LEVEL" default:"debug"`
	HTTPPort                int           `envconfig:"HTTP_PORT" default:"8000"`
	DebugHTTPPort           int           `envconfig:"DEBUG_HTTP_PORT" default:"8001"`
	GracefulShutdownTimeout time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT" default:"15s"`
}

// @title Water bucket challenge
// @version 1.0
// @description This is a water bucket challenge service
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email swagger@github.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func main() {
	var cfg Config
	must(envconfig.Process("", &cfg), "init config")

	logger.SetLevel(cfg.LogLevel)
	logger.Infof("start with config: %+v", cfg)
	// automatically set GOMAXPROCS to match Linux container CPU quota.
	_, _ = maxprocs.Set(maxprocs.Logger(logger.Infof))

	a, err := process.New()
	must(err, "init runner")

	a.MustAddActor(process.NewTracingRunner())

	publicRouter := initPublicRouter()
	a.MustAddActor(process.NewHTTPServerRunner(
		publicRouter,
		"public",
		cfg.HTTPPort,
		cfg.GracefulShutdownTimeout),
	)
	debugRouter := initDebugRouter(a, publicRouter)
	a.MustAddActor(process.NewHTTPServerRunner(
		debugRouter,
		"debug",
		cfg.DebugHTTPPort,
		cfg.GracefulShutdownTimeout),
	)

	if err := a.Run(); err != nil {
		logger.Error(err)
	}
}

// public http
func initPublicRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(tracing.Middleware)
	// enable request structured logs
	r.Use(logger.NewStructuredLogger(logger.Logger()))
	r.Use(middleware.Recoverer)
	// enable metrics mw
	r.Use(chiprometheus.NewMiddleware(""))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	controller.MountRoutes(r)
	return r
}

// debug http
func initDebugRouter(a *process.Runner, publicRouter chi.Router) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	// healthcheck handlers
	r.Handle("/live", http.HandlerFunc(a.HealthCheck.LiveEndpoint))   // is container alive
	r.Handle("/ready", http.HandlerFunc(a.HealthCheck.ReadyEndpoint)) // is container ready to accept requests

	// prometheus metrics handler
	r.Mount("/metrics", promhttp.Handler())

	// pprof debug handlers
	r.Mount("/debug", middleware.Profiler())

	// swagger ui and swagger.json
	r.Get("/swagger/*", httpswagger.Handler(httpswagger.URL("/swagger/doc.json")))

	// mount public router to serve requests from swagger ui
	r.Mount("/", publicRouter)

	return r
}

func must(err error, message string) {
	if err != nil {
		logger.Fatalf("%s failure: %s", message, err)
	}
}
