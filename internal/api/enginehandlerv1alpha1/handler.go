package enginehandlerv1alpha1

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1"
	commandv1 "github.com/makasim/flowstatesrv/protogen/flowstate/command/v1"
	v1alpha1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1alpha1"
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
	stateCtxs := make([]*flowstate.StateCtx, 0, len(req.Msg.StateContexts))
	for _, apiS := range req.Msg.StateContexts {
		stateCtxs = append(stateCtxs, convertorv1.ConvertAPIToStateCtx(apiS))
	}
	datas := make([]*flowstate.Data, 0, len(req.Msg.Data))
	for _, apiD := range req.Msg.Data {
		datas = append(datas, convertorv1.ConvertAPIToData(apiD))
	}

	cmds := make([]flowstate.Command, 0, len(req.Msg.Commands))
	for _, apiC := range req.Msg.Commands {
		cmd, err := convertorv1.APICommandToCommand(apiC, stateCtxs, datas)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		cmds = append(cmds, cmd)
	}

	conflictErr := &flowstate.ErrCommitConflict{}
	if err := s.e.Do(cmds...); errors.As(err, conflictErr) {
		apiConflictErr := &v1alpha1.ErrorConflict{}
		for _, stateID := range conflictErr.TaskIDs() {
			apiConflictErr.CommittableStateIds = append(apiConflictErr.CommittableStateIds, string(stateID))
		}
		ed, edErr := connect.NewErrorDetail(apiConflictErr)
		if edErr != nil {
			return nil, connect.NewError(connect.CodeInternal, edErr)
		}

		connErr := connect.NewError(connect.CodeAborted, err)
		connErr.AddDetail(ed)

		return nil, connErr
	} else if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	results := make([]*commandv1.AnyResult, 0, len(cmds))
	for _, cmd := range cmds {
		cmdRes, err := convertorv1.CommandToAPIResult(cmd)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		results = append(results, cmdRes)
	}

	return connect.NewResponse(&v1alpha1.DoResponse{
		StateContexts: convertorv1.ConvertStateCtxsToAPI(stateCtxs),
		Data:          convertorv1.ConvertDatasToAPI(datas),
		Results:       results,
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
			apiS := convertorv1.ConvertStateToAPI(state)
			if err := stream.Send(&v1alpha1.WatchResponse{
				State: apiS,
			}); err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
