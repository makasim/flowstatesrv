package flowhandlerv1alpha1

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/internal/remotecallflow"
	"github.com/makasim/flowstatesrv/protogen/flowstate/flow/v1alpha1/flowv1alpha1connect"
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
	fc := flowv1alpha1connect.NewFlowServiceClient(http.DefaultClient, req.Msg.HttpHost)

	s.fr.SetFlow(
		flowstate.FlowID(req.Msg.FlowId),
		remotecallflow.New(fc),
	)

	return connect.NewResponse(&v1alpha1.RegisterResponse{}), nil
}
