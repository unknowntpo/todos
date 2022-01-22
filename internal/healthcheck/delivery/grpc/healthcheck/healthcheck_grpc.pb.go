// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package healthcheck

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

// HealthcheckClient is the client API for Healthcheck service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HealthcheckClient interface {
	Healthcheck(ctx context.Context, in *HealthcheckRequest, opts ...grpc.CallOption) (*HealthcheckResponse, error)
}

type healthcheckClient struct {
	cc grpc.ClientConnInterface
}

func NewHealthcheckClient(cc grpc.ClientConnInterface) HealthcheckClient {
	return &healthcheckClient{cc}
}

func (c *healthcheckClient) Healthcheck(ctx context.Context, in *HealthcheckRequest, opts ...grpc.CallOption) (*HealthcheckResponse, error) {
	out := new(HealthcheckResponse)
	err := c.cc.Invoke(ctx, "/healthcheck.Healthcheck/Healthcheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HealthcheckServer is the server API for Healthcheck service.
// All implementations must embed UnimplementedHealthcheckServer
// for forward compatibility
type HealthcheckServer interface {
	Healthcheck(context.Context, *HealthcheckRequest) (*HealthcheckResponse, error)
	mustEmbedUnimplementedHealthcheckServer()
}

// UnimplementedHealthcheckServer must be embedded to have forward compatible implementations.
type UnimplementedHealthcheckServer struct {
}

func (UnimplementedHealthcheckServer) Healthcheck(context.Context, *HealthcheckRequest) (*HealthcheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Healthcheck not implemented")
}
func (UnimplementedHealthcheckServer) mustEmbedUnimplementedHealthcheckServer() {}

// UnsafeHealthcheckServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HealthcheckServer will
// result in compilation errors.
type UnsafeHealthcheckServer interface {
	mustEmbedUnimplementedHealthcheckServer()
}

func RegisterHealthcheckServer(s grpc.ServiceRegistrar, srv HealthcheckServer) {
	s.RegisterService(&Healthcheck_ServiceDesc, srv)
}

func _Healthcheck_Healthcheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthcheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthcheckServer).Healthcheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/healthcheck.Healthcheck/Healthcheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthcheckServer).Healthcheck(ctx, req.(*HealthcheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Healthcheck_ServiceDesc is the grpc.ServiceDesc for Healthcheck service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Healthcheck_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "healthcheck.Healthcheck",
	HandlerType: (*HealthcheckServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Healthcheck",
			Handler:    _Healthcheck_Healthcheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/healthcheck/delivery/grpc/healthcheck/healthcheck.proto",
}
