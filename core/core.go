package core

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/tiny-sky/Tvote/log"
	"golang.org/x/sync/errgroup"
)

type Option func(core *Core)

type Core struct {
	server       []Server
	stopCtx      context.Context
	runWaitGroup sync.WaitGroup
	errGroup     *errgroup.Group
	cancel       func()
	once         sync.Once
}

func WithServers(srvs ...Server) Option {
	return func(core *Core) {
		core.server = append(core.server, srvs...)
	}
}

func New(opts ...Option) *Core {
	core := &Core{
		runWaitGroup: sync.WaitGroup{},
		once:         sync.Once{},
	}
	for _, opt := range opts {
		opt(core)
	}
	return core
}

func (core *Core) Run(ctx context.Context) error {
	var c1 context.Context

	c1, core.cancel = context.WithCancel(ctx)
	core.errGroup, core.stopCtx = errgroup.WithContext(c1)

	for _, server := range core.server {
		core.runWaitGroup.Add(1)
		srv := server
		core.errGroup.Go(func() error {
			<-core.stopCtx.Done()
			return srv.Stop(ctx)
		})

		core.errGroup.Go(func() error {
			defer core.runWaitGroup.Done()
			return srv.Run(ctx)
		})
	}

	core.runWaitGroup.Wait()
	log.Infof("start")

	// 优雅关停
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	core.errGroup.Go(func() error {
		select {
		case <-core.stopCtx.Done():
			return core.stopCtx.Err()
		case <-c:
			return core.Stop()
		}
	})
	if err := core.errGroup.Wait(); err != nil {
		return err
	}
	return nil
}

func (core *Core) Stop() (err error) {
	if core.cancel == nil {
		return nil
	}
	core.once.Do(func() {
		core.cancel()
	})
	return
}
