package process

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kannman/waterbucketchallenge/internal/pkg/logger"
)

// performs interrupt and graceful shutdown for all actors.
func interruptActor(setReady func(state bool), sig ...os.Signal) (func() error, func(error)) {
	sig = append(sig, syscall.SIGTERM, syscall.SIGINT)
	sigCh := make(chan os.Signal, 1)
	doneCh := make(chan struct{})
	return func() error {
			signal.Notify(sigCh, sig...)
			select {
			case <-doneCh:
			case <-sigCh:
				logger.Warnf("graceful shutdown")
			}
			return nil
		}, func(err error) {
			setReady(false)
			signal.Stop(sigCh)
			close(doneCh)
		}
}
