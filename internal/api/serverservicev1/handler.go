package serverservicev1

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1"
	"github.com/makasim/flowstatesrv/internal/remotecallflow"
	"github.com/makasim/flowstatesrv/protogen/flowstate/client/v1/clientv1connect"
	v1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1"
)

type flowRegistry interface {
	SetFlow(id flowstate.FlowID, f flowstate.Flow)
}

type Service struct {
	e  flowstate.Engine
	fr flowRegistry
}

func New(e flowstate.Engine, fr flowRegistry) *Service {
	return &Service{
		e:  e,
		fr: fr,
	}
}

func (s *Service) DoCommand(_ context.Context, req *connect.Request[v1.DoCommandRequest]) (*connect.Response[v1.DoCommandResponse], error) {
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

	revMismatchErr := &flowstate.ErrRevMismatch{}
	if err := s.e.Do(cmds...); errors.As(err, revMismatchErr) {
		apiConflictErr := &v1.ErrorConflict{}
		for _, stateID := range revMismatchErr.TaskIDs() {
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

	results := make([]*v1.AnyResult, 0, len(cmds))
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
	fc := clientv1connect.NewClientServiceClient(http.DefaultClient, req.Msg.HttpHost)

	s.fr.SetFlow(
		flowstate.FlowID(req.Msg.FlowId),
		remotecallflow.New(fc),
	)

	return connect.NewResponse(&v1.RegisterFlowResponse{}), nil
}
