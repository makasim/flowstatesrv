package srvdriver

import (
	"context"
	"errors"
	"log"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1"
	v1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1"
	"github.com/makasim/flowstatesrv/protogen/flowstate/v1/flowstatev1connect"
)

type Watcher struct {
	sc flowstatev1connect.ServerServiceClient
}

func newWatcher(sc flowstatev1connect.ServerServiceClient) *Watcher {
	return &Watcher{
		sc: sc,
	}
}

func (d *Watcher) Init(_ *flowstate.Engine) error {
	return nil
}

func (d *Watcher) Shutdown(_ context.Context) error {
	return nil
}

func (d *Watcher) Do(cmd0 flowstate.Command) error {
	cmd, ok := cmd0.(*flowstate.GetWatcherCommand)
	if !ok {
		return flowstate.ErrCommandNotSupported
	}

	lis := &listener{
		sinceRev:    cmd.SinceRev,
		labels:      make(map[string]string),
		sinceLatest: cmd.SinceLatest,

		ec:       d.sc,
		watchCh:  make(chan flowstate.State, 1),
		closeCh:  make(chan struct{}),
		closedCh: make(chan struct{}),
	}
	for k, v := range cmd.Labels {
		lis.labels[k] = v
	}

	go lis.listen()

	cmd.Watcher = lis
	return nil
}

type listener struct {
	sinceRev    int64
	sinceLatest bool
	labels      map[string]string

	ec       flowstatev1connect.ServerServiceClient
	watchCh  chan flowstate.State
	closeCh  chan struct{}
	closedCh chan struct{}
}

func (lis *listener) Watch() <-chan flowstate.State {
	return lis.watchCh
}

func (lis *listener) Close() {
	close(lis.closeCh)
	<-lis.closedCh
}

func (lis *listener) listen() {
	defer close(lis.closedCh)

	wCtx, wCtxCancel := context.WithCancel(context.Background())

	srvS, err := lis.ec.WatchStates(wCtx, connect.NewRequest(&v1.WatchStatesRequest{
		SinceRev:    lis.sinceRev,
		SinceLatest: lis.sinceLatest,
		Labels:      lis.labels,
	}))
	if err != nil {
		wCtxCancel()
		log.Println("WARN: call: watch: ", err)
		return
	}
	defer srvS.Close()

	go func() {
		<-lis.closeCh
		wCtxCancel()
	}()

	for srvS.Receive() {
		state := convertorv1.ConvertAPIToState(srvS.Msg().State)
		select {
		case lis.watchCh <- state:
			continue
		case <-lis.closeCh:
			return
		}
	}
	if srvS.Err() != nil && !errors.Is(srvS.Err(), context.Canceled) {
		log.Println("WARN: watch: receive: ", srvS.Err())
	}
}
