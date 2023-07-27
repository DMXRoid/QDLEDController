// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.3
// source: proto/led.proto

package v2

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

// LEDControllerClient is the client API for LEDController service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LEDControllerClient interface {
	GetLEDs(ctx context.Context, in *GetLEDsRequest, opts ...grpc.CallOption) (*GetLEDsResponse, error)
	UpdateLEDs(ctx context.Context, in *UpdateLEDsRequest, opts ...grpc.CallOption) (*UpdateLEDsResponse, error)
}

type lEDControllerClient struct {
	cc grpc.ClientConnInterface
}

func NewLEDControllerClient(cc grpc.ClientConnInterface) LEDControllerClient {
	return &lEDControllerClient{cc}
}

func (c *lEDControllerClient) GetLEDs(ctx context.Context, in *GetLEDsRequest, opts ...grpc.CallOption) (*GetLEDsResponse, error) {
	out := new(GetLEDsResponse)
	err := c.cc.Invoke(ctx, "/LEDController/GetLEDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lEDControllerClient) UpdateLEDs(ctx context.Context, in *UpdateLEDsRequest, opts ...grpc.CallOption) (*UpdateLEDsResponse, error) {
	out := new(UpdateLEDsResponse)
	err := c.cc.Invoke(ctx, "/LEDController/UpdateLEDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LEDControllerServer is the server API for LEDController service.
// All implementations must embed UnimplementedLEDControllerServer
// for forward compatibility
type LEDControllerServer interface {
	GetLEDs(context.Context, *GetLEDsRequest) (*GetLEDsResponse, error)
	UpdateLEDs(context.Context, *UpdateLEDsRequest) (*UpdateLEDsResponse, error)
	mustEmbedUnimplementedLEDControllerServer()
}

// UnimplementedLEDControllerServer must be embedded to have forward compatible implementations.
type UnimplementedLEDControllerServer struct {
}

func (UnimplementedLEDControllerServer) GetLEDs(context.Context, *GetLEDsRequest) (*GetLEDsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLEDs not implemented")
}
func (UnimplementedLEDControllerServer) UpdateLEDs(context.Context, *UpdateLEDsRequest) (*UpdateLEDsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLEDs not implemented")
}
func (UnimplementedLEDControllerServer) mustEmbedUnimplementedLEDControllerServer() {}

// UnsafeLEDControllerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LEDControllerServer will
// result in compilation errors.
type UnsafeLEDControllerServer interface {
	mustEmbedUnimplementedLEDControllerServer()
}

func RegisterLEDControllerServer(s grpc.ServiceRegistrar, srv LEDControllerServer) {
	s.RegisterService(&LEDController_ServiceDesc, srv)
}

func _LEDController_GetLEDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLEDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LEDControllerServer).GetLEDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/LEDController/GetLEDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LEDControllerServer).GetLEDs(ctx, req.(*GetLEDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LEDController_UpdateLEDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLEDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LEDControllerServer).UpdateLEDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/LEDController/UpdateLEDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LEDControllerServer).UpdateLEDs(ctx, req.(*UpdateLEDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LEDController_ServiceDesc is the grpc.ServiceDesc for LEDController service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LEDController_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "LEDController",
	HandlerType: (*LEDControllerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLEDs",
			Handler:    _LEDController_GetLEDs_Handler,
		},
		{
			MethodName: "UpdateLEDs",
			Handler:    _LEDController_UpdateLEDs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/led.proto",
}
