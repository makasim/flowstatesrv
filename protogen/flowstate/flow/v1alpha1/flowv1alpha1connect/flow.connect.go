// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: flowstate/flow/v1alpha1/flow.proto

package flowv1alpha1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1alpha1 "github.com/makasim/flowstatesrv/protogen/flowstate/flow/v1alpha1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// FlowServiceName is the fully-qualified name of the FlowService service.
	FlowServiceName = "flowstate.flow.v1alpha1.FlowService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// FlowServiceExecuteProcedure is the fully-qualified name of the FlowService's Execute RPC.
	FlowServiceExecuteProcedure = "/flowstate.flow.v1alpha1.FlowService/Execute"
)

// FlowServiceClient is a client for the flowstate.flow.v1alpha1.FlowService service.
type FlowServiceClient interface {
	Execute(context.Context, *connect.Request[v1alpha1.ExecuteRequest]) (*connect.Response[v1alpha1.ExecuteResponse], error)
}

// NewFlowServiceClient constructs a client for the flowstate.flow.v1alpha1.FlowService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewFlowServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) FlowServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &flowServiceClient{
		execute: connect.NewClient[v1alpha1.ExecuteRequest, v1alpha1.ExecuteResponse](
			httpClient,
			baseURL+FlowServiceExecuteProcedure,
			opts...,
		),
	}
}

// flowServiceClient implements FlowServiceClient.
type flowServiceClient struct {
	execute *connect.Client[v1alpha1.ExecuteRequest, v1alpha1.ExecuteResponse]
}

// Execute calls flowstate.flow.v1alpha1.FlowService.Execute.
func (c *flowServiceClient) Execute(ctx context.Context, req *connect.Request[v1alpha1.ExecuteRequest]) (*connect.Response[v1alpha1.ExecuteResponse], error) {
	return c.execute.CallUnary(ctx, req)
}

// FlowServiceHandler is an implementation of the flowstate.flow.v1alpha1.FlowService service.
type FlowServiceHandler interface {
	Execute(context.Context, *connect.Request[v1alpha1.ExecuteRequest]) (*connect.Response[v1alpha1.ExecuteResponse], error)
}

// NewFlowServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewFlowServiceHandler(svc FlowServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	flowServiceExecuteHandler := connect.NewUnaryHandler(
		FlowServiceExecuteProcedure,
		svc.Execute,
		opts...,
	)
	return "/flowstate.flow.v1alpha1.FlowService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case FlowServiceExecuteProcedure:
			flowServiceExecuteHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedFlowServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedFlowServiceHandler struct{}

func (UnimplementedFlowServiceHandler) Execute(context.Context, *connect.Request[v1alpha1.ExecuteRequest]) (*connect.Response[v1alpha1.ExecuteResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("flowstate.flow.v1alpha1.FlowService.Execute is not implemented"))
}