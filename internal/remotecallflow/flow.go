package remotecallflow

import (
	"bytes"
	"io"
	"net/http"

	"github.com/makasim/flowstate"
	"github.com/makasim/flowstatesrv/convertorv1alpha1"
	flowv1alpha1 "github.com/makasim/flowstatesrv/protogen/flowstate/flow/v1alpha1"
	"google.golang.org/protobuf/encoding/protojson"
)

type Config struct {
}

type Flow struct {
	callURL string
	hc      *http.Client
}

func New(callURL string) *Flow {
	return &Flow{
		callURL: callURL,
		hc:      &http.Client{},
	}
}

func (f *Flow) Execute(stateCtx *flowstate.StateCtx, e *flowstate.Engine) (flowstate.Command, error) {

	apiStateCtx := convertorv1alpha1.ConvertStateCtxToAPI(stateCtx)

	apiReq := &flowv1alpha1.ExecuteRequest{
		StateContext: apiStateCtx,
	}

	b, err := protojson.Marshal(apiReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, f.callURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	resp, err := f.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &flowv1alpha1.ExecuteResponse{}
	if err := protojson.Unmarshal(b, res); err != nil {
		return nil, err
	}

	resStateCtx := convertorv1alpha1.ConvertAPIToStateCtx(res.StateContext)

	cmd, err := convertorv1alpha1.ConvertAPIToCommand(res.Command, []*flowstate.StateCtx{resStateCtx})
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
