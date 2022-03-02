package process

import (
	"sync"

	"github.com/heptiolabs/healthcheck"
	"github.com/kannman/waterbucket/internal/pkg/logger"
	"github.com/pkg/errors"

	"github.com/oklog/run"
)

type Runner struct {
	runGroup    run.Group
	HealthCheck healthcheck.Handler
	closer      *closerRunner

	mu    sync.RWMutex // protects flag below
	ready bool
}

func New() (*Runner, error) {
	closerRunner := &closerRunner{}
	closerRunner.add(func() error {
		// sync logger before app is terminated.
		return logger.Logger().Sync()
	})

	app := &Runner{
		HealthCheck: healthcheck.NewHandler(),
		closer:      closerRunner,
	}
	app.HealthCheck.AddReadinessCheck("ready", app.checkReadiness)

	// init actors
	app.AddActor(interruptActor(app.setReadiness))
	app.AddActor(closerRunner.actor())

	return app, nil
}

// AddActor (function) to the application group. Each actor must be pre-emptable by an
// interrupt function. That is, if interrupt is invoked, execute should return.
// Also, it must be safe to call interrupt even after execute has returned.
//
// The first actor (function) to return interrupts all running actors.
// The error is passed to the interrupt functions, and is returned by Run.
func (r *Runner) AddActor(execute func() error, interrupt func(error)) {
	r.runGroup.Add(execute, interrupt)
}

func (r *Runner) MustAddActor(execute func() error, interrupt func(error), err error) {
	if err != nil {
		logger.Fatalf("add actor failure: %s", err)
	}
	r.runGroup.Add(execute, interrupt)
}

func (r *Runner) AddCloser(closer func() error) {
	r.closer.add(closer)
}

func (r *Runner) Run() error {
	logger.Warn("application started")
	defer logger.Warn("application stopped")
	r.setReadiness(true)
	return r.runGroup.Run()
}

func (r *Runner) setReadiness(state bool) {
	r.mu.Lock()
	r.ready = state
	r.mu.Unlock()
}

func (r *Runner) checkReadiness() error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if !r.ready {
		return errors.New("application is not ready")
	}
	return nil
}
