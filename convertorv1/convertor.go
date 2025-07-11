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

func APICommandToCommand(apiAnyCmd *flowstatev1.Command, stateCtxs []*flowstate.StateCtx, datas []*flowstate.Data) (flowstate.Command, error) {
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

		cmd := flowstate.DelayUntil(stateCtx, time.Unix(apiCmd.ExecuteAtSec, 0))

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
	case apiAnyCmd.GetGetStateById() != nil:
		apiCmd := apiAnyCmd.GetGetStateById()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.GetStateByID(stateCtx, flowstate.StateID(apiCmd.Id), apiCmd.Rev), nil
	case apiAnyCmd.GetGetStateByLabels() != nil:
		apiCmd := apiAnyCmd.GetGetStateByLabels()

		stateCtx, err := FindStateCtxByRef(apiCmd.StateRef, stateCtxs)
		if err != nil {
			return nil, err
		}

		return flowstate.GetStateByLabels(stateCtx, copyMap(apiCmd.Labels)), nil
	case apiAnyCmd.GetGetStates() != nil:
		apiCmd := apiAnyCmd.GetGetStates()

		var orLabels []map[string]string
		for _, apiLabel := range apiCmd.GetLabels() {
			if len(apiLabel.GetLabels()) == 0 {
				continue
			}

			orLabels = append(orLabels, copyMap(apiLabel.Labels))
		}

		cmd := &flowstate.GetStatesCommand{
			SinceRev:   apiCmd.SinceRev,
			Labels:     orLabels,
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
		return nil, fmt.Errorf("at least one know command should be set")
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

func ConvertAPIToDelayedState(apiDelayedState *flowstatev1.DelayedState) flowstate.DelayedState {
	if apiDelayedState == nil {
		return flowstate.DelayedState{}
	}

	return flowstate.DelayedState{
		State:     ConvertAPIToState(apiDelayedState.State),
		Offset:    apiDelayedState.Offset,
		ExecuteAt: time.Unix(apiDelayedState.ExecuteAtSec, 0),
	}
}

func ConvertAPIToDelayedStates(apiDelayedStates []*flowstatev1.DelayedState) []flowstate.DelayedState {
	delayedStates := make([]flowstate.DelayedState, 0, len(apiDelayedStates))
	for _, apiDelayedState := range apiDelayedStates {
		delayedStates = append(delayedStates, ConvertAPIToDelayedState(apiDelayedState))
	}
	return delayedStates
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
		From:        flowstate.FlowID(apiT.From),
		To:          flowstate.FlowID(apiT.To),
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

func CommandToAPICommand(cmd flowstate.Command) (*flowstatev1.Command, error) {
	switch cmd1 := cmd.(type) {
	case *flowstate.TransitCommand:
		return &flowstatev1.Command{
			Transit: &flowstatev1.Transit{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
				FlowId:   string(cmd1.FlowID),
			},
		}, nil
	case *flowstate.PauseCommand:
		apiCmd := &flowstatev1.Command{
			Pause: &flowstatev1.Pause{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}

		if cmd1.FlowID != "" {
			apiCmd.GetPause().FlowId = string(cmd1.FlowID)
		}

		return apiCmd, nil
	case *flowstate.ResumeCommand:
		return &flowstatev1.Command{
			Resume: &flowstatev1.Resume{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.EndCommand:
		return &flowstatev1.Command{
			End: &flowstatev1.End{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.ExecuteCommand:
		return &flowstatev1.Command{
			Execute: &flowstatev1.Execute{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.DelayCommand:
		apiCmd := &flowstatev1.Command{
			Delay: &flowstatev1.Delay{
				StateRef:     ConvertStateCtxToRefAPI(cmd1.StateCtx),
				ExecuteAtSec: cmd1.ExecuteAt.Unix(),
				Commit:       cmd1.Commit,
			},
		}

		return apiCmd, nil
	case *flowstate.NoopCommand:
		return &flowstatev1.Command{
			Noop: &flowstatev1.Noop{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.SerializeCommand:
		return &flowstatev1.Command{
			Serialize: &flowstatev1.Serialize{
				SerializableStateRef: ConvertStateCtxToRefAPI(cmd1.SerializableStateCtx),
				StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
				Annotation:           cmd1.Annotation,
			},
		}, nil
	case *flowstate.DeserializeCommand:
		return &flowstatev1.Command{
			Deserialize: &flowstatev1.Deserialize{
				StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
				DeserializedStateRef: ConvertStateCtxToRefAPI(cmd1.DeserializedStateCtx),
				Annotation:           cmd1.Annotation,
			},
		}, nil
	case *flowstate.CommitCommand:
		apiCmtCmds := make([]*flowstatev1.Command, 0, len(cmd1.Commands))
		for _, cmtCmd := range cmd1.Commands {
			apiCmtCmd, err := CommandToAPICommand(cmtCmd)
			if err != nil {
				return nil, err
			}

			apiCmtCmds = append(apiCmtCmds, apiCmtCmd)
		}

		return &flowstatev1.Command{
			Commit: &flowstatev1.Commit{
				Commands: apiCmtCmds,
			},
		}, nil
	case *flowstate.StoreDataCommand:
		return &flowstatev1.Command{
			StoreData: &flowstatev1.StoreData{
				DataRef: ConvertDataToRefAPI(cmd1.Data),
			},
		}, nil
	case *flowstate.GetDataCommand:
		return &flowstatev1.Command{
			GetData: &flowstatev1.GetData{
				DataRef: ConvertDataToRefAPI(cmd1.Data),
			},
		}, nil
	case *flowstate.ReferenceDataCommand:
		return &flowstatev1.Command{
			ReferenceData: &flowstatev1.ReferenceData{
				StateRef:   ConvertStateCtxToRefAPI(cmd1.StateCtx),
				DataRef:    ConvertDataToRefAPI(cmd1.Data),
				Annotation: cmd1.Annotation,
			},
		}, nil
	case *flowstate.DereferenceDataCommand:
		return &flowstatev1.Command{
			DereferenceData: &flowstatev1.DereferenceData{
				StateRef:   ConvertStateCtxToRefAPI(cmd1.StateCtx),
				DataRef:    ConvertDataToRefAPI(cmd1.Data),
				Annotation: cmd1.Annotation,
			},
		}, nil
	case *flowstate.GetStateByIDCommand:
		return &flowstatev1.Command{
			GetStateById: &flowstatev1.GetStateByID{
				Id:       string(cmd1.ID),
				Rev:      cmd1.Rev,
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.GetStateByLabelsCommand:
		return &flowstatev1.Command{
			GetStateByLabels: &flowstatev1.GetStateByLabels{
				Labels:   copyMap(cmd1.Labels),
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.GetStatesCommand:
		var apiLabels []*flowstatev1.GetStates_Labels
		for _, labels := range cmd1.Labels {
			apiLabels = append(apiLabels, &flowstatev1.GetStates_Labels{
				Labels: copyMap(labels),
			})
		}

		return &flowstatev1.Command{
			GetStates: &flowstatev1.GetStates{
				SinceRev:      cmd1.SinceRev,
				SinceTimeUsec: cmd1.SinceTime.UnixMicro(),
				LatestOnly:    cmd1.LatestOnly,
				Labels:        apiLabels,
				Limit:         int64(cmd1.Limit),
			},
		}, nil
	case *flowstate.CommitStateCtxCommand:
		return &flowstatev1.Command{
			CommitState: &flowstatev1.CommitState{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown command type %T", cmd)
	}
}

func CommandToAPIResult(cmd flowstate.Command) (*flowstatev1.Result, error) {
	switch cmd1 := cmd.(type) {
	case *flowstate.TransitCommand:
		return &flowstatev1.Result{
			Transit: &flowstatev1.TransitResult{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.PauseCommand:
		return &flowstatev1.Result{
			Pause: &flowstatev1.PauseResult{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.ResumeCommand:
		return &flowstatev1.Result{
			Resume: &flowstatev1.ResumeResult{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.EndCommand:
		return &flowstatev1.Result{
			End: &flowstatev1.EndResult{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.ExecuteCommand:
		return &flowstatev1.Result{
			Execute: &flowstatev1.ExecuteResult{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.DelayCommand:
		return &flowstatev1.Result{
			Delay: &flowstatev1.DelayResult{
				StateRef:      ConvertStateCtxToRefAPI(cmd1.StateCtx),
				DelayingState: ConvertStateToAPI(cmd1.DelayingState),
			},
		}, nil
	case *flowstate.NoopCommand:
		return &flowstatev1.Result{
			Noop: &flowstatev1.NoopResult{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.SerializeCommand:
		return &flowstatev1.Result{
			Serialize: &flowstatev1.SerializeResult{
				SerializableStateRef: ConvertStateCtxToRefAPI(cmd1.SerializableStateCtx),
				StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
				Annotation:           cmd1.Annotation,
			},
		}, nil
	case *flowstate.DeserializeCommand:
		return &flowstatev1.Result{
			Deserialize: &flowstatev1.DeserializeResult{
				StateRef:             ConvertStateCtxToRefAPI(cmd1.StateCtx),
				DeserializedStateRef: ConvertStateCtxToRefAPI(cmd1.DeserializedStateCtx),
				Annotation:           cmd1.Annotation,
			},
		}, nil
	case *flowstate.CommitCommand:
		apiCmtResults := make([]*flowstatev1.Result, 0, len(cmd1.Commands))
		for _, subCmd := range cmd1.Commands {
			apiCmtRes, err := CommandToAPIResult(subCmd)
			if err != nil {
				return nil, err
			}

			apiCmtResults = append(apiCmtResults, apiCmtRes)
		}

		return &flowstatev1.Result{
			Commit: &flowstatev1.CommitResult{
				Results: apiCmtResults,
			},
		}, nil
	case *flowstate.StoreDataCommand:
		return &flowstatev1.Result{
			StoreData: &flowstatev1.StoreDataResult{
				DataRef: ConvertDataToRefAPI(cmd1.Data),
			},
		}, nil
	case *flowstate.GetDataCommand:
		return &flowstatev1.Result{
			GetData: &flowstatev1.GetDataResult{
				DataRef: ConvertDataToRefAPI(cmd1.Data),
			},
		}, nil
	case *flowstate.ReferenceDataCommand:
		return &flowstatev1.Result{
			ReferenceData: &flowstatev1.ReferenceDataResult{
				StateRef:   ConvertStateCtxToRefAPI(cmd1.StateCtx),
				DataRef:    ConvertDataToRefAPI(cmd1.Data),
				Annotation: cmd1.Annotation,
			},
		}, nil
	case *flowstate.DereferenceDataCommand:
		return &flowstatev1.Result{
			DereferenceData: &flowstatev1.DereferenceDataResult{
				StateRef:   ConvertStateCtxToRefAPI(cmd1.StateCtx),
				DataRef:    ConvertDataToRefAPI(cmd1.Data),
				Annotation: cmd1.Annotation,
			},
		}, nil
	case *flowstate.GetStateByIDCommand:
		return &flowstatev1.Result{
			GetStateById: &flowstatev1.GetStateByIDResult{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.GetStateByLabelsCommand:
		return &flowstatev1.Result{
			GetStateByLabels: &flowstatev1.GetStateByLabelsResult{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	case *flowstate.GetStatesCommand:
		if cmd1.Result == nil {
			return nil, fmt.Errorf("command result is nil for GetStatesCommand")
		}

		return &flowstatev1.Result{
			GetStates: &flowstatev1.GetStatesResult{
				States: ConvertStatesToAPI(cmd1.Result.States),
				More:   cmd1.Result.More,
			},
		}, nil
	case *flowstate.CommitStateCtxCommand:
		return &flowstatev1.Result{
			CommitState: &flowstatev1.CommitStateResult{
				StateRef: ConvertStateCtxToRefAPI(cmd1.StateCtx),
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown command type %T", cmd)
	}
}

func ConvertCommandsToAPI(cmds []flowstate.Command) ([]*flowstatev1.Command, error) {
	apiCmds := make([]*flowstatev1.Command, 0, len(cmds))
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
	case *flowstate.GetStateByIDCommand:
		apiStateCtxs = append(apiStateCtxs, ConvertStateCtxToAPI(cmd1.StateCtx))
	case *flowstate.GetStateByLabelsCommand:
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

func ConvertDelayedStatesToAPI(delayedStates []flowstate.DelayedState) []*flowstatev1.DelayedState {
	apiDelayedStates := make([]*flowstatev1.DelayedState, 0, len(delayedStates))
	for _, ds := range delayedStates {
		apiDelayedStates = append(apiDelayedStates, ConvertDelayedStateToAPI(ds))
	}
	return apiDelayedStates
}

func ConvertDelayedStateToAPI(ds flowstate.DelayedState) *flowstatev1.DelayedState {
	return &flowstatev1.DelayedState{
		State:        ConvertStateToAPI(ds.State),
		Offset:       ds.Offset,
		ExecuteAtSec: ds.ExecuteAt.Unix(),
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
		From:        string(ts.From),
		To:          string(ts.To),
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
