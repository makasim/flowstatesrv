package srvdriver

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1"
	v1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1"
	"github.com/makasim/flowstatesrv/protogen/flowstate/v1/flowstatev1connect"
)

type RemoteDoer struct {
	e  flowstate.Engine
	sc flowstatev1connect.ServerServiceClient
}

func newRemoteDoer(sc flowstatev1connect.ServerServiceClient) *RemoteDoer {
	return &RemoteDoer{
		sc: sc,
	}
}

func (d *RemoteDoer) Init(e flowstate.Engine) error {
	d.e = e
	return nil
}

func (d *RemoteDoer) Shutdown(_ context.Context) error {
	return nil
}

func (d *RemoteDoer) Do(cmd0 flowstate.Command) error {
	switch cmd0.(type) {
	case *flowstate.CommitCommand:
		return d.do(cmd0)
	case *flowstate.DelayCommand:
		return d.do(cmd0)
	case *flowstate.StoreDataCommand:
		return d.do(cmd0)
	case *flowstate.GetDataCommand:
		return d.do(cmd0)
	case *flowstate.GetCommand:
		return d.do(cmd0)
	case *flowstate.GetManyCommand:
		return d.do(cmd0)
	case *flowstate.CommitStateCtxCommand:
		return d.do(cmd0)
	default:
		return flowstate.ErrCommandNotSupported
	}
}

func (d *RemoteDoer) do(cmd0 flowstate.Command) error {
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
	if conflictErr := asRevMismatchError(err); conflictErr != nil {
		return conflictErr
	} else if err != nil {
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
	case *flowstate.GetCommand:
		if res.GetGet() == nil {
			return fmt.Errorf("got get nil result")
		}

		apiRes := res.GetGet()

		stateCtx, err := convertorv1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.GetManyCommand:
		if res.GetGetMany() == nil {
			return fmt.Errorf("got get many nil result")
		}

		apiRes := res.GetGetMany()

		cmd.SetResult(&flowstate.GetManyResult{
			States: convertorv1.ConvertAPIToStates(apiRes.States),
			More:   apiRes.More,
		})

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

func asRevMismatchError(err error) *flowstate.ErrRevMismatch {
	// See https://connectrpc.com/docs/go/errors/#error-details

	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		return nil
	}

	for _, detail := range connectErr.Details() {
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
