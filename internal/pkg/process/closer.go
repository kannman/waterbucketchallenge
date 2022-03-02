package process

import (
	"sync"

	"github.com/kannman/waterbucketchallenge/internal/pkg/logger"
)

type closerRunner struct {
	lock    sync.Mutex
	closers []func() error
}

func (a *closerRunner) add(c ...func() error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.closers = append(c, a.closers...)
}

func (a *closerRunner) actor() (func() error, func(error)) {
	var runWg sync.WaitGroup
	runWg.Add(1)
	return func() error {
			runWg.Wait()
			return nil
		}, func(err error) {
			runWg.Done()

			a.lock.Lock()
			defer a.lock.Unlock()
			var wg sync.WaitGroup
			wg.Add(len(a.closers))
			for _, c := range a.closers {
				closer := c
				go func() {
					defer wg.Done()
					if err := closer(); err != nil {
						logger.Errorf("closer: %s", err)
					}
				}()
			}
			wg.Wait()
		}
}
