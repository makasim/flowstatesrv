package convertorv1alpha1

import (
	"fmt"
	"slices"
	"time"

	"github.com/makasim/flowstate"
	"github.com/makasim/flowstate/exptcmd"
	v1alpha1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1alpha1"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

func FindStateCtxByRef(ref *v1alpha1.StateRef, stateCtxs []*flowstate.StateCtx) (*flowstate.StateCtx, error) {
	for _, stateCtx := range stateCtxs {
		if string(stateCtx.Current.ID) == ref.Id && stateCtx.Current.Rev == ref.Rev {
			return stateCtx, nil
		}
	}

	return nil, fmt.Errorf("there is no state ctx provided for ref: %s:%d", ref.Id, ref.Rev)
}

func APICommandToCommand(apiC *anypb.Any, stateCtxs []*flowstate.StateCtx) (flowstate.Command, error) {
	switch apiC.TypeUrl {
	case `type.googleapis.com/flowstate.v1alpha1.Transit`:
		apiCmd := &v1alpha1.Transit{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Transit(stateCtx, flowstate.FlowID(apiCmd.FlowId)), nil
	case `type.googleapis.com/flowstate.v1alpha1.Pause`:
		apiCmd := &v1alpha1.Pause{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		cmd := flowstate.Pause(stateCtx)
		if apiCmd.FlowId != "" {
			cmd = cmd.WithTransit(flowstate.FlowID(apiCmd.FlowId))
		}

		return cmd, nil
	case `type.googleapis.com/flowstate.v1alpha1.Resume`:
		apiCmd := &v1alpha1.Resume{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Resume(stateCtx), nil
	case `type.googleapis.com/flowstate.v1alpha1.End`:
		apiCmd := &v1alpha1.End{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.End(stateCtx), nil
	case `type.googleapis.com/flowstate.v1alpha1.Execute`:
		apiCmd := &v1alpha1.Execute{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Execute(stateCtx), nil
	case `type.googleapis.com/flowstate.v1alpha1.Delay`:
		apiCmd := &v1alpha1.Delay{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
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
	case `type.googleapis.com/flowstate.v1alpha1.Stack`:
		apiCmd := &v1alpha1.Stack{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stackedStateCtx, err := FindStateCtxByRef(apiCmd.StackedStateRef, stateCtxs)
		if err != nil {
			return nil, err
		}
		nextStateCtx, err := FindStateCtxByRef(apiCmd.NextStateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return exptcmd.Stack(stackedStateCtx, nextStateCtx), nil
	case `type.googleapis.com/flowstate.v1alpha1.Unstack`:
		apiCmd := &v1alpha1.Unstack{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}
		unstackStateCtx, err := FindStateCtxByRef(apiCmd.UnstackStateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return exptcmd.Unstack(stateCtx, unstackStateCtx), nil
	case `type.googleapis.com/flowstate.v1alpha1.Commit`:
		apiCmd := &v1alpha1.Commit{}
		if err := apiC.UnmarshalTo(apiCmd); err != nil {
			return nil, err
		}

		subCmds := make([]flowstate.Command, 0, len(apiCmd.Commands))
		for _, subCmd := range apiCmd.Commands {
			subCmd, err := APICommandToCommand(subCmd, stateCtxs)
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

func ConvertAPIToStateCtx(apiS *v1alpha1.StateContext) *flowstate.StateCtx {
	return &flowstate.StateCtx{
		Current:     ConvertAPIToState(apiS.Current),
		Committed:   ConvertAPIToState(apiS.Committed),
		Transitions: ConvertAPIToTransitions(apiS.Transitions),
	}
}

func ConvertAPIToStateCtxs(apiS []*v1alpha1.StateContext) []*flowstate.StateCtx {
	stateCtxs := make([]*flowstate.StateCtx, 0, len(apiS))
	for _, apiS := range apiS {
		stateCtxs = append(stateCtxs, ConvertAPIToStateCtx(apiS))
	}
	return stateCtxs
}

func ConvertAPIToState(apiS *v1alpha1.State) flowstate.State {
	if apiS == nil {
		return flowstate.State{}
	}

	return flowstate.State{
		ID:                   flowstate.StateID(apiS.Id),
		Rev:                  apiS.Rev,
		Annotations:          copyMap(apiS.Annotations),
		Labels:               copyMap(apiS.Labels),
		CommittedAtUnixMilli: apiS.CommittedAtUnixMilli,
		Transition:           ConvertAPIToTransition(apiS.Transition),
	}
}

func ConvertAPIToTransitions(apiTs []*v1alpha1.Transition) []flowstate.Transition {
	ts := make([]flowstate.Transition, 0, len(apiTs))
	for _, apiT := range apiTs {
		ts = append(ts, ConvertAPIToTransition(apiT))
	}
	return ts
}

func ConvertAPIToTransition(apiT *v1alpha1.Transition) flowstate.Transition {
	if apiT == nil {
		return flowstate.Transition{}
	}

	return flowstate.Transition{
		FromID:      flowstate.FlowID(apiT.From),
		ToID:        flowstate.FlowID(apiT.To),
		Annotations: copyMap(apiT.Annotations),
	}
}

func CommandToAPICommand(cmd flowstate.Command) (*anypb.Any, error) {
	switch cmd1 := cmd.(type) {
	case *flowstate.TransitCommand:
		return anypb.New(&v1alpha1.Transit{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			FlowId:   string(cmd1.FlowID),
		})
	case *flowstate.PauseCommand:
		return anypb.New(&v1alpha1.Pause{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			FlowId:   string(cmd1.FlowID),
		})
	case *flowstate.ResumeCommand:
		return anypb.New(&v1alpha1.Resume{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.EndCommand:
		return anypb.New(&v1alpha1.End{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.ExecuteCommand:
		return anypb.New(&v1alpha1.Execute{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.DelayCommand:
		return anypb.New(&v1alpha1.Delay{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			Duration: cmd1.Duration.String(),
			Commit:   cmd1.Commit,
		})
	case *flowstate.NoopCommand:
		return anypb.New(&v1alpha1.Noop{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *exptcmd.StackCommand:
		return anypb.New(&v1alpha1.Stack{
			StackedStateRef: ConvertStateCtxToRefAPI(cmd1.StackedStateCtx),
			NextStateRef:    ConvertStateCtxToRefAPI(cmd1.NextStateCtx),
		})
	case *exptcmd.UnstackCommand:
		return anypb.New(&v1alpha1.Unstack{
			StateRef:        ConvertStateCtxToRefAPI(cmd1.StateCtx),
			UnstackStateRef: ConvertStateCtxToRefAPI(cmd1.UnstackStateCtx),
		})
	case *flowstate.CommitCommand:
		subCmds := make([]*anypb.Any, 0, len(cmd1.Commands))
		for _, subCmd := range cmd1.Commands {
			subRes, err := CommandToAPICommand(subCmd)
			if err != nil {
				return nil, err
			}

			subCmds = append(subCmds, subRes)
		}

		return anypb.New(&v1alpha1.Commit{
			Commands: subCmds,
		})
	default:
		return nil, fmt.Errorf("unknown command type %T", cmd)
	}
}

func CommandToAPIResult(cmd flowstate.Command) (*anypb.Any, error) {
	switch cmd1 := cmd.(type) {
	case *flowstate.TransitCommand:
		return anypb.New(&v1alpha1.TransitResult{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.PauseCommand:
		return anypb.New(&v1alpha1.PauseResult{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.ResumeCommand:
		return anypb.New(&v1alpha1.ResumeResult{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.EndCommand:
		return anypb.New(&v1alpha1.EndResult{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.ExecuteCommand:
		return anypb.New(&v1alpha1.ExecuteResult{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *flowstate.DelayCommand:
		return anypb.New(&v1alpha1.DelayResult{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			Duration: cmd1.Duration.String(),
			Commit:   cmd1.Commit,
		})
	case *flowstate.NoopCommand:
		return anypb.New(&v1alpha1.NoopResult{
			StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
		})
	case *exptcmd.StackCommand:
		return anypb.New(&v1alpha1.StackResult{
			StackedStateRef: ConvertStateCtxToRefAPI(cmd1.StackedStateCtx),
			NextStateRef:    ConvertStateCtxToRefAPI(cmd1.NextStateCtx),
		})
	case *exptcmd.UnstackCommand:
		return anypb.New(&v1alpha1.UnstackResult{
			StateRef:        ConvertStateCtxToRefAPI(cmd1.StateCtx),
			UnstackStateRef: ConvertStateCtxToRefAPI(cmd1.UnstackStateCtx),
		})
	case *flowstate.CommitCommand:
		subResults := make([]*anypb.Any, 0, len(cmd1.Commands))
		for _, subCmd := range cmd1.Commands {
			subRes, err := CommandToAPIResult(subCmd)
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

func ConvertCommandsToAPI(cmds []flowstate.Command) ([]*anypb.Any, error) {
	apiCmds := make([]*anypb.Any, 0, len(cmds))
	for _, cmd := range cmds {
		apiCmd, err := CommandToAPICommand(cmd)
		if err != nil {
			return nil, err
		}

		apiCmds = append(apiCmds, apiCmd)
	}

	return apiCmds, nil
}

func ConvertCommandToStateContexts(cmd flowstate.Command) []*v1alpha1.StateContext {
	apiStateCtxs := make([]*v1alpha1.StateContext, 0)

	switch cmd1 := cmd.(type) {
	case *flowstate.TransitCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.PauseCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.ResumeCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.EndCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.ExecuteCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.DelayCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.NoopCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	case *exptcmd.StackCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StackedStateCtx), ConvertStateCtxToAPI(cmd1.NextStateCtx))
	case *exptcmd.UnstackCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx), ConvertStateCtxToAPI(cmd1.UnstackStateCtx))
	case *flowstate.CommitCommand:
		for _, subCmd := range cmd1.Commands {
			apiStateCtxs = append(apiStateCtxs, ConvertCommandToStateContexts(subCmd)...)
		}
	default:
		return nil
	}

	slices.CompactFunc(apiStateCtxs, func(l, r *v1alpha1.StateContext) bool {
		return l.Current.Id == r.Current.Id && l.Current.Rev == r.Current.Rev
	})

	return apiStateCtxs
}

func ConvertStateCtxsToAPI(ss []*flowstate.StateCtx) []*v1alpha1.StateContext {
	apiS := make([]*v1alpha1.StateContext, 0, len(ss))
	for _, stateCtx := range ss {
		apiS = append(apiS, ConvertStateCtxToAPI(stateCtx))
	}
	return apiS
}

func ConvertStateCtxToAPI(s *flowstate.StateCtx) *v1alpha1.StateContext {
	return &v1alpha1.StateContext{
		Committed:   ConvertStateToAPI(s.Committed),
		Current:     ConvertStateToAPI(s.Current),
		Transitions: ConvertTransitionsToAPI(s.Transitions),
	}
}

func ConvertStateCtxToRefAPI(s *flowstate.StateCtx) *v1alpha1.StateRef {
	return &v1alpha1.StateRef{
		Id:  string(s.Current.ID),
		Rev: s.Current.Rev,
	}
}

func ConvertStateToAPI(s flowstate.State) *v1alpha1.State {
	return &v1alpha1.State{
		Id:                   string(s.ID),
		Rev:                  s.Rev,
		CommittedAtUnixMilli: s.CommittedAtUnixMilli,
		Transition:           ConvertTransitionToAPI(s.Transition),
		Labels:               copyMap(s.Labels),
		Annotations:          copyMap(s.Annotations),
	}
}

func ConvertTransitionsToAPI(tss []flowstate.Transition) []*v1alpha1.Transition {
	apiTss := make([]*v1alpha1.Transition, 0, len(tss))
	for _, ts := range tss {
		apiTss = append(apiTss, ConvertTransitionToAPI(ts))
	}
	return apiTss
}

func ConvertTransitionToAPI(ts flowstate.Transition) *v1alpha1.Transition {
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
