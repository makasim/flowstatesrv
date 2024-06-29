package srvdriver

import (
	"context"
	"log"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1alpha1"
	v1alpha1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1alpha1"
	"github.com/makasim/flowstatesrv/protogen/flowstate/v1alpha1/flowstatev1alpha1connect"
)

type Watcher struct {
	ec flowstatev1alpha1connect.EngineServiceClient
}

func newWatcher(ec flowstatev1alpha1connect.EngineServiceClient) *Watcher {
	return &Watcher{
		ec: ec,
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

		ec:       d.ec,
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

	ec       flowstatev1alpha1connect.EngineServiceClient
	watchCh  chan flowstate.State
	closeCh  chan struct{}
	closedCh chan struct{}
}

func (lis *listener) Watch() <-chan flowstate.State {
	return lis.watchCh
}

func (lis *listener) Close() {
	log.Println(1000000)
	close(lis.closeCh)
	log.Println(1000001)

	<-lis.closedCh
}

func (lis *listener) listen() {
	srvS, err := lis.ec.Watch(context.Background(), connect.NewRequest(&v1alpha1.WatchRequest{
		SinceRev:    lis.sinceRev,
		SinceLatest: lis.sinceLatest,
		Labels:      lis.labels,
	}))
	if err != nil {
		log.Println("WARN: call: watch: ", err)
		return
	}
	go func() {
		defer close(lis.closedCh)

		log.Println(1000002)
		<-lis.closeCh
		log.Println(1000003)
		if err := srvS.Close(); err != nil {
			log.Println("WARN: watch stream: close: ", err)
		}
		log.Println(1000004)
	}()

	for srvS.Receive() {
		state := convertorv1alpha1.ConvertAPIToState(srvS.Msg().State)
		select {
		case lis.watchCh <- state:
			continue
		case <-lis.closeCh:
			return
		}
	}
	if srvS.Err() != nil {
		log.Println("WARN: watch: receive: ", srvS.Err())
	}
}
