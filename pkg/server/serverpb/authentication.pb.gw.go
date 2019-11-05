// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: server/serverpb/authentication.proto

/*
Package serverpb is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package serverpb

import (
	"io"
	"net/http"

	"context"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray

func request_LogIn_UserLogin_0(ctx context.Context, marshaler runtime.Marshaler, client LogInClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq UserLoginRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.UserLogin(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func request_LogOut_UserLogout_0(ctx context.Context, marshaler runtime.Marshaler, client LogOutClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq UserLogoutRequest
	var metadata runtime.ServerMetadata

	msg, err := client.UserLogout(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

// RegisterLogInHandlerFromEndpoint is same as RegisterLogInHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterLogInHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterLogInHandler(ctx, mux, conn)
}

// RegisterLogInHandler registers the http handlers for service LogIn to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterLogInHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterLogInHandlerClient(ctx, mux, NewLogInClient(conn))
}

// RegisterLogInHandler registers the http handlers for service LogIn to "mux".
// The handlers forward requests to the grpc endpoint over the given implementation of "LogInClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "LogInClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "LogInClient" to call the correct interceptors.
func RegisterLogInHandlerClient(ctx context.Context, mux *runtime.ServeMux, client LogInClient) error {

	mux.Handle("POST", pattern_LogIn_UserLogin_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_LogIn_UserLogin_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_LogIn_UserLogin_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_LogIn_UserLogin_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"login"}, ""))
)

var (
	forward_LogIn_UserLogin_0 = runtime.ForwardResponseMessage
)

// RegisterLogOutHandlerFromEndpoint is same as RegisterLogOutHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterLogOutHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterLogOutHandler(ctx, mux, conn)
}

// RegisterLogOutHandler registers the http handlers for service LogOut to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterLogOutHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterLogOutHandlerClient(ctx, mux, NewLogOutClient(conn))
}

// RegisterLogOutHandler registers the http handlers for service LogOut to "mux".
// The handlers forward requests to the grpc endpoint over the given implementation of "LogOutClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "LogOutClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "LogOutClient" to call the correct interceptors.
func RegisterLogOutHandlerClient(ctx context.Context, mux *runtime.ServeMux, client LogOutClient) error {

	mux.Handle("GET", pattern_LogOut_UserLogout_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		if cn, ok := w.(http.CloseNotifier); ok {
			go func(done <-chan struct{}, closed <-chan bool) {
				select {
				case <-done:
				case <-closed:
					cancel()
				}
			}(ctx.Done(), cn.CloseNotify())
		}
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_LogOut_UserLogout_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_LogOut_UserLogout_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_LogOut_UserLogout_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"logout"}, ""))
)

var (
	forward_LogOut_UserLogout_0 = runtime.ForwardResponseMessage
)