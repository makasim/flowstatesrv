package remotecallflow

import (
	"context"

	"connectrpc.com/connect"
	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1"
	clientv1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1"
	"github.com/makasim/flowstatesrv/protogen/flowstate/v1/flowstatev1connect"
)

type Config struct {
}

type Flow struct {
	fc flowstatev1connect.FlowServiceClient
}

func New(fc flowstatev1connect.FlowServiceClient) *Flow {
	return &Flow{
		fc: fc,
	}
}

func (f *Flow) Execute(stateCtx *flowstate.StateCtx, _ flowstate.Engine) (flowstate.Command, error) {
	apiStateCtx := convertorv1.ConvertStateCtxToAPI(stateCtx)

	resp, err := f.fc.Execute(context.Background(), connect.NewRequest(&clientv1.ExecuteRequest{
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
