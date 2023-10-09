// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: postman/games/wheel/wheel_service.proto

package wheel

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// WheelServiceClient is the client API for WheelService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WheelServiceClient interface {
	JoinWheelRound(ctx context.Context, in *JoinWheelRoundRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type wheelServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWheelServiceClient(cc grpc.ClientConnInterface) WheelServiceClient {
	return &wheelServiceClient{cc}
}

func (c *wheelServiceClient) JoinWheelRound(ctx context.Context, in *JoinWheelRoundRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/wheel_service.WheelService/JoinWheelRound", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WheelServiceServer is the server API for WheelService service.
// All implementations must embed UnimplementedWheelServiceServer
// for forward compatibility
type WheelServiceServer interface {
	JoinWheelRound(context.Context, *JoinWheelRoundRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedWheelServiceServer()
}

// UnimplementedWheelServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWheelServiceServer struct {
}

func (UnimplementedWheelServiceServer) JoinWheelRound(context.Context, *JoinWheelRoundRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinWheelRound not implemented")
}
func (UnimplementedWheelServiceServer) mustEmbedUnimplementedWheelServiceServer() {}

// UnsafeWheelServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WheelServiceServer will
// result in compilation errors.
type UnsafeWheelServiceServer interface {
	mustEmbedUnimplementedWheelServiceServer()
}

func RegisterWheelServiceServer(s grpc.ServiceRegistrar, srv WheelServiceServer) {
	s.RegisterService(&WheelService_ServiceDesc, srv)
}

func _WheelService_JoinWheelRound_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinWheelRoundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WheelServiceServer).JoinWheelRound(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wheel_service.WheelService/JoinWheelRound",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WheelServiceServer).JoinWheelRound(ctx, req.(*JoinWheelRoundRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WheelService_ServiceDesc is the grpc.ServiceDesc for WheelService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WheelService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "wheel_service.WheelService",
	HandlerType: (*WheelServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "JoinWheelRound",
			Handler:    _WheelService_JoinWheelRound_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "postman/games/wheel/wheel_service.proto",
}