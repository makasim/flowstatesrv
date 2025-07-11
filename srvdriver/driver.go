package srvdriver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1"
	flowstatev1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1"
	v1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1"
	"github.com/makasim/flowstatesrv/protogen/flowstate/v1/flowstatev1connect"
)

type Driver struct {
	*flowstate.FlowRegistry

	hs *http.Server
	hc *http.Client
	sc flowstatev1connect.ServerServiceClient
	l  *slog.Logger
}

func New(serverHttpHost string, l *slog.Logger) (*Driver, error) {
	hc := &http.Client{}

	d := &Driver{
		FlowRegistry: &flowstate.FlowRegistry{},

		hc: hc,
		sc: flowstatev1connect.NewServerServiceClient(hc, serverHttpHost, connect.WithProtoJSON()),

		l: l,
	}

	return d, nil
}

func (d *Driver) Init(e flowstate.Engine) error {
	mux := http.NewServeMux()
	mux.Handle(flowstatev1connect.NewFlowServiceHandler(newHandler(e)))

	d.hs = &http.Server{
		Addr:    `:23654`,
		Handler: mux,
	}

	go func() {
		if err := d.hs.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			d.l.Warn(fmt.Sprintf("http server: listen and serve: %s", err))
		}
	}()

	return nil
}

func (d *Driver) GetStateByID(cmd *flowstate.GetStateByIDCommand) error {
	return d.do(cmd)
}

func (d *Driver) GetStateByLabels(cmd *flowstate.GetStateByLabelsCommand) error {
	return d.do(cmd)
}

func (d *Driver) GetStates(cmd *flowstate.GetStatesCommand) (*flowstate.GetStatesResult, error) {
	if err := d.do(cmd); err != nil {
		return nil, fmt.Errorf("get states: %w", err)
	}
	return cmd.Result, nil
}

func (d *Driver) GetDelayedStates(cmd *flowstate.GetDelayedStatesCommand) (*flowstate.GetDelayedStatesResult, error) {
	if err := d.do(cmd); err != nil {
		return nil, fmt.Errorf("get states: %w", err)
	}
	return cmd.Result, nil
}

func (d *Driver) GetData(cmd *flowstate.GetDataCommand) error {
	return d.do(cmd)
}

func (d *Driver) Delay(cmd *flowstate.DelayCommand) error {
	return d.do(cmd)
}

func (d *Driver) StoreData(cmd *flowstate.StoreDataCommand) error {
	return d.do(cmd)
}

func (d *Driver) Commit(cmd *flowstate.CommitCommand) error {
	return d.do(cmd)
}

func (d *Driver) Flow(id flowstate.FlowID) (flowstate.Flow, error) {
	flow, err := d.FlowRegistry.Flow(id)
	if err == nil {
		return flow, nil
	}

	// TODO look for flow on server
	return nil, err
}

func (d *Driver) SetFlow(id flowstate.FlowID, flow flowstate.Flow) error {
	if _, err := d.sc.RegisterFlow(context.Background(), connect.NewRequest(&flowstatev1.RegisterFlowRequest{
		FlowId:   string(id),
		HttpHost: `http://127.0.0.1:23654`,
	})); err != nil {
		return fmt.Errorf("register flow: %w", err)
	}

	if err := d.FlowRegistry.SetFlow(id, flow); err != nil {
		return fmt.Errorf("set flow: %w", err)
	}

	return nil
}

func (d *Driver) Shutdown(_ context.Context) error {
	if d.hs != nil {
		if err := d.hs.Shutdown(context.Background()); err != nil {
			return err
		}
	}

	d.hc.CloseIdleConnections()

	return nil
}

