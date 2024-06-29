package srvdriver

import (
	"context"
	"fmt"

	"github.com/makasim/flowstate"
)

type Watcher struct {
}

func newWatcher() *Watcher {
	return &Watcher{}
}

func (d *Watcher) Init(_ *flowstate.Engine) error {
	return nil
}

func (d *Watcher) Shutdown(_ context.Context) error {
	return nil
}

func (d *Watcher) Do(cmd0 flowstate.Command) error {
	_, ok := cmd0.(*flowstate.GetWatcherCommand)
	if !ok {
		return flowstate.ErrCommandNotSupported
	}

	return fmt.Errorf("wathcer not implemented")
}
