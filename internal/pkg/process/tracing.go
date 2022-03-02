package process

import (
	"github.com/kannman/waterbucket/internal/pkg/logger"
	"github.com/kannman/waterbucket/internal/pkg/tracing"
	"github.com/pkg/errors"
)

func NewTracingRunner() (func() error, func(error), error) {
	closer, err := tracing.InitTracer()
	if err != nil {
		return nil, nil, errors.Wrap(err, "init tracing")
	}
	doneCh := make(chan struct{})
	return func() error {
			<-doneCh
			return nil
		}, func(err error) {
			close(doneCh)
			if err := closer.Close(); err != nil {
				logger.Errorf("close tracing: %s", err)
			}
			logger.Warn("tracing closed")
		}, nil
}
