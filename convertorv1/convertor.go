package convertorv1

import (
	"fmt"
	"slices"
	"time"

	"github.com/makasim/flowstate"
	flowstatev1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1"
)

func FindStateCtxByRef(ref *flowstatev1.StateRef, stateCtxs []*flowstate.StateCtx) (*flowstate.StateCtx, error) {
	for _, stateCtx := range stateCtxs {
		if string(stateCtx.Current.ID) == ref.Id && stateCtx.Current.Rev == ref.Rev {
			return stateCtx, nil
		}
	}

	return nil, fmt.Errorf("there is no state ctx provided for ref: %s:%d", ref.Id, ref.Rev)
}

func FindDataByRef(ref *flowstatev1.DataRef, datas []*flowstate.Data) (*flowstate.Data, error) {
	for _, d := range datas {
		if d.ID == flowstate.DataID(ref.Id) && d.Rev == ref.Rev {
			return d, nil
		}
	}

	return nil, fmt.Errorf("there is no data provided for ref: %s:%d", ref.Id, ref.Rev)
}

func APICommandToCommand(apiAnyCmd *flowstatev1.AnyCommand, stateCtxs []*flowstate.StateCtx, datas []*flowstate.Data) (flowstate.Command, error) {
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

		if apiCmd.FlowId != `` {
			cmd = cmd.WithTransit(flowstate.FlowID(apiCmd.FlowId))
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

		if apiCmd.Commit {
			cmd.WithCommit(true)
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
			cmtCmd, err := APICommandToCommand(apiCmtCmd, stateCtxs, datas)
			if err != nil {
				return nil, err
			}

			cmtCmds = append(cmtCmds, cmtCmd)
		}

		return flowstate.Commit(cmtCmds...), nil

	case apiAnyCmd.GetStoreData() != nil:
		apiCmd := apiAnyCmd.GetStoreData()

		d, err := FindDataByRef(apiCmd.DataRef, datas)
		if err != nil {
			return nil, err
		}

		return flowstate.StoreData(d), nil
	case apiAnyCmd.GetGetData() != nil:
		apiCmd := apiAnyCmd.GetGetData()

		d, err := FindDataByRef(apiCmd.DataRef, datas)
		if err != nil {
			return nil, err
		}

		return flowstate.GetData(d), nil
	case apiAnyCmd.GetReferenceData() != nil:
		apiCmd := apiAnyCmd.GetReferenceData()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}
		d, err := FindDataByRef(apiCmd.DataRef, datas)
		if err != nil {
			return nil, err
		}

		return flowstate.ReferenceData(stateCtx, d, apiCmd.Annotation), nil
	case apiAnyCmd.GetDereferenceData() != nil:
		apiCmd := apiAnyCmd.GetDereferenceData()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}
		d, err := FindDataByRef(apiCmd.DataRef, datas)
		if err != nil {
			return nil, err
		}

		return flowstate.DereferenceData(stateCtx, d, apiCmd.Annotation), nil
	case apiAnyCmd.GetGet() != nil:
		apiCmd := apiAnyCmd.GetGet()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		if apiCmd.Id != "" {
			return flowstate.GetByID(stateCtx, flowstate.StateID(apiCmd.Id), apiCmd.Rev), nil
		}

		return flowstate.GetByLabels(stateCtx, copyMap(apiCmd.Labels)), nil
	case apiAnyCmd.GetGetMany() != nil:
		apiCmd := apiAnyCmd.GetGetMany()

		cmd := &flowstate.GetManyCommand{
			SinceRev:   apiCmd.SinceRev,
			SinceTime:  time.UnixMicro(apiCmd.SinceTimeUsec),
			LatestOnly: apiCmd.LatestOnly,
			Limit:      int(apiCmd.Limit),
		}

		for _, labels := range apiCmd.GetLabels() {
			if len(labels.GetLabels()) == 0 {
				continue
			}

			cmd.Labels = append(cmd.Labels, copyMap(labels.Labels))
		}

		return cmd, nil
	case apiAnyCmd.GetCommitState() != nil:
		apiCmd := apiAnyCmd.GetCommitState()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.CommitStateCtx(stateCtx), nil
	default:
		return nil, fmt.Errorf("unknown command %T", apiAnyCmd.Command)
	}
}

