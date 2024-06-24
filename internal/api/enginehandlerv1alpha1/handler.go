package enginehandlerv1alpha1

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	v1alpha1 "github.com/makasim/flowstatesrv/internal/protogen/flowstate/v1alpha1"
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
	stateCtxs := make([]*flowstate.StateCtx, 0, len(req.Msg.StateCtxs))
	for _, apiS := range req.Msg.StateCtxs {
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

	results := make([]*v1alpha1.CommandResult, 0, len(cmds))
	for _, cmd := range cmds {
		cmdRes, err := convCommandToAPI(cmd)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		results = append(results, cmdRes)
	}

	return connect.NewResponse(&v1alpha1.DoResponse{
		StateCtxs: convStateCtxsToAPI(stateCtxs),
		Results:   results,
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

func findStateCtxByRef(ref *v1alpha1.StateCtxRef, stateCtxs []*flowstate.StateCtx) (*flowstate.StateCtx, error) {
	for _, stateCtx := range stateCtxs {
		if string(stateCtx.Current.ID) == ref.Id && stateCtx.Current.Rev == ref.Rev {
			return stateCtx, nil
		}
	}

	return nil, fmt.Errorf("there is no state ctx provided for ref: %s:%d", ref.Id, ref.Rev)
}

func convAPIToCommand(apiC *v1alpha1.Command, stateCtxs []*flowstate.StateCtx) (flowstate.Command, error) {
	switch {
	case apiC.GetTransit() != nil:
		stateCtx, err := findStateCtxByRef(apiC.GetTransit().StateCtx, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Transit(stateCtx, flowstate.FlowID(apiC.GetTransit().FlowId)), nil
	case apiC.GetCommit() != nil:
		subCmds := make([]flowstate.Command, 0, len(apiC.GetCommit().Commands))
		for _, subCmd := range apiC.GetCommit().Commands {
			cmd, err := convAPIToCommand(subCmd, stateCtxs)
			if err != nil {
				return nil, err
			}

			subCmds = append(subCmds, cmd)
		}

		return flowstate.Commit(subCmds...), nil
	default:
		return nil, fmt.Errorf("unknown command type %T", apiC.GetCommand())
	}
}

func convAPIToStateCtx(apiS *v1alpha1.StateCtx) *flowstate.StateCtx {
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
	return flowstate.Transition{
		FromID:      flowstate.FlowID(apiT.From),
		ToID:        flowstate.FlowID(apiT.To),
		Annotations: copyMap(apiT.Annotations),
	}
}

func convCommandToAPI(cmd flowstate.Command) (*v1alpha1.CommandResult, error) {
	switch cmd1 := cmd.(type) {
	case *flowstate.TransitCommand:
		return &v1alpha1.CommandResult{
			Result: &v1alpha1.CommandResult_Transit{
				Transit: &v1alpha1.TransitResult{
					StateCtx: convStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.CommitCommand:
		apiCmdResults := make([]*v1alpha1.CommandResult, 0, len(cmd1.Commands))
		for _, subCmd := range cmd1.Commands {
			apiCmd, err := convCommandToAPI(subCmd)
			if err != nil {
				return nil, err
			}

			apiCmdResults = append(apiCmdResults, apiCmd)
		}

		return &v1alpha1.CommandResult{
			Result: &v1alpha1.CommandResult_Commit{
				Commit: &v1alpha1.CommitResult{
					Results: apiCmdResults,
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown command type %T", cmd)
	}
}

func convStateCtxsToAPI(ss []*flowstate.StateCtx) []*v1alpha1.StateCtx {
	apiS := make([]*v1alpha1.StateCtx, 0, len(ss))
	for _, stateCtx := range ss {
		apiS = append(apiS, convStateCtxToAPI(stateCtx))
	}
	return apiS
}

func convStateCtxToAPI(s *flowstate.StateCtx) *v1alpha1.StateCtx {
	return &v1alpha1.StateCtx{
		Committed:   convStateToAPI(s.Committed),
		Current:     convStateToAPI(s.Current),
		Transitions: convTransitionsToAPI(s.Transitions),
	}
}

func convStateCtxToRefAPI(s *flowstate.StateCtx) *v1alpha1.StateCtxRef {
	return &v1alpha1.StateCtxRef{
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
