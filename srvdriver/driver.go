package srvdriver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/bufbuild/httplb"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstate/exptcmd"
	"github.com/makasim/flowstate/stddoer"
	"github.com/makasim/flowstatesrv/protogen/flowstate/flow/v1alpha1/flowv1alpha1connect"
	v1alpha1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1alpha1"
	"github.com/makasim/flowstatesrv/protogen/flowstate/v1alpha1/flowstatev1alpha1connect"
)

type Driver struct {
	*FlowRegistry
	hs    *http.Server
	hc    *httplb.Client
	doers []flowstate.Doer
	ec    flowstatev1alpha1connect.EngineServiceClient
	fc    flowstatev1alpha1connect.FlowServiceClient
}

func New(serverHttpHost string) *Driver {
	hc := httplb.NewClient()

	d := &Driver{
		FlowRegistry: &FlowRegistry{},

		hc: hc,
		ec: flowstatev1alpha1connect.NewEngineServiceClient(hc, serverHttpHost, connect.WithProtoJSON()),
		fc: flowstatev1alpha1connect.NewFlowServiceClient(hc, serverHttpHost, connect.WithProtoJSON()),
	}

	doers := []flowstate.Doer{
		stddoer.Transit(),
		stddoer.Pause(),
		stddoer.Resume(),
		stddoer.End(),
		stddoer.Noop(),

		exptcmd.NewStacker(),
		exptcmd.UnstackDoer(),

		newFlowGetter(d.FlowRegistry),
		newWatcher(d.ec),
		newRemoteDoer(d.ec),
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
	mux := http.NewServeMux()
	mux.Handle(flowv1alpha1connect.NewFlowServiceHandler(newHandler(e)))

	d.hs = &http.Server{
		Addr:    `:23654`,
		Handler: mux,
	}

	go func() {
		if err := d.hs.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("WARN: http server: listen and serve: %s", err)
		}
	}()

	for fID := range d.FlowRegistry.flows {
		if _, err := d.fc.Register(context.Background(), connect.NewRequest(&v1alpha1.RegisterRequest{
			FlowId:   string(fID),
			HttpHost: `http://127.0.0.1:23654`,
		})); err != nil {
			return fmt.Errorf("register flow: %w", err)
		}
	}

	return nil
}

func (d *Driver) Shutdown(_ context.Context) error {
	if err := d.hs.Shutdown(context.Background()); err != nil {
		return err
	}

	if err := d.hc.Close(); err != nil {
		return err
	}

	return nil
}
