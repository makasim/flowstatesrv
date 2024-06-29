package srvdriver

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1alpha1"
	v1alpha1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1alpha1"
	"github.com/makasim/flowstatesrv/protogen/flowstate/v1alpha1/flowstatev1alpha1connect"
	anypb "google.golang.org/protobuf/types/known/anypb"
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
	apiCmd, err := convertorv1alpha1.CommandToAPICommand(cmd0)
	if err != nil {
		return err
	}
	apiStateCtxs := convertorv1alpha1.ConvertCommandToStateContexts(cmd0)

	resp, err := d.ec.Do(context.Background(), connect.NewRequest(&v1alpha1.DoRequest{
		StateContexts: apiStateCtxs,
		Commands:      []*anypb.Any{apiCmd},
	}))
	if err != nil {
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

	stateCtxs := convertorv1alpha1.ConvertAPIToStateCtxs(resp.StateContexts)
	for i, cmd0 := range cmds {
		if err := syncCommandWithResult(cmd0, resp.Results[i], stateCtxs); err != nil {
			return err
		}
	}

	return nil
}

func syncCommandWithResult(cmd0 flowstate.Command, res *anypb.Any, stateCtxs []*flowstate.StateCtx) error {
	switch cmd := cmd0.(type) {
	case *flowstate.TransitCommand:
		if res.TypeUrl != `type.googleapis.com/flowstate.v1alpha1.TransitResult` {
			return fmt.Errorf("unexpected result type %s", res.TypeUrl)
		}

		apiRes := &v1alpha1.TransitResult{}
		if err := res.UnmarshalTo(apiRes); err != nil {
			return err
		}

		stateCtx, err := convertorv1alpha1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.PauseCommand:
		if res.TypeUrl != `type.googleapis.com/flowstate.v1alpha1.PauseResult` {
			return fmt.Errorf("unexpected result type %s", res.TypeUrl)
		}

		apiRes := &v1alpha1.PauseResult{}
		if err := res.UnmarshalTo(apiRes); err != nil {
			return err
		}

		stateCtx, err := convertorv1alpha1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.ResumeCommand:
		if res.TypeUrl != `type.googleapis.com/flowstate.v1alpha1.ResumeResult` {
			return fmt.Errorf("unexpected result type %s", res.TypeUrl)
		}

		apiRes := &v1alpha1.ResumeResult{}
		if err := res.UnmarshalTo(apiRes); err != nil {
			return err
		}

		stateCtx, err := convertorv1alpha1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.EndCommand:
		if res.TypeUrl != `type.googleapis.com/flowstate.v1alpha1.EndResult` {
			return fmt.Errorf("unexpected result type %s", res.TypeUrl)
		}

		apiRes := &v1alpha1.EndResult{}
		if err := res.UnmarshalTo(apiRes); err != nil {
			return err
		}

		stateCtx, err := convertorv1alpha1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.ExecuteCommand:
		if res.TypeUrl != `type.googleapis.com/flowstate.v1alpha1.ExecuteResult` {
			return fmt.Errorf("unexpected result type %s", res.TypeUrl)
		}

		apiRes := &v1alpha1.ExecuteResult{}
		if err := res.UnmarshalTo(apiRes); err != nil {
			return err
		}

		stateCtx, err := convertorv1alpha1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.DelayCommand:
		if res.TypeUrl != `type.googleapis.com/flowstate.v1alpha1.DelayResult` {
			return fmt.Errorf("unexpected result type %s", res.TypeUrl)
		}

		apiRes := &v1alpha1.DelayResult{}
		if err := res.UnmarshalTo(apiRes); err != nil {
			return err
		}

		stateCtx, err := convertorv1alpha1.FindStateCtxByRef(apiRes.StateRef, stateCtxs)
		if err != nil {
			return err
		}

		stateCtx.CopyTo(cmd.StateCtx)
		return nil
	case *flowstate.CommitCommand:
		if res.TypeUrl != `type.googleapis.com/flowstate.v1alpha1.CommitResult` {
			return fmt.Errorf("unexpected result type %s", res.TypeUrl)
		}

		apiRes := &v1alpha1.CommitResult{}
		if err := res.UnmarshalTo(apiRes); err != nil {
			return err
		}

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