func (d *Driver) do(cmd0 flowstate.Command) error {
	apiCmd, err := convertorv1.CommandToAPICommand(cmd0)
	if err != nil {
		return err
	}
	apiStateCtxs := convertorv1.ConvertCommandToStateContexts(cmd0)
	apiDatas := convertorv1.ConvertCommandToDatas(cmd0)

	resp, err := d.sc.DoCommand(context.Background(), connect.NewRequest(&v1.DoCommandRequest{
		StateContexts: apiStateCtxs,
		Data:          apiDatas,
		Commands:      []*v1.Command{apiCmd},
	}))
	if err != nil {
		var connectErr *connect.Error
		if !errors.As(err, &connectErr) {
			return err
		}

		if revMismatchErr := asRevMismatchError(connectErr); revMismatchErr != nil {
			return revMismatchErr
		} else if connectErr.Code() == connect.CodeNotFound {
			return flowstate.ErrNotFound
		}

		return err
	}

	if err := syncCommandWithDoResponse([]flowstate.Command{cmd0}, resp.Msg); err != nil {
		return err
	}

	return nil
}

func syncCommandWithDoResponse(cmds []flowstate.Command, resp *v1.DoCommandResponse) error {
	if len(cmds) != len(resp.Results) {
		return fmt.Errorf("commands and results count mismatch")
	}

	stateCtxs := convertorv1.ConvertAPIToStateCtxs(resp.StateContexts)
	datas := convertorv1.ConvertAPIToDatas(resp.Data)

	for i, cmd0 := range cmds {
		if err := syncCommandWithResult(cmd0, resp.Results[i], stateCtxs, datas); err != nil {
			return err
		}
	}

	return nil
}

