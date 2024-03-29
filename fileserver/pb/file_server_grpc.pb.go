// file_server.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: file_server.proto

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
	FileServer_GetFiles_FullMethodName     = "/pb.FileServer/GetFiles"
	FileServer_CheckFileMD5_FullMethodName = "/pb.FileServer/CheckFileMD5"
	FileServer_UploadFile_FullMethodName   = "/pb.FileServer/UploadFile"
	FileServer_DownloadFile_FullMethodName = "/pb.FileServer/DownloadFile"
)

// FileServerClient is the client API for FileServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileServerClient interface {
	GetFiles(ctx context.Context, in *Path, opts ...grpc.CallOption) (*FileList, error)
	CheckFileMD5(ctx context.Context, in *FileCheckRequest, opts ...grpc.CallOption) (*FileCheckResponse, error)
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (FileServer_UploadFileClient, error)
	DownloadFile(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (FileServer_DownloadFileClient, error)
}

type fileServerClient struct {
	cc grpc.ClientConnInterface
}

func NewFileServerClient(cc grpc.ClientConnInterface) FileServerClient {
	return &fileServerClient{cc}
}

func (c *fileServerClient) GetFiles(ctx context.Context, in *Path, opts ...grpc.CallOption) (*FileList, error) {
	out := new(FileList)
	err := c.cc.Invoke(ctx, FileServer_GetFiles_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServerClient) CheckFileMD5(ctx context.Context, in *FileCheckRequest, opts ...grpc.CallOption) (*FileCheckResponse, error) {
	out := new(FileCheckResponse)
	err := c.cc.Invoke(ctx, FileServer_CheckFileMD5_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServerClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (FileServer_UploadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileServer_ServiceDesc.Streams[0], FileServer_UploadFile_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &fileServerUploadFileClient{stream}
	return x, nil
}

type FileServer_UploadFileClient interface {
	Send(*FileData) error
	CloseAndRecv() (*Status, error)
	grpc.ClientStream
}

type fileServerUploadFileClient struct {
	grpc.ClientStream
}

func (x *fileServerUploadFileClient) Send(m *FileData) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileServerUploadFileClient) CloseAndRecv() (*Status, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Status)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileServerClient) DownloadFile(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (FileServer_DownloadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileServer_ServiceDesc.Streams[1], FileServer_DownloadFile_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &fileServerDownloadFileClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileServer_DownloadFileClient interface {
	Recv() (*FileChunk, error)
	grpc.ClientStream
}

type fileServerDownloadFileClient struct {
	grpc.ClientStream
}

func (x *fileServerDownloadFileClient) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FileServerServer is the server API for FileServer service.
// All implementations must embed UnimplementedFileServerServer
// for forward compatibility
type FileServerServer interface {
	GetFiles(context.Context, *Path) (*FileList, error)
	CheckFileMD5(context.Context, *FileCheckRequest) (*FileCheckResponse, error)
	UploadFile(FileServer_UploadFileServer) error
	DownloadFile(*DownloadRequest, FileServer_DownloadFileServer) error
	mustEmbedUnimplementedFileServerServer()
}

// UnimplementedFileServerServer must be embedded to have forward compatible implementations.
type UnimplementedFileServerServer struct {
}

func (UnimplementedFileServerServer) GetFiles(context.Context, *Path) (*FileList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFiles not implemented")
}
func (UnimplementedFileServerServer) CheckFileMD5(context.Context, *FileCheckRequest) (*FileCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckFileMD5 not implemented")
}
func (UnimplementedFileServerServer) UploadFile(FileServer_UploadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedFileServerServer) DownloadFile(*DownloadRequest, FileServer_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedFileServerServer) mustEmbedUnimplementedFileServerServer() {}

// UnsafeFileServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileServerServer will
// result in compilation errors.
type UnsafeFileServerServer interface {
	mustEmbedUnimplementedFileServerServer()
}

func RegisterFileServerServer(s grpc.ServiceRegistrar, srv FileServerServer) {
	s.RegisterService(&FileServer_ServiceDesc, srv)
}

func _FileServer_GetFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Path)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServerServer).GetFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileServer_GetFiles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServerServer).GetFiles(ctx, req.(*Path))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileServer_CheckFileMD5_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServerServer).CheckFileMD5(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileServer_CheckFileMD5_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServerServer).CheckFileMD5(ctx, req.(*FileCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileServer_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileServerServer).UploadFile(&fileServerUploadFileServer{stream})
}

type FileServer_UploadFileServer interface {
	SendAndClose(*Status) error
	Recv() (*FileData, error)
	grpc.ServerStream
}

type fileServerUploadFileServer struct {
	grpc.ServerStream
}

func (x *fileServerUploadFileServer) SendAndClose(m *Status) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileServerUploadFileServer) Recv() (*FileData, error) {
	m := new(FileData)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FileServer_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileServerServer).DownloadFile(m, &fileServerDownloadFileServer{stream})
}

type FileServer_DownloadFileServer interface {
	Send(*FileChunk) error
	grpc.ServerStream
}

type fileServerDownloadFileServer struct {
	grpc.ServerStream
}

func (x *fileServerDownloadFileServer) Send(m *FileChunk) error {
	return x.ServerStream.SendMsg(m)
}

// FileServer_ServiceDesc is the grpc.ServiceDesc for FileServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.FileServer",
	HandlerType: (*FileServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFiles",
			Handler:    _FileServer_GetFiles_Handler,
		},
		{
			MethodName: "CheckFileMD5",
			Handler:    _FileServer_CheckFileMD5_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFile",
			Handler:       _FileServer_UploadFile_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadFile",
			Handler:       _FileServer_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "file_server.proto",
}
