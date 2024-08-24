package srvdriver

import (
	"context"
	"errors"
	"log"
	"time"

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
	cmd, ok := cmd0.(*flowstate.WatchCommand)
	if !ok {
		return flowstate.ErrCommandNotSupported
	}

	lis := &listener{
		sinceRev:    cmd.SinceRev,
		labels:      make([]map[string]string, 0),
		sinceLatest: cmd.SinceLatest,
		sinceTime:   cmd.SinceTime,

		ec:       d.sc,
		watchCh:  make(chan flowstate.State, 1),
		closeCh:  make(chan struct{}),
		closedCh: make(chan struct{}),
	}
	for i := range cmd.Labels {
		lis.labels = append(lis.labels, make(map[string]string))
		for k, v := range cmd.Labels[i] {
			lis.labels[i][k] = v
		}
	}

	go lis.listen()

	cmd.Listener = lis
	return nil
}

type listener struct {
	sinceRev    int64
	sinceLatest bool
	sinceTime   time.Time
	labels      []map[string]string

	ec       flowstatev1connect.ServerServiceClient
	watchCh  chan flowstate.State
	closeCh  chan struct{}
	closedCh chan struct{}
}

func (lis *listener) Listen() <-chan flowstate.State {
	return lis.watchCh
}

func (lis *listener) Close() {
	close(lis.closeCh)
	<-lis.closedCh
}

func (lis *listener) listen() {
	defer close(lis.closedCh)

	wCtx, wCtxCancel := context.WithCancel(context.Background())

	var sinceTimeMSec int64
	if !lis.sinceTime.IsZero() {
		sinceTimeMSec = lis.sinceTime.UnixMilli()
	}

	apiLabels := make([]*v1.WatchStatesRequest_Labels, 0, len(lis.labels))
	for i := range lis.labels {
		apiLabels = append(apiLabels, &v1.WatchStatesRequest_Labels{
			Labels: make(map[string]string),
		})
		for k, v := range lis.labels[i] {
			apiLabels[i].Labels[k] = v
		}
	}

	log.Println(lis.sinceTime.Unix())
	srvS, err := lis.ec.WatchStates(wCtx, connect.NewRequest(&v1.WatchStatesRequest{
		SinceRev:      lis.sinceRev,
		SinceLatest:   lis.sinceLatest,
		SinceTimeMsec: sinceTimeMSec,
		Labels:        apiLabels,
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
