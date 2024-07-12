package convertorv1

import (
	"fmt"
	"slices"
	"time"

	"github.com/makasim/flowstate"
	commandv1 "github.com/makasim/flowstatesrv/protogen/flowstate/command/v1"
	v1 "github.com/makasim/flowstatesrv/protogen/flowstate/state/v1"
)

func FindStateCtxByRef(ref *v1.StateRef, stateCtxs []*flowstate.StateCtx) (*flowstate.StateCtx, error) {
	for _, stateCtx := range stateCtxs {
		if string(stateCtx.Current.ID) == ref.Id && stateCtx.Current.Rev == ref.Rev {
			return stateCtx, nil
		}
	}

	return nil, fmt.Errorf("there is no state ctx provided for ref: %s:%d", ref.Id, ref.Rev)
}

func APICommandToCommand(apiAnyCmd *commandv1.AnyCommand, stateCtxs []*flowstate.StateCtx) (flowstate.Command, error) {
	switch {
	case apiAnyCmd.GetTransit() != nil:
		apiCmd := apiAnyCmd.GetTransit()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Transit(stateCtx, flowstate.FlowID(apiCmd.FlowId)), nil
	case apiAnyCmd.GetPause() != nil:
		apiCmd := apiAnyCmd.GetPause()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		cmd := flowstate.Pause(stateCtx)

		if apiCmd.WithTransit != nil {
			cmd = cmd.WithTransit(flowstate.FlowID(apiCmd.WithTransit.FlowId))
		}

		return cmd, nil
	case apiAnyCmd.GetResume() != nil:
		apiCmd := apiAnyCmd.GetResume()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Resume(stateCtx), nil
	case apiAnyCmd.GetEnd() != nil:
		apiCmd := apiAnyCmd.GetEnd()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.End(stateCtx), nil
	case apiAnyCmd.GetExecute() != nil:
		apiCmd := apiAnyCmd.GetExecute()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Execute(stateCtx), nil
	case apiAnyCmd.GetDelay() != nil:
		apiCmd := apiAnyCmd.GetDelay()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		dur, err := time.ParseDuration(apiCmd.Duration)
		if err != nil {
			return nil, err
		}

		cmd := flowstate.Delay(stateCtx, dur)

		if apiCmd.WithCommit != nil {
			cmd.Commit = apiCmd.WithCommit.Commit
		}

		return cmd, nil
	case apiAnyCmd.GetSerialize() != nil:
		apiCmd := apiAnyCmd.GetSerialize()

		serializableStateCtx, err := FindStateCtxByRef(apiCmd.SerializableStateRef, stateCtxs)
		if err != nil {
			return nil, err
		}
		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Serialize(serializableStateCtx, stateCtx, apiCmd.Annotation), nil
	case apiAnyCmd.GetDeserialize() != nil:
		apiCmd := apiAnyCmd.GetDeserialize()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}
		deserializedStateCtx, err := FindStateCtxByRef(apiCmd.DeserializedStateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.Deserialize(stateCtx, deserializedStateCtx, apiCmd.Annotation), nil
	case apiAnyCmd.GetCommit() != nil:
		apiCmd := apiAnyCmd.GetCommit()

		cmtCmds := make([]flowstate.Command, 0, len(apiCmd.Commands))
		for _, apiCmtCmd := range apiCmd.Commands {
			cmtCmd, err := APICommandToCommand(apiCmtCmd, stateCtxs)
			if err != nil {
				return nil, err
			}

			cmtCmds = append(cmtCmds, cmtCmd)
		}

		return flowstate.Commit(cmtCmds...), nil
	default:
		return nil, fmt.Errorf("unknown command %T", apiAnyCmd.Command)
	}
}

func ConvertAPIToStateCtx(apiS *v1.StateContext) *flowstate.StateCtx {
	return &flowstate.StateCtx{
		Current:     ConvertAPIToState(apiS.Current),
		Committed:   ConvertAPIToState(apiS.Committed),
		Transitions: ConvertAPIToTransitions(apiS.Transitions),
	}
}

func ConvertAPIToStateCtxs(apiS []*v1.StateContext) []*flowstate.StateCtx {
	stateCtxs := make([]*flowstate.StateCtx, 0, len(apiS))
	for _, apiS := range apiS {
		stateCtxs = append(stateCtxs, ConvertAPIToStateCtx(apiS))
	}
	return stateCtxs
}