func syncCommandWithResult(cmd0 flowstate.Command, res *v1.Result, stateCtxs []*flowstate.StateCtx, datas []*flowstate.Data) error {
	switch cmd := cmd0.(type) {
	case *flowstate.TransitCommand:
		if res.GetTransit() == nil {
			return fmt.Errorf("got transit nil result")
		}

		apiRes := res.GetTransit()
		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.PauseCommand:
		if res.GetPause() == nil {
			return fmt.Errorf("got pause nil result")
		}

		apiRes := res.GetPause()

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.ResumeCommand:
		if res.GetResume() == nil {
			return fmt.Errorf("got resume nil result")
		}

		apiRes := res.GetResume()

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.EndCommand:
		if res.GetEnd() == nil {
			return fmt.Errorf("got end nil result")
		}

		apiRes := res.GetEnd()

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.ExecuteCommand:
		if res.GetExecute() == nil {
			return fmt.Errorf("got execute nil result")
		}

		apiRes := res.GetExecute()

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.DelayCommand:
		if res.GetDelay() == nil {
			return fmt.Errorf("got delay nil result")
		}

		apiRes := res.GetDelay()

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.SerializeCommand:
		if res.GetSerialize() == nil {
			return fmt.Errorf("got serialize nil result")
		}

		apiRes := res.GetSerialize()

		serializableStateCtx, err := convertorv1.FindStateCtxByRef(apiRes.SerializableStateRef, stateCtxs)
		if err != nil {
			return err
		}
		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		serializableStateCtx.CopyTo(cmd.SerializableStateCtx)
		stateCtx.CopyTo(cmd.StateCtx)
		cmd.Annotation = apiRes.Annotation
		return nil
	case *flowstate.DeserializeCommand:
		if res.GetDeserialize() == nil {
			return fmt.Errorf("got deserialize nil result")
		}

		apiRes := res.GetDeserialize()

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}
		deserializedStateCtx, err := convertorv1.FindStateCtxByRef(apiRes.DeserializedStateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		deserializedStateCtx.CopyTo(cmd.DeserializedStateCtx)
		cmd.Annotation = apiRes.Annotation
		return nil
	case *flowstate.CommitCommand:
		if res.GetCommit() == nil {
			return fmt.Errorf("got commit nil result")
		}

		apiRes := res.GetCommit()

		if len(apiRes.Results) != len(cmd.Commands) {
			return fmt.Errorf("commands and results count mismatch")
		}

		for i, subCmd := range cmd.Commands {
			if err := syncCommandWithResult(subCmd, apiRes.Results[i], stateCtxs, datas); err != nil {
				return err
			}
		}
		return nil
	case *flowstate.StoreDataCommand:
		if res.GetStoreData() == nil {
			return fmt.Errorf("got store data nil result")
		}
		apiRes := res.GetStoreData()
		d, err := convertorv1.FindDataByRef(apiRes.DataRef, datas)
		if err != nil {
			return err
		}
		d.CopyTo(cmd.Data)
		return nil
	case *flowstate.GetDataCommand:
		if res.GetGetData() == nil {
			return fmt.Errorf("got get data nil result")
		}

		apiRes := res.GetGetData()
		d, err := convertorv1.FindDataByRef(apiRes.DataRef, datas)
		if err != nil {
			return err
		}
		d.CopyTo(cmd.Data)

		return nil
	case *flowstate.ReferenceDataCommand:
		if res.GetReferenceData() == nil {
			return fmt.Errorf("get reference data nil result")
		}

		apiRes := res.GetReferenceData()

		d, err := convertorv1.FindDataByRef(apiRes.DataRef, datas)
		if err != nil {
			return err
		}
		d.CopyTo(cmd.Data)

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}
		stateCtx.CopyTo(cmd.StateCtx)

		return nil
	case *flowstate.DereferenceDataCommand:
		if res.GetDereferenceData() == nil {
			return fmt.Errorf("got dereference data nil result")
		}

		apiRes := res.GetDereferenceData()

		d, err := convertorv1.FindDataByRef(apiRes.DataRef, datas)
		if err != nil {
			return err
		}
		d.CopyTo(cmd.Data)

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.GetStateByIDCommand:
		if res.GetGetStateById() == nil {
			return fmt.Errorf("got get state by id nil result")
		}

		apiRes := res.GetGetStateById()

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.GetStateByLabelsCommand:
		if res.GetGetStateByLabels() == nil {
			return fmt.Errorf("got get state by labels nil result")
		}

		apiRes := res.GetGetStateByLabels()

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.GetStatesCommand:
		if res.GetGetStates() == nil {
			return fmt.Errorf("got get states nil result")
		}

		apiRes := res.GetGetStates()

		cmd.Result = &flowstate.GetStatesResult{
			States: convertorv1.ConvertAPIToStates(apiRes.States),
			More:   apiRes.More,
		}

		return nil
	case *flowstate.GetDelayedStatesCommand:
		if res.GetGetDelayedStates() == nil {
			return fmt.Errorf("got get delayed states nil result")
		}

		apiRes := res.GetGetDelayedStates()

		cmd.Result = &flowstate.GetDelayedStatesResult{
			States: convertorv1.ConvertAPIToDelayedStates(apiRes.DelayedStates),
			More:   apiRes.More,
		}

		return nil
	case *flowstate.CommitStateCtxCommand:
		if res.GetCommitState() == nil {
			return fmt.Errorf("got commit state nil result")
		}

		apiRes := res.GetCommitState()

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	default:
		return fmt.Errorf("unknown command %T", cmd0)
	}
}

func asRevMismatchError(err *connect.Error) *flowstate.ErrRevMismatch {
	// See https://connectrpc.com/docs/go/errors/#error-details

	for _, detail := range err.Details() {
		msg, valueErr := detail.Value()
		if valueErr != nil {
			continue // usually, errors here mean that we don't have the schema for this Protobuf message
		}

		if apiRevMismatchErr, ok := msg.(*v1.ErrorRevMismatch); ok {
			revMismatchErr := &flowstate.ErrRevMismatch{}
			for _, stateID := range apiRevMismatchErr.CommittableStateIds {
				revMismatchErr.Add("", flowstate.StateID(stateID), nil)
			}

			return revMismatchErr
		}
	}

	return nil
}
