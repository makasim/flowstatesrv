package srvdriver

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1"
	commandv1 "github.com/makasim/flowstatesrv/protogen/flowstate/command/v1"
	v1alpha1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1alpha1"
	"github.com/makasim/flowstatesrv/protogen/flowstate/v1alpha1/flowstatev1alpha1connect"
)

type RemoteDoer struct {
	e  *flowstate.Engine
	ec flowstatev1alpha1connect.EngineServiceClient
}

func newRemoteDoer(ec flowstatev1alpha1connect.EngineServiceClient) *RemoteDoer {
	return &RemoteDoer{
		ec: ec,
	}
}

func (d *RemoteDoer) Init(e *flowstate.Engine) error {
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
	case *flowstate.GetWatcherCommand:
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

	resp, err := d.ec.Do(context.Background(), connect.NewRequest(&v1alpha1.DoRequest{
		StateContexts: apiStateCtxs,
		Commands:      []*commandv1.AnyCommand{apiCmd},
	}))
	if conflictErr := asCommitConflictError(err); conflictErr != nil {
		return conflictErr
	} else if err != nil {
		return err
	}

	if err := syncCommandWithDoResponse([]flowstate.Command{cmd0}, resp.Msg); err != nil {
		return err
	}

	return nil
}

func syncCommandWithDoResponse(cmds []flowstate.Command, resp *v1alpha1.DoResponse) error {
	if len(cmds) != len(resp.Results) {
		return fmt.Errorf("commands and results count mismatch")
	}

	stateCtxs := convertorv1.ConvertAPIToStateCtxs(resp.StateContexts)
	for i, cmd0 := range cmds {
		if err := syncCommandWithResult(cmd0, resp.Results[i], stateCtxs); err != nil {
			return err
		}
	}

	return nil
}

func syncCommandWithResult(cmd0 flowstate.Command, res *commandv1.AnyResult, stateCtxs []*flowstate.StateCtx) error {
	switch cmd := cmd0.(type) {
	case *flowstate.TransitCommand:
		if res.GetTransit() == nil {
			return fmt.Errorf("unexpected result type %T", res.Result)
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
			return fmt.Errorf("unexpected result type %T", res.Result)
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
			return fmt.Errorf("unexpected result type %T", res.Result)
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
			return fmt.Errorf("unexpected result type %T", res.Result)
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
			return fmt.Errorf("unexpected result type %T", res.Result)
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
			return fmt.Errorf("unexpected result type %T", res.Result)
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
			return fmt.Errorf("unexpected result type %T", res.Result)
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
			return fmt.Errorf("unexpected result type %T", res.Result)
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
			return fmt.Errorf("unexpected result type %T", res.Result)
		}

		apiRes := res.GetCommit()

		if len(apiRes.Results) != len(cmd.Commands) {
			return fmt.Errorf("commands and results count mismatch")
		}

		for i, subCmd := range cmd.Commands {
			if err := syncCommandWithResult(subCmd, apiRes.Results[i], stateCtxs); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown command %T", cmd0)
	}
}

func asCommitConflictError(err error) *flowstate.ErrCommitConflict {
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

		if apiErrConflict, ok := msg.(*v1alpha1.ErrorConflict); ok {
			conflictErr := &flowstate.ErrCommitConflict{}
			for _, stateID := range apiErrConflict.CommittableStateIds {
				conflictErr.Add("", flowstate.StateID(stateID), nil)
			}

			return conflictErr
		}
	}

	return nil
}
