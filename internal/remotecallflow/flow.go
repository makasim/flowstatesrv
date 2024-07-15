package remotecallflow

import (
	"context"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1"
	flowv1alpha1 "github.com/makasim/flowstatesrv/protogen/flowstate/flow/v1alpha1"
	"github.com/makasim/flowstatesrv/protogen/flowstate/flow/v1alpha1/flowv1alpha1connect"
)

type Config struct {
}

type Flow struct {
	fc flowv1alpha1connect.FlowServiceClient
}

func New(fc flowv1alpha1connect.FlowServiceClient) *Flow {
	return &Flow{
		fc: fc,
	}
}

func (f *Flow) Execute(stateCtx *flowstate.StateCtx, _ *flowstate.Engine) (flowstate.Command, error) {
	apiStateCtx := convertorv1.ConvertStateCtxToAPI(stateCtx)

	resp, err := f.fc.Execute(context.Background(), connect.NewRequest(&flowv1alpha1.ExecuteRequest{
		StateContext: apiStateCtx,
	}))
	if err != nil {
		return nil, err
	}

	resStateCtx := convertorv1.ConvertAPIToStateCtx(resp.Msg.StateContext)
	resStateCtx.CopyTo(stateCtx)

	cmd, err := convertorv1.APICommandToCommand(resp.Msg.Command, []*flowstate.StateCtx{stateCtx}, nil)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
