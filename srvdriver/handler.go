package srvdriver

import (
	"context"
	"log"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1"
	clientv1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1"
)

type handler struct {
	e flowstate.Engine
}

func newHandler(e flowstate.Engine) *handler {
	return &handler{
		e: e,
	}
}

func (h *handler) Execute(_ context.Context, req *connect.Request[clientv1.ExecuteRequest]) (*connect.Response[clientv1.ExecuteResponse], error) {
	stateCtx := convertorv1.ConvertAPIToStateCtx(req.Msg.StateContext)
	resStateCtx := stateCtx.CopyTo(&flowstate.StateCtx{})
	go func() {
		if err := h.e.Execute(stateCtx); err != nil {
			log.Printf("ERROR: engine: execute: %s", err)
		}
	}()

	noopCmd, err := convertorv1.CommandToAPICommand(flowstate.End(resStateCtx))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&clientv1.ExecuteResponse{
		StateContext: convertorv1.ConvertStateCtxToAPI(resStateCtx),
		Command:      noopCmd,
	}), nil
}