func ConvertAPIToStateCtx(apiS *flowstatev1.StateContext) *flowstate.StateCtx {
	return &flowstate.StateCtx{
		Current:     ConvertAPIToState(apiS.Current),
		Committed:   ConvertAPIToState(apiS.Committed),
		Transitions: ConvertAPIToTransitions(apiS.Transitions),
	}
}

func ConvertAPIToStateCtxs(apiS []*flowstatev1.StateContext) []*flowstate.StateCtx {
	stateCtxs := make([]*flowstate.StateCtx, 0, len(apiS))
	for _, apiS := range apiS {
		stateCtxs = append(stateCtxs, ConvertAPIToStateCtx(apiS))
	}
	return stateCtxs
}

func ConvertAPIToState(apiS *flowstatev1.State) flowstate.State {
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

func ConvertAPIToStates(apiS []*flowstatev1.State) []flowstate.State {
	ss := make([]flowstate.State, 0, len(apiS))
	for _, apiS := range apiS {
		ss = append(ss, ConvertAPIToState(apiS))
	}
	return ss
}

func ConvertAPIToTransitions(apiTs []*flowstatev1.Transition) []flowstate.Transition {
	ts := make([]flowstate.Transition, 0, len(apiTs))
	for _, apiT := range apiTs {
		ts = append(ts, ConvertAPIToTransition(apiT))
	}
	return ts
}

func ConvertAPIToTransition(apiT *flowstatev1.Transition) flowstate.Transition {
	if apiT == nil {
		return flowstate.Transition{}
	}

	return flowstate.Transition{
		FromID:      flowstate.FlowID(apiT.From),
		ToID:        flowstate.FlowID(apiT.To),
		Annotations: copyMap(apiT.Annotations),
	}
}

func ConvertAPIToDatas(apiDatas []*flowstatev1.Data) []*flowstate.Data {
	datas := make([]*flowstate.Data, 0, len(apiDatas))
	for _, apiD := range apiDatas {
		datas = append(datas, ConvertAPIToData(apiD))
	}
	return datas
}

func CommandToAPICommand(cmd flowstate.Command) (*flowstatev1.AnyCommand, error) {
	switch cmd1 := cmd.(type) {
	case *flowstate.TransitCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_Transit{
				Transit: &flowstatev1.Transit{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
					FlowId:   string(cmd1.FlowID),
				},
			},
		}, nil
	case *flowstate.PauseCommand:
		apiCmd := &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_Pause{
				Pause: &flowstatev1.Pause{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}

		if cmd1.FlowID != "" {
			apiCmd.GetPause().FlowId = string(cmd1.FlowID)
		}

		return apiCmd, nil
	case *flowstate.ResumeCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_Resume{
				Resume: &flowstatev1.Resume{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.EndCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_End{
				End: &flowstatev1.End{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.ExecuteCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_Execute{
				Execute: &flowstatev1.Execute{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.DelayCommand:
		apiCmd := &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_Delay{
				Delay: &flowstatev1.Delay{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
					Duration: cmd1.Duration.String(),
				},
			},
		}

		if cmd1.Commit {
			apiCmd.GetDelay().Commit = true
		}

		return apiCmd, nil
	case *flowstate.NoopCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_Noop{
				Noop: &flowstatev1.Noop{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.SerializeCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_Serialize{
				Serialize: &flowstatev1.Serialize{
					SerializableStateRef: ConvertStateCtxToRefAPI(cmd1.SerializableStateCtx),
					StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
					Annotation:           cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.DeserializeCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_Deserialize{
				Deserialize: &flowstatev1.Deserialize{
					StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
					DeserializedStateRef: ConvertStateCtxToRefAPI(cmd1.DeserializedStateCtx),
					Annotation:           cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.CommitCommand:
		apiCmtCmds := make([]*flowstatev1.AnyCommand, 0, len(cmd1.Commands))
		for _, cmtCmd := range cmd1.Commands {
			apiCmtCmd, err := CommandToAPICommand(cmtCmd)
			if err != nil {
				return nil, err
			}

			apiCmtCmds = append(apiCmtCmds, apiCmtCmd)
		}

		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_Commit{
				Commit: &flowstatev1.Commit{
					Commands: apiCmtCmds,
				},
			},
		}, nil
	case *flowstate.StoreDataCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_StoreData{
				StoreData: &flowstatev1.StoreData{
					DataRef: ConvertDataToRefAPI(cmd1.Data),
				},
			},
		}, nil
	case *flowstate.GetDataCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_GetData{
				GetData: &flowstatev1.GetData{
					DataRef: ConvertDataToRefAPI(cmd1.Data),
				},
			},
		}, nil
	case *flowstate.ReferenceDataCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_ReferenceData{
				ReferenceData: &flowstatev1.ReferenceData{
					StateRef:   ConvertStateCtxToRefAPI(cmd1.StateCtx),
					DataRef:    ConvertDataToRefAPI(cmd1.Data),
					Annotation: cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.DereferenceDataCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_DereferenceData{
				DereferenceData: &flowstatev1.DereferenceData{
					StateRef:   ConvertStateCtxToRefAPI(cmd1.StateCtx),
					DataRef:    ConvertDataToRefAPI(cmd1.Data),
					Annotation: cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.GetCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_Get{
				Get: &flowstatev1.Get{
					Id:       string(cmd1.ID),
					Rev:      cmd1.Rev,
					Labels:   copyMap(cmd1.Labels),
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.GetManyCommand:
		var apiLabels []*flowstatev1.GetMany_Labels
		for _, labels := range cmd1.Labels {
			apiLabels = append(apiLabels, &flowstatev1.GetMany_Labels{
				Labels: copyMap(labels),
			})
		}

		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_GetMany{
				GetMany: &flowstatev1.GetMany{
					SinceRev:      cmd1.SinceRev,
					SinceTimeUsec: cmd1.SinceTime.UnixMicro(),
					LatestOnly:    cmd1.LatestOnly,
					Labels:        apiLabels,
					Limit:         int64(cmd1.Limit),
				},
			},
		}, nil
	case *flowstate.CommitStateCtxCommand:
		return &flowstatev1.AnyCommand{
			Command: &flowstatev1.AnyCommand_CommitState{
				CommitState: &flowstatev1.CommitState{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown command type %T", cmd)
	}
}

func CommandToAPIResult(cmd flowstate.Command) (*flowstatev1.AnyResult, error) {
	switch cmd1 := cmd.(type) {
	case *flowstate.TransitCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_Transit{
				Transit: &flowstatev1.TransitResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.PauseCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_Pause{
				Pause: &flowstatev1.PauseResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.ResumeCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_Resume{
				Resume: &flowstatev1.ResumeResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.EndCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_End{
				End: &flowstatev1.EndResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.ExecuteCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_Execute{
				Execute: &flowstatev1.ExecuteResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.DelayCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_Delay{
				Delay: &flowstatev1.DelayResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
					Duration: cmd1.Duration.String(),
					Commit:   cmd1.Commit,
				},
			},
		}, nil
	case *flowstate.NoopCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_Noop{
				Noop: &flowstatev1.NoopResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.SerializeCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_Serialize{
				Serialize: &flowstatev1.SerializeResult{
					SerializableStateRef: ConvertStateCtxToRefAPI(cmd1.SerializableStateCtx),
					StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
					Annotation:           cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.DeserializeCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_Deserialize{
				Deserialize: &flowstatev1.DeserializeResult{
					StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
					DeserializedStateRef: ConvertStateCtxToRefAPI(cmd1.DeserializedStateCtx),
					Annotation:           cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.CommitCommand:
		apiCmtResults := make([]*flowstatev1.AnyResult, 0, len(cmd1.Commands))
		for _, subCmd := range cmd1.Commands {
			apiCmtRes, err := CommandToAPIResult(subCmd)
			if err != nil {
				return nil, err
			}

			apiCmtResults = append(apiCmtResults, apiCmtRes)
		}

		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_Commit{
				Commit: &flowstatev1.CommitResult{
					Results: apiCmtResults,
				},
			},
		}, nil
	case *flowstate.StoreDataCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_StoreData{
				StoreData: &flowstatev1.StoreDataResult{
					DataRef: ConvertDataToRefAPI(cmd1.Data),
				},
			},
		}, nil
	case *flowstate.GetDataCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_GetData{
				GetData: &flowstatev1.GetDataResult{
					DataRef: ConvertDataToRefAPI(cmd1.Data),
				},
			},
		}, nil
	case *flowstate.ReferenceDataCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_ReferenceData{
				ReferenceData: &flowstatev1.ReferenceDataResult{
					StateRef:   ConvertStateCtxToRefAPI(cmd1.StateCtx),
					DataRef:    ConvertDataToRefAPI(cmd1.Data),
					Annotation: cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.DereferenceDataCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_DereferenceData{
				DereferenceData: &flowstatev1.DereferenceDataResult{
					StateRef:   ConvertStateCtxToRefAPI(cmd1.StateCtx),
					DataRef:    ConvertDataToRefAPI(cmd1.Data),
					Annotation: cmd1.Annotation,
				},
			},
		}, nil
	case *flowstate.GetCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_Get{
				Get: &flowstatev1.GetResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	case *flowstate.GetManyCommand:
		res, err := cmd1.Result()
		if err != nil {
			return nil, err
		}

		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_GetMany{
				GetMany: &flowstatev1.GetManyResult{
					States: ConvertStatesToAPI(res.States),
					More:   res.More,
				},
			},
		}, nil
	case *flowstate.CommitStateCtxCommand:
		return &flowstatev1.AnyResult{
			Result: &flowstatev1.AnyResult_CommitState{
				CommitState: &flowstatev1.CommitStateResult{
					StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown command type %T", cmd)
	}
}

func ConvertCommandsToAPI(cmds []flowstate.Command) ([]*flowstatev1.AnyCommand, error) {
	apiCmds := make([]*flowstatev1.AnyCommand, 0, len(cmds))
	for _, cmd := range cmds {
		apiCmd, err := CommandToAPICommand(cmd)
		if err != nil {
			return nil, err
		}

		apiCmds = append(apiCmds, apiCmd)
	}

	return apiCmds, nil
}

func ConvertCommandToStateContexts(cmd flowstate.Command) []*flowstatev1.StateContext {
	apiStateCtxs := make([]*flowstatev1.StateContext, 0)

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
	case *flowstate.ReferenceDataCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.DereferenceDataCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.GetCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.CommitStateCtxCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	default:
		return nil
	}

	slices.CompactFunc(apiStateCtxs, func(l, r *flowstatev1.StateContext) bool {
		return l.Current.Id == r.Current.Id && l.Current.Rev == r.Current.Rev
	})

	return apiStateCtxs
}

func ConvertStateCtxsToAPI(ss []*flowstate.StateCtx) []*flowstatev1.StateContext {
	apiS := make([]*flowstatev1.StateContext, 0, len(ss))
	for _, stateCtx := range ss {
		apiS = append(apiS, ConvertStateCtxToAPI(stateCtx))
	}
	return apiS
}

func ConvertStateCtxToAPI(s *flowstate.StateCtx) *flowstatev1.StateContext {
	return &flowstatev1.StateContext{
		Committed:   ConvertStateToAPI(s.Committed),
		Current:     ConvertStateToAPI(s.Current),
		Transitions: ConvertTransitionsToAPI(s.Transitions),
	}
}

func ConvertStateCtxToRefAPI(s *flowstate.StateCtx) *flowstatev1.StateRef {
	return &flowstatev1.StateRef{
		Id:  string(s.Current.ID),
		Rev: s.Current.Rev,
	}
}

func ConvertStateToAPI(s flowstate.State) *flowstatev1.State {
	return &flowstatev1.State{
		Id:                   string(s.ID),
		Rev:                  s.Rev,
		CommittedAtUnixMilli: s.CommittedAtUnixMilli,
		Transition:           ConvertTransitionToAPI(s.Transition),
		Labels:               copyMap(s.Labels),
		Annotations:          copyMap(s.Annotations),
	}
}

func ConvertStatesToAPI(ss []flowstate.State) []*flowstatev1.State {
	apiSS := make([]*flowstatev1.State, 0, len(ss))
	for _, s := range ss {
		apiSS = append(apiSS, ConvertStateToAPI(s))
	}
	return apiSS
}

func ConvertTransitionsToAPI(tss []flowstate.Transition) []*flowstatev1.Transition {
	apiTss := make([]*flowstatev1.Transition, 0, len(tss))
	for _, ts := range tss {
		apiTss = append(apiTss, ConvertTransitionToAPI(ts))
	}
	return apiTss
}

func ConvertTransitionToAPI(ts flowstate.Transition) *flowstatev1.Transition {
	return &flowstatev1.Transition{
		From:        string(ts.FromID),
		To:          string(ts.ToID),
		Annotations: copyMap(ts.Annotations),
	}
}

func ConvertCommandToDatas(cmd flowstate.Command) []*flowstatev1.Data {
	apiDatas := make([]*flowstatev1.Data, 0)

	switch cmd1 := cmd.(type) {
	case *flowstate.StoreDataCommand:
		apiDatas = append(apiDatas, ConvertDataToAPI(cmd1.Data))
	case *flowstate.GetDataCommand:
		apiDatas = append(apiDatas, ConvertDataToAPI(cmd1.Data))
	case *flowstate.ReferenceDataCommand:
		apiDatas = append(apiDatas, ConvertDataToAPI(cmd1.Data))
	case *flowstate.DereferenceDataCommand:
		apiDatas = append(apiDatas, ConvertDataToAPI(cmd1.Data))
	case *flowstate.CommitCommand:
		for _, subCmd := range cmd1.Commands {
			apiDatas = append(apiDatas, ConvertCommandToDatas(subCmd)...)
		}
	default:
		return nil
	}

	slices.CompactFunc(apiDatas, func(l, r *flowstatev1.Data) bool {
		return l.Id == r.Id && l.Rev == r.Rev
	})

	return apiDatas
}

func ConvertAPIToData(data *flowstatev1.Data) *flowstate.Data {
	return &flowstate.Data{
		ID:     flowstate.DataID(data.Id),
		Rev:    data.Rev,
		Binary: data.Binary,
		B:      append([]byte(nil), data.B...),
	}
}

func ConvertDataToAPI(data *flowstate.Data) *flowstatev1.Data {
	return &flowstatev1.Data{
		Id:     string(data.ID),
		Rev:    data.Rev,
		Binary: data.Binary,
		B:      string(data.B),
	}
}

func ConvertDatasToAPI(datas []*flowstate.Data) []*flowstatev1.Data {
	apiD := make([]*flowstatev1.Data, 0, len(datas))
	for _, d := range datas {
		apiD = append(apiD, ConvertDataToAPI(d))
	}
	return apiD
}

func ConvertDataToRefAPI(data *flowstate.Data) *flowstatev1.DataRef {
	return &flowstatev1.DataRef{
		Id:  string(data.ID),
		Rev: data.Rev,
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
