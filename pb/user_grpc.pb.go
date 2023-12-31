// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.5
// source: user.proto

package pb

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

const (
	UserService_CreateUser_FullMethodName   = "/UserService/CreateUser"
	UserService_LoginUser_FullMethodName    = "/UserService/LoginUser"
	UserService_RefreshToken_FullMethodName = "/UserService/RefreshToken"
	UserService_VerifyEmail_FullMethodName  = "/UserService/VerifyEmail"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	CreateUser(ctx context.Context, in *CreateUserReq, opts ...grpc.CallOption) (*CreateUserRes, error)
	LoginUser(ctx context.Context, in *LoginUserReq, opts ...grpc.CallOption) (*LoginUserRes, error)
	RefreshToken(ctx context.Context, in *RefreshTokenReq, opts ...grpc.CallOption) (*RefreshTokenRes, error)
	VerifyEmail(ctx context.Context, in *VerifyEmailReq, opts ...grpc.CallOption) (*VerifyEmailRes, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserReq, opts ...grpc.CallOption) (*CreateUserRes, error) {
	out := new(CreateUserRes)
	err := c.cc.Invoke(ctx, UserService_CreateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) LoginUser(ctx context.Context, in *LoginUserReq, opts ...grpc.CallOption) (*LoginUserRes, error) {
	out := new(LoginUserRes)
	err := c.cc.Invoke(ctx, UserService_LoginUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RefreshToken(ctx context.Context, in *RefreshTokenReq, opts ...grpc.CallOption) (*RefreshTokenRes, error) {
	out := new(RefreshTokenRes)
	err := c.cc.Invoke(ctx, UserService_RefreshToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) VerifyEmail(ctx context.Context, in *VerifyEmailReq, opts ...grpc.CallOption) (*VerifyEmailRes, error) {
	out := new(VerifyEmailRes)
	err := c.cc.Invoke(ctx, UserService_VerifyEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	CreateUser(context.Context, *CreateUserReq) (*CreateUserRes, error)
	LoginUser(context.Context, *LoginUserReq) (*LoginUserRes, error)
	RefreshToken(context.Context, *RefreshTokenReq) (*RefreshTokenRes, error)
	VerifyEmail(context.Context, *VerifyEmailReq) (*VerifyEmailRes, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) CreateUser(context.Context, *CreateUserReq) (*CreateUserRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedUserServiceServer) LoginUser(context.Context, *LoginUserReq) (*LoginUserRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedUserServiceServer) RefreshToken(context.Context, *RefreshTokenReq) (*RefreshTokenRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshToken not implemented")
}
func (UnimplementedUserServiceServer) VerifyEmail(context.Context, *VerifyEmailReq) (*VerifyEmailRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyEmail not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateUser(ctx, req.(*CreateUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_LoginUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).LoginUser(ctx, req.(*LoginUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RefreshToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_RefreshToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RefreshToken(ctx, req.(*RefreshTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_VerifyEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyEmailReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).VerifyEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_VerifyEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).VerifyEmail(ctx, req.(*VerifyEmailReq))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _UserService_CreateUser_Handler,
		},
		{
			MethodName: "LoginUser",
			Handler:    _UserService_LoginUser_Handler,
		},
		{
			MethodName: "RefreshToken",
			Handler:    _UserService_RefreshToken_Handler,
		},
		{
			MethodName: "VerifyEmail",
			Handler:    _UserService_VerifyEmail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
