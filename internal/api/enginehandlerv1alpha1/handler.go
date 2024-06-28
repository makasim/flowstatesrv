package enginehandlerv1alpha1

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	v1alpha1 "github.com/makasim/flowstatesrv/internal/protogen/flowstate/v1alpha1"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

type Handler struct {
	e *flowstate.Engine
}

func New(e *flowstate.Engine) *Handler {
	return &Handler{
		e: e,
	}
}

func (s *Handler) Do(_ context.Context, req *connect.Request[v1alpha1.DoRequest]) (*connect.Response[v1alpha1.DoResponse], error) {
	stateCtxs := make([]*flowstate.StateCtx, 0, len(req.Msg.Contexts))
	for _, apiS := range req.Msg.Contexts {
		stateCtxs = append(stateCtxs, convAPIToStateCtx(apiS))
	}

	cmds := make([]flowstate.Command, 0, len(req.Msg.Commands))
	for _, apiC := range req.Msg.Commands {
		cmd, err := convAPIToCommand(apiC, stateCtxs)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		cmds = append(cmds, cmd)
	}

	if err := s.e.Do(cmds...); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	results := make([]*anypb.Any, 0, len(cmds))
	for _, cmd := range cmds {
		cmdRes, err := convCommandToAPI(cmd)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		anyCmdRes, err := anypb.New(cmdRes)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		results = append(results, anyCmdRes)
	}

	return connect.NewResponse(&v1alpha1.DoResponse{
		Contexts: convStateCtxsToAPI(stateCtxs),
		Results:  results,
	}), nil
}

