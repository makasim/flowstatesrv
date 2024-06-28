package flowhandlerv1alpha1

import (
	"context"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/internal/remotecallflow"
	v1alpha1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1alpha1"
)

type flowRegistry interface {
	SetFlow(id flowstate.FlowID, f flowstate.Flow)
}

type Handler struct {
	fr flowRegistry
}

func New(fr flowRegistry) *Handler {
	return &Handler{
		fr: fr,
	}
}

func (s *Handler) Register(_ context.Context, req *connect.Request[v1alpha1.RegisterRequest]) (*connect.Response[v1alpha1.RegisterResponse], error) {
	callURL := req.Msg.HttpHost + `/flowstate.v1alpha1.FlowService/Register`

	s.fr.SetFlow(
		flowstate.FlowID(req.Msg.FlowId),
		remotecallflow.New(callURL),
	)

	return connect.NewResponse(&v1alpha1.RegisterResponse{}), nil
}
