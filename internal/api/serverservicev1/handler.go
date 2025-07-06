package serverservicev1

import (
	"context"
	"errors"
	"net/http"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1"
	"github.com/makasim/flowstatesrv/internal/remotecallflow"
	v1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1"
	"github.com/makasim/flowstatesrv/protogen/flowstate/v1/flowstatev1connect"
)

type Service struct {
	e flowstate.Engine
	d flowstate.Driver
}

func New(e flowstate.Engine, d flowstate.Driver) *Service {
	return &Service{
		e: e,
		d: d,
	}
}

func (s *Service) DoCommand(_ context.Context, req *connect.Request[v1.DoCommandRequest]) (*connect.Response[v1.DoCommandResponse], error) {
	if err := protovalidate.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

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

	err := s.e.Do(cmds...)
	if revMismatchErr := (&flowstate.ErrRevMismatch{}); errors.As(err, revMismatchErr) {
		apiRevMismatchErr := &v1.ErrorRevMismatch{}
		for _, stateID := range revMismatchErr.TaskIDs() {
			apiRevMismatchErr.CommittableStateIds = append(apiRevMismatchErr.CommittableStateIds, string(stateID))
		}
		ed, edErr := connect.NewErrorDetail(apiRevMismatchErr)
		if edErr != nil {
			return nil, connect.NewError(connect.CodeInternal, edErr)
		}

		connErr := connect.NewError(connect.CodeAborted, err)
		connErr.AddDetail(ed)

		return nil, connErr
	} else if errors.Is(err, flowstate.ErrNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	} else if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	results := make([]*v1.Result, 0, len(cmds))
	for _, cmd := range cmds {
		cmdRes, err := convertorv1.CommandToAPIResult(cmd)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		results = append(results, cmdRes)
	}

	return connect.NewResponse(&v1.DoCommandResponse{
		StateContexts: convertorv1.ConvertStateCtxsToAPI(stateCtxs),
		Data:          convertorv1.ConvertDatasToAPI(datas),
		Results:       results,
	}), nil
}

func (s *Service) RegisterFlow(_ context.Context, req *connect.Request[v1.RegisterFlowRequest]) (*connect.Response[v1.RegisterFlowResponse], error) {
	fc := flowstatev1connect.NewFlowServiceClient(http.DefaultClient, req.Msg.HttpHost)

	if err := s.d.SetFlow(
		flowstate.FlowID(req.Msg.FlowId),
		remotecallflow.New(fc),
	); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.RegisterFlowResponse{}), nil
}
