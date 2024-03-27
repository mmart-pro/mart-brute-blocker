// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: mbb/mbb.proto

package pb

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MBBService_Allow_FullMethodName       = "/api.MBBService/Allow"
	MBBService_Deny_FullMethodName        = "/api.MBBService/Deny"
	MBBService_Remove_FullMethodName      = "/api.MBBService/Remove"
	MBBService_Exists_FullMethodName      = "/api.MBBService/Exists"
	MBBService_Contains_FullMethodName    = "/api.MBBService/Contains"
	MBBService_ClearBucket_FullMethodName = "/api.MBBService/ClearBucket"
	MBBService_Check_FullMethodName       = "/api.MBBService/Check"
)

// MBBServiceClient is the client API for MBBService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MBBServiceClient interface {
	Allow(ctx context.Context, in *SubnetReq, opts ...grpc.CallOption) (*empty.Empty, error)
	Deny(ctx context.Context, in *SubnetReq, opts ...grpc.CallOption) (*empty.Empty, error)
	Remove(ctx context.Context, in *SubnetReq, opts ...grpc.CallOption) (*empty.Empty, error)
	Exists(ctx context.Context, in *SubnetReq, opts ...grpc.CallOption) (*ExistsResponse, error)
	Contains(ctx context.Context, in *IpReq, opts ...grpc.CallOption) (*ContainsResponse, error)
	ClearBucket(ctx context.Context, in *ClearBucketRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
}

type mBBServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMBBServiceClient(cc grpc.ClientConnInterface) MBBServiceClient {
	return &mBBServiceClient{cc}
}

func (c *mBBServiceClient) Allow(ctx context.Context, in *SubnetReq, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, MBBService_Allow_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mBBServiceClient) Deny(ctx context.Context, in *SubnetReq, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, MBBService_Deny_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mBBServiceClient) Remove(ctx context.Context, in *SubnetReq, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, MBBService_Remove_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mBBServiceClient) Exists(ctx context.Context, in *SubnetReq, opts ...grpc.CallOption) (*ExistsResponse, error) {
	out := new(ExistsResponse)
	err := c.cc.Invoke(ctx, MBBService_Exists_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mBBServiceClient) Contains(ctx context.Context, in *IpReq, opts ...grpc.CallOption) (*ContainsResponse, error) {
	out := new(ContainsResponse)
	err := c.cc.Invoke(ctx, MBBService_Contains_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mBBServiceClient) ClearBucket(ctx context.Context, in *ClearBucketRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, MBBService_ClearBucket_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mBBServiceClient) Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, MBBService_Check_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MBBServiceServer is the server API for MBBService service.
// All implementations must embed UnimplementedMBBServiceServer
// for forward compatibility
type MBBServiceServer interface {
	Allow(context.Context, *SubnetReq) (*empty.Empty, error)
	Deny(context.Context, *SubnetReq) (*empty.Empty, error)
	Remove(context.Context, *SubnetReq) (*empty.Empty, error)
	Exists(context.Context, *SubnetReq) (*ExistsResponse, error)
	Contains(context.Context, *IpReq) (*ContainsResponse, error)
	ClearBucket(context.Context, *ClearBucketRequest) (*empty.Empty, error)
	Check(context.Context, *CheckRequest) (*CheckResponse, error)
	mustEmbedUnimplementedMBBServiceServer()
}

// UnimplementedMBBServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMBBServiceServer struct {
}

func (UnimplementedMBBServiceServer) Allow(context.Context, *SubnetReq) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Allow not implemented")
}
func (UnimplementedMBBServiceServer) Deny(context.Context, *SubnetReq) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deny not implemented")
}
func (UnimplementedMBBServiceServer) Remove(context.Context, *SubnetReq) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (UnimplementedMBBServiceServer) Exists(context.Context, *SubnetReq) (*ExistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (UnimplementedMBBServiceServer) Contains(context.Context, *IpReq) (*ContainsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Contains not implemented")
}
func (UnimplementedMBBServiceServer) ClearBucket(context.Context, *ClearBucketRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearBucket not implemented")
}
func (UnimplementedMBBServiceServer) Check(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedMBBServiceServer) mustEmbedUnimplementedMBBServiceServer() {}

// UnsafeMBBServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MBBServiceServer will
// result in compilation errors.
type UnsafeMBBServiceServer interface {
	mustEmbedUnimplementedMBBServiceServer()
}

func RegisterMBBServiceServer(s grpc.ServiceRegistrar, srv MBBServiceServer) {
	s.RegisterService(&MBBService_ServiceDesc, srv)
}

func _MBBService_Allow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubnetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MBBServiceServer).Allow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MBBService_Allow_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MBBServiceServer).Allow(ctx, req.(*SubnetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MBBService_Deny_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubnetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MBBServiceServer).Deny(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MBBService_Deny_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MBBServiceServer).Deny(ctx, req.(*SubnetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MBBService_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubnetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MBBServiceServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MBBService_Remove_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MBBServiceServer).Remove(ctx, req.(*SubnetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MBBService_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubnetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MBBServiceServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MBBService_Exists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MBBServiceServer).Exists(ctx, req.(*SubnetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MBBService_Contains_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IpReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MBBServiceServer).Contains(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MBBService_Contains_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MBBServiceServer).Contains(ctx, req.(*IpReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MBBService_ClearBucket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearBucketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MBBServiceServer).ClearBucket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MBBService_ClearBucket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MBBServiceServer).ClearBucket(ctx, req.(*ClearBucketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MBBService_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MBBServiceServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MBBService_Check_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MBBServiceServer).Check(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MBBService_ServiceDesc is the grpc.ServiceDesc for MBBService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MBBService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.MBBService",
	HandlerType: (*MBBServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Allow",
			Handler:    _MBBService_Allow_Handler,
		},
		{
			MethodName: "Deny",
			Handler:    _MBBService_Deny_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _MBBService_Remove_Handler,
		},
		{
			MethodName: "Exists",
			Handler:    _MBBService_Exists_Handler,
		},
		{
			MethodName: "Contains",
			Handler:    _MBBService_Contains_Handler,
		},
		{
			MethodName: "ClearBucket",
			Handler:    _MBBService_ClearBucket_Handler,
		},
		{
			MethodName: "Check",
			Handler:    _MBBService_Check_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mbb/mbb.proto",
}
