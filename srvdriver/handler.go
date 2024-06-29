package srvdriver

import (
	"context"
	"log"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1alpha1"
	v1alpha1 "github.com/makasim/flowstatesrv/protogen/flowstate/flow/v1alpha1"
)

type Handler struct {
	e *flowstate.Engine
}

func newHandler(e *flowstate.Engine) *Handler {
	return &Handler{
		e: e,
	}
}

func (h *Handler) Execute(_ context.Context, req *connect.Request[v1alpha1.ExecuteRequest]) (*connect.Response[v1alpha1.ExecuteResponse], error) {
	stateCtx := convertorv1alpha1.ConvertAPIToStateCtx(req.Msg.StateContext)
	resStateCtx := stateCtx.CopyTo(&flowstate.StateCtx{})
	go func() {
		if err := h.e.Execute(stateCtx); err != nil {
			log.Printf("ERROR: engine: execute: %s", err)
		}
	}()

	noopCmd, err := convertorv1alpha1.CommandToAPICommand(flowstate.End(resStateCtx))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1alpha1.ExecuteResponse{
		StateContext: convertorv1alpha1.ConvertStateCtxToAPI(resStateCtx),
		Command:      noopCmd,
	}), nil
}
