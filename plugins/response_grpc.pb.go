// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package plugins

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ResponseClient is the client API for Response service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ResponseClient interface {
	Modify(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Output, error)
}

type responseClient struct {
	cc grpc.ClientConnInterface
}

func NewResponseClient(cc grpc.ClientConnInterface) ResponseClient {
	return &responseClient{cc}
}

func (c *responseClient) Modify(ctx context.Context, in *Input, opts ...grpc.CallOption) (*Output, error) {
	out := new(Output)
	err := c.cc.Invoke(ctx, "/ncs.protobuf.Response/Modify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResponseServer is the server API for Response service.
// All implementations must embed UnimplementedResponseServer
// for forward compatibility
type ResponseServer interface {
	Modify(context.Context, *Input) (*Output, error)
	mustEmbedUnimplementedResponseServer()
}

// UnimplementedResponseServer must be embedded to have forward compatible implementations.
type UnimplementedResponseServer struct {
}

func (UnimplementedResponseServer) Modify(context.Context, *Input) (*Output, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Modify not implemented")
}
func (UnimplementedResponseServer) mustEmbedUnimplementedResponseServer() {}

// UnsafeResponseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ResponseServer will
// result in compilation errors.
type UnsafeResponseServer interface {
	mustEmbedUnimplementedResponseServer()
}

func RegisterResponseServer(s grpc.ServiceRegistrar, srv ResponseServer) {
	s.RegisterService(&Response_ServiceDesc, srv)
}

func _Response_Modify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Input)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResponseServer).Modify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ncs.protobuf.Response/Modify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResponseServer).Modify(ctx, req.(*Input))
	}
	return interceptor(ctx, in, info, handler)
}

// Response_ServiceDesc is the grpc.ServiceDesc for Response service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Response_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ncs.protobuf.Response",
	HandlerType: (*ResponseServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Modify",
			Handler:    _Response_Modify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "response.proto",
}
