package driver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/makasim/flowstate"
	"github.com/makasim/flowstate/memdriver"
)

type Driver struct {
	*memdriver.FlowRegistry

	l     *memdriver.Log
	doers []flowstate.Doer
}

func New() *Driver {
	l := &memdriver.Log{}

	d := &Driver{
		l:            l,
		FlowRegistry: &memdriver.FlowRegistry{},
	}

	doers := []flowstate.Doer{
		TransitDoer,
		PauseDoer,
		ResumeDoer,
		EndDoer,
		flowstate.DefaultNoopDoer,
		flowstate.DefaultSerializerDoer,
		flowstate.DefaultDeserializeDoer,
		flowstate.DefaultDereferenceDataDoer,
		flowstate.DefaultReferenceDataDoer,

		flowstate.Recoverer(time.Millisecond * 500),

		memdriver.NewDataLog(),
		memdriver.NewFlowGetter(d.FlowRegistry),
		memdriver.NewCommiter(l),
		memdriver.NewGetter(l),
		memdriver.NewWatcher(l),
		memdriver.NewDelayer(),
	}
	d.doers = doers

	return d
}

func (d *Driver) Do(cmd0 flowstate.Command) error {
	for _, doer := range d.doers {
		if err := doer.Do(cmd0); errors.Is(err, flowstate.ErrCommandNotSupported) {
			continue
		} else if err != nil {
			return fmt.Errorf("%T: do: %w", doer, err)
		}

		return nil
	}

	return fmt.Errorf("no doer for command %T", cmd0)
}

func (d *Driver) Init(e *flowstate.Engine) error {
	for _, doer := range d.doers {
		if err := doer.Init(e); err != nil {
			return fmt.Errorf("%T: init: %w", doer, err)
		}
	}
	return nil
}

func (d *Driver) Shutdown(_ context.Context) error {
	var res error
	for _, doer := range d.doers {
		if err := doer.Shutdown(context.Background()); err != nil {
			res = errors.Join(res, fmt.Errorf("%T: shutdown: %w", doer, err))
		}
	}

	return res
}
