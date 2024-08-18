// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: proxy/proxy.proto

package proxy

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	ProxyService_CreateProxy_FullMethodName              = "/proxy.ProxyService/CreateProxy"
	ProxyService_GetProxy_FullMethodName                 = "/proxy.ProxyService/GetProxy"
	ProxyService_UpdateProxy_FullMethodName              = "/proxy.ProxyService/UpdateProxy"
	ProxyService_DeleteProxy_FullMethodName              = "/proxy.ProxyService/DeleteProxy"
	ProxyService_CreateProxyBlock_FullMethodName         = "/proxy.ProxyService/CreateProxyBlock"
	ProxyService_GetProxyBlock_FullMethodName            = "/proxy.ProxyService/GetProxyBlock"
	ProxyService_UpdateProxyBlock_FullMethodName         = "/proxy.ProxyService/UpdateProxyBlock"
	ProxyService_DeleteProxyBlock_FullMethodName         = "/proxy.ProxyService/DeleteProxyBlock"
	ProxyService_GetProxiesByProxyBlockId_FullMethodName = "/proxy.ProxyService/GetProxiesByProxyBlockId"
)

// ProxyServiceClient is the client API for ProxyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProxyServiceClient interface {
	CreateProxy(ctx context.Context, in *CreateProxyRequest, opts ...grpc.CallOption) (*CreateProxyResponse, error)
	GetProxy(ctx context.Context, in *GetProxyRequest, opts ...grpc.CallOption) (*GetProxyResponse, error)
	UpdateProxy(ctx context.Context, in *UpdateProxyRequest, opts ...grpc.CallOption) (*UpdateProxyResponse, error)
	DeleteProxy(ctx context.Context, in *DeleteProxyRequest, opts ...grpc.CallOption) (*DeleteProxyResponse, error)
	CreateProxyBlock(ctx context.Context, in *CreateProxyBlockRequest, opts ...grpc.CallOption) (*CreateProxyBlockResponse, error)
	GetProxyBlock(ctx context.Context, in *GetProxyBlockRequest, opts ...grpc.CallOption) (*GetProxyBlockResponse, error)
	UpdateProxyBlock(ctx context.Context, in *UpdateProxyBlockRequest, opts ...grpc.CallOption) (*UpdateProxyBlockResponse, error)
	DeleteProxyBlock(ctx context.Context, in *DeleteProxyBlockRequest, opts ...grpc.CallOption) (*DeleteProxyBlockResponse, error)
	GetProxiesByProxyBlockId(ctx context.Context, in *GetProxiesByProxyBlockIdRequest, opts ...grpc.CallOption) (*GetProxiesByProxyBlockIdResponse, error)
}

type proxyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProxyServiceClient(cc grpc.ClientConnInterface) ProxyServiceClient {
	return &proxyServiceClient{cc}
}