func ConvertAPIToState(apiS *v1.State) flowstate.State {
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

func ConvertAPIToTransitions(apiTs []*v1.Transition) []flowstate.Transition {
	ts := make([]flowstate.Transition, 0, len(apiTs))
	for _, apiT := range apiTs {
		ts = append(ts, ConvertAPIToTransition(apiT))
	}
	return ts
}

func ConvertAPIToTransition(apiT *v1.Transition) flowstate.Transition {
	if apiT == nil {
		return flowstate.Transition{}
	}

	return flowstate.Transition{
		FromID:      flowstate.FlowID(apiT.From),
		ToID:        flowstate.FlowID(apiT.To),
		Annotations: copyMap(apiT.Annotations),
	}
}

func CommandToAPICommand(cmd flowstate.Command) (*commandv1.AnyCommand, error) {
	switch cmd1 := cmd.(type) {
	case *flowstate.TransitCommand:
		return &commandv1.AnyCommand{
			Command: &commandv1.AnyCommand_Transit{
				Transit: &commandv1.Transit{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
					FlowId:   string(cmd1.FlowID),
				},
			},
		}, nil
	case *flowstate.PauseCommand:
		apiCmd := &commandv1.AnyCommand{
			Command: &commandv1.AnyCommand_Pause{
				Pause: &commandv1.Pause{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}

		if cmd1.FlowID != "" {
			apiCmd.GetPause().WithTransit = &commandv1.Pause_WithTransit{
				FlowId: string(cmd1.FlowID),
			}
		}

		return apiCmd, nil
	case *flowstate.ResumeCommand:
		return &commandv1.AnyCommand{
			Command: &commandv1.AnyCommand_Resume{
				Resume: &commandv1.Resume{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.EndCommand:
		return &commandv1.AnyCommand{
			Command: &commandv1.AnyCommand_End{
				End: &commandv1.End{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.ExecuteCommand:
		return &commandv1.AnyCommand{
			Command: &commandv1.AnyCommand_Execute{
				Execute: &commandv1.Execute{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.DelayCommand:
		apiCmd := &commandv1.AnyCommand{
			Command: &commandv1.AnyCommand_Delay{
				Delay: &commandv1.Delay{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
					Duration: cmd1.Duration.String(),
				},
			},
		}

		if cmd1.Commit {
			apiCmd.GetDelay().WithCommit = &commandv1.Delay_WithCommit{
				Commit: true,
			}
		}

		return apiCmd, nil
	case *flowstate.NoopCommand:
		return &commandv1.AnyCommand{
			Command: &commandv1.AnyCommand_Noop{
				Noop: &commandv1.Noop{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.SerializeCommand:
		return &commandv1.AnyCommand{
			Command: &commandv1.AnyCommand_Serialize{
				Serialize: &commandv1.Serialize{
					SerializableStateRef: ConvertStateCtxToRefAPI(cmd1.SerializableStateCtx),
					StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
					Annotation:           cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.DeserializeCommand:
		return &commandv1.AnyCommand{
			Command: &commandv1.AnyCommand_Deserialize{
				Deserialize: &commandv1.Deserialize{
					StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
					DeserializedStateRef: ConvertStateCtxToRefAPI(cmd1.DeserializedStateCtx),
					Annotation:           cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.CommitCommand:
		apiCmtCmds := make([]*commandv1.AnyCommand, 0, len(cmd1.Commands))
		for _, cmtCmd := range cmd1.Commands {
			apiCmtCmd, err := CommandToAPICommand(cmtCmd)
			if err != nil {
				return nil, err
			}

			apiCmtCmds = append(apiCmtCmds, apiCmtCmd)
		}

		return &commandv1.AnyCommand{
			Command: &commandv1.AnyCommand_Commit{
				Commit: &commandv1.Commit{
					Commands: apiCmtCmds,
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown command type %T", cmd)
	}
}

func CommandToAPIResult(cmd flowstate.Command) (*commandv1.AnyResult, error) {
	switch cmd1 := cmd.(type) {
	case *flowstate.TransitCommand:
		return &commandv1.AnyResult{
			Result: &commandv1.AnyResult_Transit{
				Transit: &commandv1.TransitResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.PauseCommand:
		return &commandv1.AnyResult{
			Result: &commandv1.AnyResult_Pause{
				Pause: &commandv1.PauseResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.ResumeCommand:
		return &commandv1.AnyResult{
			Result: &commandv1.AnyResult_Resume{
				Resume: &commandv1.ResumeResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.EndCommand:
		return &commandv1.AnyResult{
			Result: &commandv1.AnyResult_End{
				End: &commandv1.EndResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.ExecuteCommand:
		return &commandv1.AnyResult{
			Result: &commandv1.AnyResult_Execute{
				Execute: &commandv1.ExecuteResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.DelayCommand:
		return &commandv1.AnyResult{
			Result: &commandv1.AnyResult_Delay{
				Delay: &commandv1.DelayResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
					Duration: cmd1.Duration.String(),
					Commit:   cmd1.Commit,
				},
			},
		}, nil
	case *flowstate.NoopCommand:
		return &commandv1.AnyResult{
			Result: &commandv1.AnyResult_Noop{
				Noop: &commandv1.NoopResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.SerializeCommand:
		return &commandv1.AnyResult{
			Result: &commandv1.AnyResult_Serialize{
				Serialize: &commandv1.SerializeResult{
					SerializableStateRef: ConvertStateCtxToRefAPI(cmd1.SerializableStateCtx),
					StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
					Annotation:           cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.DeserializeCommand:
		return &commandv1.AnyResult{
			Result: &commandv1.AnyResult_Deserialize{
				Deserialize: &commandv1.DeserializeResult{
					StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
					DeserializedStateRef: ConvertStateCtxToRefAPI(cmd1.DeserializedStateCtx),
					Annotation:           cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.CommitCommand:
		apiCmtResults := make([]*commandv1.AnyResult, 0, len(cmd1.Commands))
		for _, subCmd := range cmd1.Commands {
			apiCmtRes, err := CommandToAPIResult(subCmd)
			if err != nil {
				return nil, err
			}

			apiCmtResults = append(apiCmtResults, apiCmtRes)
		}

		return &commandv1.AnyResult{
			Result: &commandv1.AnyResult_Commit{
				Commit: &commandv1.CommitResult{
					Results: apiCmtResults,
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown command type %T", cmd)
	}
}

func ConvertCommandsToAPI(cmds []flowstate.Command) ([]*commandv1.AnyCommand, error) {
	apiCmds := make([]*commandv1.AnyCommand, 0, len(cmds))
	for _, cmd := range cmds {
		apiCmd, err := CommandToAPICommand(cmd)
		if err != nil {
			return nil, err
		}

		apiCmds = append(apiCmds, apiCmd)
	}

	return apiCmds, nil
}

func ConvertCommandToStateContexts(cmd flowstate.Command) []*v1.StateContext {
	apiStateCtxs := make([]*v1.StateContext, 0)

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
	case *flowstate.SerializeCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.SerializableStateCtx), ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.DeserializeCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.DeserializedStateCtx), ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.CommitCommand:
		for _, subCmd := range cmd1.Commands {
			apiStateCtxs = append(apiStateCtxs, ConvertCommandToStateContexts(subCmd)...)
		}
	default:
		return nil
	}

	slices.CompactFunc(apiStateCtxs, func(l, r *v1.StateContext) bool {
		return l.Current.Id == r.Current.Id && l.Current.Rev == r.Current.Rev
	})

	return apiStateCtxs
}

func ConvertStateCtxsToAPI(ss []*flowstate.StateCtx) []*v1.StateContext {
	apiS := make([]*v1.StateContext, 0, len(ss))
	for _, stateCtx := range ss {
		apiS = append(apiS, ConvertStateCtxToAPI(stateCtx))
	}
	return apiS
}

func ConvertStateCtxToAPI(s *flowstate.StateCtx) *v1.StateContext {
	return &v1.StateContext{
		Committed:   ConvertStateToAPI(s.Committed),
		Current:     ConvertStateToAPI(s.Current),
		Transitions: ConvertTransitionsToAPI(s.Transitions),
	}
}

func ConvertStateCtxToRefAPI(s *flowstate.StateCtx) *v1.StateRef {
	return &v1.StateRef{
		Id:  string(s.Current.ID),
		Rev: s.Current.Rev,
	}
}

func ConvertStateToAPI(s flowstate.State) *v1.State {
	return &v1.State{
		Id:                   string(s.ID),
		Rev:                  s.Rev,
		CommittedAtUnixMilli: s.CommittedAtUnixMilli,
		Transition:           ConvertTransitionToAPI(s.Transition),
		Labels:               copyMap(s.Labels),
		Annotations:          copyMap(s.Annotations),
	}
}

func ConvertTransitionsToAPI(tss []flowstate.Transition) []*v1.Transition {
	apiTss := make([]*v1.Transition, 0, len(tss))
	for _, ts := range tss {
		apiTss = append(apiTss, ConvertTransitionToAPI(ts))
	}
	return apiTss
}

func ConvertTransitionToAPI(ts flowstate.Transition) *v1.Transition {
	return &v1.Transition{
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
