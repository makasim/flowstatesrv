// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: flowstate/v1/server.proto

package flowstatev1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/makasim/flowstatesrv/protogen/flowstate/v1"
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
	// ServerServiceName is the fully-qualified name of the ServerService service.
	ServerServiceName = "flowstate.v1.ServerService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ServerServiceDoCommandProcedure is the fully-qualified name of the ServerService's DoCommand RPC.
	ServerServiceDoCommandProcedure = "/flowstate.v1.ServerService/DoCommand"
	// ServerServiceWatchStatesProcedure is the fully-qualified name of the ServerService's WatchStates
	// RPC.
	ServerServiceWatchStatesProcedure = "/flowstate.v1.ServerService/WatchStates"
	// ServerServiceRegisterFlowProcedure is the fully-qualified name of the ServerService's
	// RegisterFlow RPC.
	ServerServiceRegisterFlowProcedure = "/flowstate.v1.ServerService/RegisterFlow"
)

// ServerServiceClient is a client for the flowstate.v1.ServerService service.
type ServerServiceClient interface {
	DoCommand(context.Context, *connect.Request[v1.DoCommandRequest]) (*connect.Response[v1.DoCommandResponse], error)
	WatchStates(context.Context, *connect.Request[v1.WatchStatesRequest]) (*connect.ServerStreamForClient[v1.WatchStatesResponse], error)
	RegisterFlow(context.Context, *connect.Request[v1.RegisterFlowRequest]) (*connect.Response[v1.RegisterFlowResponse], error)
}

// NewServerServiceClient constructs a client for the flowstate.v1.ServerService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewServerServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ServerServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &serverServiceClient{
		doCommand: connect.NewClient[v1.DoCommandRequest, v1.DoCommandResponse](
			httpClient,
			baseURL+ServerServiceDoCommandProcedure,
			opts...,
		),
		watchStates: connect.NewClient[v1.WatchStatesRequest, v1.WatchStatesResponse](
			httpClient,
			baseURL+ServerServiceWatchStatesProcedure,
			opts...,
		),
		registerFlow: connect.NewClient[v1.RegisterFlowRequest, v1.RegisterFlowResponse](
			httpClient,
			baseURL+ServerServiceRegisterFlowProcedure,
			opts...,
		),
	}
}

// serverServiceClient implements ServerServiceClient.
type serverServiceClient struct {
	doCommand    *connect.Client[v1.DoCommandRequest, v1.DoCommandResponse]
	watchStates  *connect.Client[v1.WatchStatesRequest, v1.WatchStatesResponse]
	registerFlow *connect.Client[v1.RegisterFlowRequest, v1.RegisterFlowResponse]
}

// DoCommand calls flowstate.v1.ServerService.DoCommand.
func (c *serverServiceClient) DoCommand(ctx context.Context, req *connect.Request[v1.DoCommandRequest]) (*connect.Response[v1.DoCommandResponse], error) {
	return c.doCommand.CallUnary(ctx, req)
}

// WatchStates calls flowstate.v1.ServerService.WatchStates.
func (c *serverServiceClient) WatchStates(ctx context.Context, req *connect.Request[v1.WatchStatesRequest]) (*connect.ServerStreamForClient[v1.WatchStatesResponse], error) {
	return c.watchStates.CallServerStream(ctx, req)
}

// RegisterFlow calls flowstate.v1.ServerService.RegisterFlow.
func (c *serverServiceClient) RegisterFlow(ctx context.Context, req *connect.Request[v1.RegisterFlowRequest]) (*connect.Response[v1.RegisterFlowResponse], error) {
	return c.registerFlow.CallUnary(ctx, req)
}

// ServerServiceHandler is an implementation of the flowstate.v1.ServerService service.
type ServerServiceHandler interface {
	DoCommand(context.Context, *connect.Request[v1.DoCommandRequest]) (*connect.Response[v1.DoCommandResponse], error)
	WatchStates(context.Context, *connect.Request[v1.WatchStatesRequest], *connect.ServerStream[v1.WatchStatesResponse]) error
	RegisterFlow(context.Context, *connect.Request[v1.RegisterFlowRequest]) (*connect.Response[v1.RegisterFlowResponse], error)
}

// NewServerServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewServerServiceHandler(svc ServerServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	serverServiceDoCommandHandler := connect.NewUnaryHandler(
		ServerServiceDoCommandProcedure,
		svc.DoCommand,
		opts...,
	)
	serverServiceWatchStatesHandler := connect.NewServerStreamHandler(
		ServerServiceWatchStatesProcedure,
		svc.WatchStates,
		opts...,
	)
	serverServiceRegisterFlowHandler := connect.NewUnaryHandler(
		ServerServiceRegisterFlowProcedure,
		svc.RegisterFlow,
		opts...,
	)
	return "/flowstate.v1.ServerService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ServerServiceDoCommandProcedure:
			serverServiceDoCommandHandler.ServeHTTP(w, r)
		case ServerServiceWatchStatesProcedure:
			serverServiceWatchStatesHandler.ServeHTTP(w, r)
		case ServerServiceRegisterFlowProcedure:
			serverServiceRegisterFlowHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedServerServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedServerServiceHandler struct{}

func (UnimplementedServerServiceHandler) DoCommand(context.Context, *connect.Request[v1.DoCommandRequest]) (*connect.Response[v1.DoCommandResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("flowstate.v1.ServerService.DoCommand is not implemented"))
}

func (UnimplementedServerServiceHandler) WatchStates(context.Context, *connect.Request[v1.WatchStatesRequest], *connect.ServerStream[v1.WatchStatesResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("flowstate.v1.ServerService.WatchStates is not implemented"))
}

func (UnimplementedServerServiceHandler) RegisterFlow(context.Context, *connect.Request[v1.RegisterFlowRequest]) (*connect.Response[v1.RegisterFlowResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("flowstate.v1.ServerService.RegisterFlow is not implemented"))
}