func (c *proxyServiceClient) CreateProxy(ctx context.Context, in *CreateProxyRequest, opts ...grpc.CallOption) (*CreateProxyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateProxyResponse)
	err := c.cc.Invoke(ctx, ProxyService_CreateProxy_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyServiceClient) GetProxy(ctx context.Context, in *GetProxyRequest, opts ...grpc.CallOption) (*GetProxyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProxyResponse)
	err := c.cc.Invoke(ctx, ProxyService_GetProxy_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyServiceClient) UpdateProxy(ctx context.Context, in *UpdateProxyRequest, opts ...grpc.CallOption) (*UpdateProxyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateProxyResponse)
	err := c.cc.Invoke(ctx, ProxyService_UpdateProxy_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyServiceClient) DeleteProxy(ctx context.Context, in *DeleteProxyRequest, opts ...grpc.CallOption) (*DeleteProxyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteProxyResponse)
	err := c.cc.Invoke(ctx, ProxyService_DeleteProxy_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyServiceClient) CreateProxyBlock(ctx context.Context, in *CreateProxyBlockRequest, opts ...grpc.CallOption) (*CreateProxyBlockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateProxyBlockResponse)
	err := c.cc.Invoke(ctx, ProxyService_CreateProxyBlock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyServiceClient) GetProxyBlock(ctx context.Context, in *GetProxyBlockRequest, opts ...grpc.CallOption) (*GetProxyBlockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProxyBlockResponse)
	err := c.cc.Invoke(ctx, ProxyService_GetProxyBlock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyServiceClient) UpdateProxyBlock(ctx context.Context, in *UpdateProxyBlockRequest, opts ...grpc.CallOption) (*UpdateProxyBlockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateProxyBlockResponse)
	err := c.cc.Invoke(ctx, ProxyService_UpdateProxyBlock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyServiceClient) DeleteProxyBlock(ctx context.Context, in *DeleteProxyBlockRequest, opts ...grpc.CallOption) (*DeleteProxyBlockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteProxyBlockResponse)
	err := c.cc.Invoke(ctx, ProxyService_DeleteProxyBlock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyServiceClient) GetProxiesByProxyBlockId(ctx context.Context, in *GetProxiesByProxyBlockIdRequest, opts ...grpc.CallOption) (*GetProxiesByProxyBlockIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProxiesByProxyBlockIdResponse)
	err := c.cc.Invoke(ctx, ProxyService_GetProxiesByProxyBlockId_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProxyServiceServer is the server API for ProxyService service.
// All implementations must embed UnimplementedProxyServiceServer
// for forward compatibility
type ProxyServiceServer interface {
	CreateProxy(context.Context, *CreateProxyRequest) (*CreateProxyResponse, error)
	GetProxy(context.Context, *GetProxyRequest) (*GetProxyResponse, error)
	UpdateProxy(context.Context, *UpdateProxyRequest) (*UpdateProxyResponse, error)
	DeleteProxy(context.Context, *DeleteProxyRequest) (*DeleteProxyResponse, error)
	CreateProxyBlock(context.Context, *CreateProxyBlockRequest) (*CreateProxyBlockResponse, error)
	GetProxyBlock(context.Context, *GetProxyBlockRequest) (*GetProxyBlockResponse, error)
	UpdateProxyBlock(context.Context, *UpdateProxyBlockRequest) (*UpdateProxyBlockResponse, error)
	DeleteProxyBlock(context.Context, *DeleteProxyBlockRequest) (*DeleteProxyBlockResponse, error)
	GetProxiesByProxyBlockId(context.Context, *GetProxiesByProxyBlockIdRequest) (*GetProxiesByProxyBlockIdResponse, error)
	mustEmbedUnimplementedProxyServiceServer()
}

// UnimplementedProxyServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProxyServiceServer struct {
}

func (UnimplementedProxyServiceServer) CreateProxy(context.Context, *CreateProxyRequest) (*CreateProxyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProxy not implemented")
}
func (UnimplementedProxyServiceServer) GetProxy(context.Context, *GetProxyRequest) (*GetProxyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProxy not implemented")
}
func (UnimplementedProxyServiceServer) UpdateProxy(context.Context, *UpdateProxyRequest) (*UpdateProxyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProxy not implemented")
}
func (UnimplementedProxyServiceServer) DeleteProxy(context.Context, *DeleteProxyRequest) (*DeleteProxyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProxy not implemented")
}
func (UnimplementedProxyServiceServer) CreateProxyBlock(context.Context, *CreateProxyBlockRequest) (*CreateProxyBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProxyBlock not implemented")
}
func (UnimplementedProxyServiceServer) GetProxyBlock(context.Context, *GetProxyBlockRequest) (*GetProxyBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProxyBlock not implemented")
}
func (UnimplementedProxyServiceServer) UpdateProxyBlock(context.Context, *UpdateProxyBlockRequest) (*UpdateProxyBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProxyBlock not implemented")
}
func (UnimplementedProxyServiceServer) DeleteProxyBlock(context.Context, *DeleteProxyBlockRequest) (*DeleteProxyBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProxyBlock not implemented")
}
func (UnimplementedProxyServiceServer) GetProxiesByProxyBlockId(context.Context, *GetProxiesByProxyBlockIdRequest) (*GetProxiesByProxyBlockIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProxiesByProxyBlockId not implemented")
}
func (UnimplementedProxyServiceServer) mustEmbedUnimplementedProxyServiceServer() {}

// UnsafeProxyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProxyServiceServer will
// result in compilation errors.
type UnsafeProxyServiceServer interface {
	mustEmbedUnimplementedProxyServiceServer()
}

func RegisterProxyServiceServer(s grpc.ServiceRegistrar, srv ProxyServiceServer) {
	s.RegisterService(&ProxyService_ServiceDesc, srv)
}

func _ProxyService_CreateProxy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProxyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServiceServer).CreateProxy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyService_CreateProxy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServiceServer).CreateProxy(ctx, req.(*CreateProxyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyService_GetProxy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProxyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServiceServer).GetProxy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyService_GetProxy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServiceServer).GetProxy(ctx, req.(*GetProxyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyService_UpdateProxy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProxyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServiceServer).UpdateProxy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyService_UpdateProxy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServiceServer).UpdateProxy(ctx, req.(*UpdateProxyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyService_DeleteProxy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProxyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServiceServer).DeleteProxy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyService_DeleteProxy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServiceServer).DeleteProxy(ctx, req.(*DeleteProxyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyService_CreateProxyBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProxyBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServiceServer).CreateProxyBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyService_CreateProxyBlock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServiceServer).CreateProxyBlock(ctx, req.(*CreateProxyBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyService_GetProxyBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProxyBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServiceServer).GetProxyBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyService_GetProxyBlock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServiceServer).GetProxyBlock(ctx, req.(*GetProxyBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyService_UpdateProxyBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProxyBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServiceServer).UpdateProxyBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyService_UpdateProxyBlock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServiceServer).UpdateProxyBlock(ctx, req.(*UpdateProxyBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyService_DeleteProxyBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProxyBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServiceServer).DeleteProxyBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyService_DeleteProxyBlock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServiceServer).DeleteProxyBlock(ctx, req.(*DeleteProxyBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProxyService_GetProxiesByProxyBlockId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProxiesByProxyBlockIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServiceServer).GetProxiesByProxyBlockId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProxyService_GetProxiesByProxyBlockId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServiceServer).GetProxiesByProxyBlockId(ctx, req.(*GetProxiesByProxyBlockIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProxyService_ServiceDesc is the grpc.ServiceDesc for ProxyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProxyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proxy.ProxyService",
	HandlerType: (*ProxyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateProxy",
			Handler:    _ProxyService_CreateProxy_Handler,
		},
		{
			MethodName: "GetProxy",
			Handler:    _ProxyService_GetProxy_Handler,
		},
		{
			MethodName: "UpdateProxy",
			Handler:    _ProxyService_UpdateProxy_Handler,
		},
		{
			MethodName: "DeleteProxy",
			Handler:    _ProxyService_DeleteProxy_Handler,
		},
		{
			MethodName: "CreateProxyBlock",
			Handler:    _ProxyService_CreateProxyBlock_Handler,
		},
		{
			MethodName: "GetProxyBlock",
			Handler:    _ProxyService_GetProxyBlock_Handler,
		},
		{
			MethodName: "UpdateProxyBlock",
			Handler:    _ProxyService_UpdateProxyBlock_Handler,
		},
		{
			MethodName: "DeleteProxyBlock",
			Handler:    _ProxyService_DeleteProxyBlock_Handler,
		},
		{
			MethodName: "GetProxiesByProxyBlockId",
			Handler:    _ProxyService_GetProxiesByProxyBlockId_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proxy/proxy.proto",
}