func (s *Handler) Watch(ctx context.Context, req *connect.Request[v1alpha1.WatchRequest], stream *connect.ServerStream[v1alpha1.WatchResponse]) error {
	wCmd := flowstate.GetWatcher(req.Msg.SinceRev, req.Msg.Labels)
	wCmd.SinceLatest = req.Msg.SinceLatest

	if err := s.e.Do(wCmd); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	w := wCmd.Watcher
	defer w.Close()

	for {
		select {
		case state := <-w.Watch():
			apiS := convStateToAPI(state)
			if err := stream.Send(&v1alpha1.WatchResponse{
				State: apiS,
			}); err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func findStateCtxByRef(ref *v1alpha1.StateRef, stateCtxs []*flowstate.StateCtx) (*flowstate.StateCtx, error) {
	for _, stateCtx := range stateCtxs {
		if string(stateCtx.Current.ID) == ref.Id && stateCtx.Current.Rev == ref.Rev {
			return stateCtx, nil
		}
	}

	return nil, fmt.Errorf("there is no state ctx provided for ref: %s:%d", ref.Id, ref.Rev)
}

func convAPIToCommand(apiC *anypb.Any, stateCtxs []*flowstate.StateCtx) (flowstate.Command, error) {
	switch apiC.TypeUrl {
	case `type.googleapis.com/flowstate.v1alpha1.Transit`:
		apiCmd := &v1alpha1.Transit{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := findStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Transit(stateCtx, flowstate.FlowID(apiCmd.FlowId)), nil
	case `type.googleapis.com/flowstate.v1alpha1.Pause`:
		apiCmd := &v1alpha1.Pause{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := findStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Pause(stateCtx, flowstate.FlowID(apiCmd.FlowId)), nil
	case `type.googleapis.com/flowstate.v1alpha1.Resume`:
		apiCmd := &v1alpha1.Resume{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := findStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Resume(stateCtx), nil
	case `type.googleapis.com/flowstate.v1alpha1.End`:
		apiCmd := &v1alpha1.End{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := findStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.End(stateCtx), nil
	case `type.googleapis.com/flowstate.v1alpha1.Execute`:
		apiCmd := &v1alpha1.Execute{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := findStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Execute(stateCtx), nil
	case `type.googleapis.com/flowstate.v1alpha1.Delay`:
		apiCmd := &v1alpha1.Delay{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := findStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		dur, err := time.ParseDuration(apiCmd.Duration)
		if err != nil {
			return nil, err
		}

		cmd := flowstate.Delay(stateCtx, dur)
		cmd.Commit = apiCmd.Commit

		return cmd, nil
	case `type.googleapis.com/flowstate.v1alpha1.Commit`:
		apiCmd := &v1alpha1.Commit{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		subCmds := make([]flowstate.Command, 0, len(apiCmd.Commands))
		for _, subCmd := range apiCmd.Commands {
			subCmd, err := convAPIToCommand(subCmd, stateCtxs)
			if err != nil {
				return nil, err
			}

			subCmds = append(subCmds, subCmd)
		}

		return flowstate.Commit(subCmds...), nil
	default:
		return nil, fmt.Errorf("unknown command %s", apiC.TypeUrl)
	}
}

func convAPIToStateCtx(apiS *v1alpha1.Context) *flowstate.StateCtx {
	return &flowstate.StateCtx{
		Current:     convAPIToState(apiS.Current),
		Committed:   convAPIToState(apiS.Committed),
		Transitions: convAPIToTransitions(apiS.Transitions),
	}
}

func convAPIToState(apiS *v1alpha1.State) flowstate.State {
	if apiS == nil {
		return flowstate.State{}
	}

	return flowstate.State{
		ID:                   flowstate.StateID(apiS.Id),
		Rev:                  apiS.Rev,
		Annotations:          copyMap(apiS.Annotations),
		Labels:               copyMap(apiS.Labels),
		CommittedAtUnixMilli: apiS.CommittedAtUnixMilli,
		Transition:           convAPIToTransition(apiS.Transition),
	}
}

func convAPIToTransitions(apiTs []*v1alpha1.Transition) []flowstate.Transition {
	ts := make([]flowstate.Transition, 0, len(apiTs))
	for _, apiT := range apiTs {
		ts = append(ts, convAPIToTransition(apiT))
	}
	return ts
}

func convAPIToTransition(apiT *v1alpha1.Transition) flowstate.Transition {
	if apiT == nil {
		return flowstate.Transition{}
	}

	return flowstate.Transition{
		FromID:      flowstate.FlowID(apiT.From),
		ToID:        flowstate.FlowID(apiT.To),
		Annotations: copyMap(apiT.Annotations),
	}
}

func convCommandToAPI(cmd flowstate.Command) (*anypb.Any, error) {
	switch cmd1 := cmd.(type) {
	case *flowstate.TransitCommand:
		return anypb.New(&v1alpha1.TransitResult{
			StateRef: convStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.PauseCommand:
		return anypb.New(&v1alpha1.PauseResult{
			StateRef: convStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.ResumeCommand:
		return anypb.New(&v1alpha1.ResumeResult{
			StateRef: convStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.EndCommand:
		return anypb.New(&v1alpha1.EndResult{
			StateRef: convStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.ExecuteCommand:
		return anypb.New(&v1alpha1.ExecuteResult{
			StateRef: convStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.DelayCommand:
		return anypb.New(&v1alpha1.DelayResult{
			StateRef: convStateCtxToRefAPI(cmd1.StateCtx),
			Duration: cmd1.Duration.String(),
			Commit:   cmd1.Commit,
		})
	case *flowstate.CommitCommand:
		subResults := make([]*anypb.Any, 0, len(cmd1.Commands))
		for _, subCmd := range cmd1.Commands {
			subRes, err := convCommandToAPI(subCmd)
			if err != nil {
				return nil, err
			}

			subResults = append(subResults, subRes)
		}

		return anypb.New(&v1alpha1.CommitResult{
			Results: subResults,
		})
	default:
		return nil, fmt.Errorf("unknown command type %T", cmd)
	}
}

func convStateCtxsToAPI(ss []*flowstate.StateCtx) []*v1alpha1.Context {
	apiS := make([]*v1alpha1.Context, 0, len(ss))
	for _, stateCtx := range ss {
		apiS = append(apiS, convStateCtxToAPI(stateCtx))
	}
	return apiS
}

func convStateCtxToAPI(s *flowstate.StateCtx) *v1alpha1.Context {
	return &v1alpha1.Context{
		Committed:   convStateToAPI(s.Committed),
		Current:     convStateToAPI(s.Current),
		Transitions: convTransitionsToAPI(s.Transitions),
	}
}

func convStateCtxToRefAPI(s *flowstate.StateCtx) *v1alpha1.StateRef {
	return &v1alpha1.StateRef{
		Id:  string(s.Current.ID),
		Rev: s.Current.Rev,
	}
}

func convStateToAPI(s flowstate.State) *v1alpha1.State {
	return &v1alpha1.State{
		Id:                   string(s.ID),
		Rev:                  s.Rev,
		CommittedAtUnixMilli: s.CommittedAtUnixMilli,
		Transition:           convTransitionToAPI(s.Transition),
		Labels:               copyMap(s.Labels),
		Annotations:          copyMap(s.Annotations),
	}
}

func convTransitionsToAPI(tss []flowstate.Transition) []*v1alpha1.Transition {
	apiTss := make([]*v1alpha1.Transition, 0, len(tss))
	for _, ts := range tss {
		apiTss = append(apiTss, convTransitionToAPI(ts))
	}
	return apiTss
}

func convTransitionToAPI(ts flowstate.Transition) *v1alpha1.Transition {
	return &v1alpha1.Transition{
		From:        string(ts.FromID),
		To:          string(ts.ToID),
		Annotations: copyMap(ts.Annotations),
	}
}

func copyMap(from map[string]string) map[string]string {
	if from == nil {
		return nil
	}

	copied := make(map[string]string, len(from))
	for k, v := range from {
		copied[k] = v
	}
	return copied
}
