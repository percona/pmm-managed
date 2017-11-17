// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rds.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type RDSNode struct {
	Region string `protobuf:"bytes,3,opt,name=region" json:"region,omitempty"`
	Name   string `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
}

func (m *RDSNode) Reset()                    { *m = RDSNode{} }
func (m *RDSNode) String() string            { return proto.CompactTextString(m) }
func (*RDSNode) ProtoMessage()               {}
func (*RDSNode) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *RDSNode) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *RDSNode) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RDSService struct {
	Address       string `protobuf:"bytes,4,opt,name=address" json:"address,omitempty"`
	Port          uint32 `protobuf:"varint,5,opt,name=port" json:"port,omitempty"`
	Engine        string `protobuf:"bytes,6,opt,name=engine" json:"engine,omitempty"`
	EngineVersion string `protobuf:"bytes,7,opt,name=engine_version,json=engineVersion" json:"engine_version,omitempty"`
}

func (m *RDSService) Reset()                    { *m = RDSService{} }
func (m *RDSService) String() string            { return proto.CompactTextString(m) }
func (*RDSService) ProtoMessage()               {}
func (*RDSService) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *RDSService) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *RDSService) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *RDSService) GetEngine() string {
	if m != nil {
		return m.Engine
	}
	return ""
}

func (m *RDSService) GetEngineVersion() string {
	if m != nil {
		return m.EngineVersion
	}
	return ""
}

type RDSInstanceID struct {
	Region string `protobuf:"bytes,1,opt,name=region" json:"region,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}

func (m *RDSInstanceID) Reset()                    { *m = RDSInstanceID{} }
func (m *RDSInstanceID) String() string            { return proto.CompactTextString(m) }
func (*RDSInstanceID) ProtoMessage()               {}
func (*RDSInstanceID) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

func (m *RDSInstanceID) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *RDSInstanceID) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RDSInstance struct {
	Node    *RDSNode    `protobuf:"bytes,1,opt,name=node" json:"node,omitempty"`
	Service *RDSService `protobuf:"bytes,2,opt,name=service" json:"service,omitempty"`
}

func (m *RDSInstance) Reset()                    { *m = RDSInstance{} }
func (m *RDSInstance) String() string            { return proto.CompactTextString(m) }
func (*RDSInstance) ProtoMessage()               {}
func (*RDSInstance) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

func (m *RDSInstance) GetNode() *RDSNode {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *RDSInstance) GetService() *RDSService {
	if m != nil {
		return m.Service
	}
	return nil
}

type RDSDiscoverRequest struct {
	AwsAccessKeyId     string `protobuf:"bytes,1,opt,name=aws_access_key_id,json=awsAccessKeyId" json:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey string `protobuf:"bytes,2,opt,name=aws_secret_access_key,json=awsSecretAccessKey" json:"aws_secret_access_key,omitempty"`
}

func (m *RDSDiscoverRequest) Reset()                    { *m = RDSDiscoverRequest{} }
func (m *RDSDiscoverRequest) String() string            { return proto.CompactTextString(m) }
func (*RDSDiscoverRequest) ProtoMessage()               {}
func (*RDSDiscoverRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

func (m *RDSDiscoverRequest) GetAwsAccessKeyId() string {
	if m != nil {
		return m.AwsAccessKeyId
	}
	return ""
}

func (m *RDSDiscoverRequest) GetAwsSecretAccessKey() string {
	if m != nil {
		return m.AwsSecretAccessKey
	}
	return ""
}

type RDSDiscoverResponse struct {
	Instances []*RDSInstance `protobuf:"bytes,1,rep,name=instances" json:"instances,omitempty"`
}

func (m *RDSDiscoverResponse) Reset()                    { *m = RDSDiscoverResponse{} }
func (m *RDSDiscoverResponse) String() string            { return proto.CompactTextString(m) }
func (*RDSDiscoverResponse) ProtoMessage()               {}
func (*RDSDiscoverResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{5} }

func (m *RDSDiscoverResponse) GetInstances() []*RDSInstance {
	if m != nil {
		return m.Instances
	}
	return nil
}

type RDSListRequest struct {
}

func (m *RDSListRequest) Reset()                    { *m = RDSListRequest{} }
func (m *RDSListRequest) String() string            { return proto.CompactTextString(m) }
func (*RDSListRequest) ProtoMessage()               {}
func (*RDSListRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{6} }

type RDSListResponse struct {
	Instances []*RDSInstance `protobuf:"bytes,1,rep,name=instances" json:"instances,omitempty"`
}

func (m *RDSListResponse) Reset()                    { *m = RDSListResponse{} }
func (m *RDSListResponse) String() string            { return proto.CompactTextString(m) }
func (*RDSListResponse) ProtoMessage()               {}
func (*RDSListResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{7} }

func (m *RDSListResponse) GetInstances() []*RDSInstance {
	if m != nil {
		return m.Instances
	}
	return nil
}

type RDSAddRequest struct {
	AwsAccessKeyId     string         `protobuf:"bytes,1,opt,name=aws_access_key_id,json=awsAccessKeyId" json:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey string         `protobuf:"bytes,2,opt,name=aws_secret_access_key,json=awsSecretAccessKey" json:"aws_secret_access_key,omitempty"`
	Id                 *RDSInstanceID `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
	Username           string         `protobuf:"bytes,4,opt,name=username" json:"username,omitempty"`
	Password           string         `protobuf:"bytes,5,opt,name=password" json:"password,omitempty"`
}

func (m *RDSAddRequest) Reset()                    { *m = RDSAddRequest{} }
func (m *RDSAddRequest) String() string            { return proto.CompactTextString(m) }
func (*RDSAddRequest) ProtoMessage()               {}
func (*RDSAddRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{8} }

func (m *RDSAddRequest) GetAwsAccessKeyId() string {
	if m != nil {
		return m.AwsAccessKeyId
	}
	return ""
}

func (m *RDSAddRequest) GetAwsSecretAccessKey() string {
	if m != nil {
		return m.AwsSecretAccessKey
	}
	return ""
}

func (m *RDSAddRequest) GetId() *RDSInstanceID {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *RDSAddRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *RDSAddRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type RDSAddResponse struct {
}

func (m *RDSAddResponse) Reset()                    { *m = RDSAddResponse{} }
func (m *RDSAddResponse) String() string            { return proto.CompactTextString(m) }
func (*RDSAddResponse) ProtoMessage()               {}
func (*RDSAddResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{9} }

type RDSRemoveRequest struct {
	Id *RDSInstanceID `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *RDSRemoveRequest) Reset()                    { *m = RDSRemoveRequest{} }
func (m *RDSRemoveRequest) String() string            { return proto.CompactTextString(m) }
func (*RDSRemoveRequest) ProtoMessage()               {}
func (*RDSRemoveRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{10} }

func (m *RDSRemoveRequest) GetId() *RDSInstanceID {
	if m != nil {
		return m.Id
	}
	return nil
}

type RDSRemoveResponse struct {
}

func (m *RDSRemoveResponse) Reset()                    { *m = RDSRemoveResponse{} }
func (m *RDSRemoveResponse) String() string            { return proto.CompactTextString(m) }
func (*RDSRemoveResponse) ProtoMessage()               {}
func (*RDSRemoveResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{11} }

func init() {
	proto.RegisterType((*RDSNode)(nil), "api.RDSNode")
	proto.RegisterType((*RDSService)(nil), "api.RDSService")
	proto.RegisterType((*RDSInstanceID)(nil), "api.RDSInstanceID")
	proto.RegisterType((*RDSInstance)(nil), "api.RDSInstance")
	proto.RegisterType((*RDSDiscoverRequest)(nil), "api.RDSDiscoverRequest")
	proto.RegisterType((*RDSDiscoverResponse)(nil), "api.RDSDiscoverResponse")
	proto.RegisterType((*RDSListRequest)(nil), "api.RDSListRequest")
	proto.RegisterType((*RDSListResponse)(nil), "api.RDSListResponse")
	proto.RegisterType((*RDSAddRequest)(nil), "api.RDSAddRequest")
	proto.RegisterType((*RDSAddResponse)(nil), "api.RDSAddResponse")
	proto.RegisterType((*RDSRemoveRequest)(nil), "api.RDSRemoveRequest")
	proto.RegisterType((*RDSRemoveResponse)(nil), "api.RDSRemoveResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for RDS service

type RDSClient interface {
	Discover(ctx context.Context, in *RDSDiscoverRequest, opts ...grpc.CallOption) (*RDSDiscoverResponse, error)
	List(ctx context.Context, in *RDSListRequest, opts ...grpc.CallOption) (*RDSListResponse, error)
	Add(ctx context.Context, in *RDSAddRequest, opts ...grpc.CallOption) (*RDSAddResponse, error)
	Remove(ctx context.Context, in *RDSRemoveRequest, opts ...grpc.CallOption) (*RDSRemoveResponse, error)
}

type rDSClient struct {
	cc *grpc.ClientConn
}

func NewRDSClient(cc *grpc.ClientConn) RDSClient {
	return &rDSClient{cc}
}

func (c *rDSClient) Discover(ctx context.Context, in *RDSDiscoverRequest, opts ...grpc.CallOption) (*RDSDiscoverResponse, error) {
	out := new(RDSDiscoverResponse)
	err := grpc.Invoke(ctx, "/api.RDS/Discover", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDSClient) List(ctx context.Context, in *RDSListRequest, opts ...grpc.CallOption) (*RDSListResponse, error) {
	out := new(RDSListResponse)
	err := grpc.Invoke(ctx, "/api.RDS/List", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDSClient) Add(ctx context.Context, in *RDSAddRequest, opts ...grpc.CallOption) (*RDSAddResponse, error) {
	out := new(RDSAddResponse)
	err := grpc.Invoke(ctx, "/api.RDS/Add", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDSClient) Remove(ctx context.Context, in *RDSRemoveRequest, opts ...grpc.CallOption) (*RDSRemoveResponse, error) {
	out := new(RDSRemoveResponse)
	err := grpc.Invoke(ctx, "/api.RDS/Remove", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RDS service

type RDSServer interface {
	Discover(context.Context, *RDSDiscoverRequest) (*RDSDiscoverResponse, error)
	List(context.Context, *RDSListRequest) (*RDSListResponse, error)
	Add(context.Context, *RDSAddRequest) (*RDSAddResponse, error)
	Remove(context.Context, *RDSRemoveRequest) (*RDSRemoveResponse, error)
}

func RegisterRDSServer(s *grpc.Server, srv RDSServer) {
	s.RegisterService(&_RDS_serviceDesc, srv)
}

func _RDS_Discover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RDSDiscoverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDSServer).Discover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RDS/Discover",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDSServer).Discover(ctx, req.(*RDSDiscoverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RDS_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RDSListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDSServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RDS/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDSServer).List(ctx, req.(*RDSListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RDS_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RDSAddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDSServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RDS/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDSServer).Add(ctx, req.(*RDSAddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RDS_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RDSRemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDSServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.RDS/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDSServer).Remove(ctx, req.(*RDSRemoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RDS_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.RDS",
	HandlerType: (*RDSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Discover",
			Handler:    _RDS_Discover_Handler,
		},
		{
			MethodName: "List",
			Handler:    _RDS_List_Handler,
		},
		{
			MethodName: "Add",
			Handler:    _RDS_Add_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _RDS_Remove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rds.proto",
}

func init() { proto.RegisterFile("rds.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 590 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x54, 0xbb, 0x6e, 0x13, 0x41,
	0x14, 0xd5, 0x3e, 0xf0, 0xe3, 0x3a, 0x8f, 0xcd, 0x35, 0x09, 0x2b, 0x8b, 0xc2, 0x5a, 0x09, 0x29,
	0x49, 0x61, 0x83, 0x91, 0x28, 0xa0, 0x72, 0xb4, 0x14, 0x4e, 0x10, 0xc5, 0xac, 0x04, 0x12, 0x8d,
	0x35, 0x78, 0x46, 0xd6, 0x08, 0xb2, 0xb3, 0xcc, 0x6c, 0x6c, 0xa5, 0xa5, 0xa5, 0xe4, 0x53, 0xf8,
	0x0a, 0x6a, 0x7e, 0x81, 0x0f, 0x41, 0x3b, 0x3b, 0xe3, 0xac, 0x0d, 0xa2, 0xa0, 0xa0, 0xdb, 0x7b,
	0xee, 0x9c, 0x33, 0xe7, 0xdc, 0x3b, 0x5a, 0xe8, 0x2a, 0xa6, 0x47, 0x85, 0x92, 0xa5, 0xc4, 0x80,
	0x16, 0x62, 0xf0, 0x70, 0x29, 0xe5, 0xf2, 0x23, 0x1f, 0xd3, 0x42, 0x8c, 0x69, 0x9e, 0xcb, 0x92,
	0x96, 0x42, 0xe6, 0xf6, 0x48, 0x32, 0x85, 0x36, 0x49, 0xb3, 0xd7, 0x92, 0x71, 0x3c, 0x81, 0x96,
	0xe2, 0x4b, 0x21, 0xf3, 0x38, 0x18, 0x7a, 0xa7, 0x5d, 0x62, 0x2b, 0x44, 0x08, 0x73, 0x7a, 0xcd,
	0xe3, 0xd0, 0xa0, 0xe6, 0xfb, 0x32, 0xec, 0x78, 0x91, 0x7f, 0x19, 0x76, 0xfc, 0x28, 0x48, 0xbe,
	0x78, 0x00, 0x24, 0xcd, 0x32, 0xae, 0x56, 0x62, 0xc1, 0x31, 0x86, 0x36, 0x65, 0x4c, 0x71, 0xad,
	0x2d, 0xc3, 0x95, 0x95, 0x50, 0x21, 0x55, 0x19, 0xdf, 0x1b, 0x7a, 0xa7, 0xfb, 0xc4, 0x7c, 0x57,
	0x97, 0xf2, 0x7c, 0x29, 0x72, 0x1e, 0xb7, 0xea, 0x4b, 0xeb, 0x0a, 0x1f, 0xc1, 0x41, 0xfd, 0x35,
	0x5f, 0x71, 0xa5, 0x2b, 0x53, 0x6d, 0xd3, 0xdf, 0xaf, 0xd1, 0x37, 0x35, 0xd8, 0xf4, 0x71, 0x19,
	0x76, 0x82, 0x28, 0x4c, 0x5e, 0xc0, 0x3e, 0x49, 0xb3, 0x59, 0xae, 0x4b, 0x9a, 0x2f, 0xf8, 0x2c,
	0x6d, 0xc4, 0xf2, 0xfe, 0x18, 0xcb, 0xbf, 0x8b, 0x95, 0xbc, 0x83, 0x5e, 0x83, 0x8c, 0x43, 0x08,
	0x73, 0xc9, 0xb8, 0x21, 0xf6, 0x26, 0x7b, 0x23, 0x5a, 0x88, 0x91, 0x9d, 0x16, 0x31, 0x1d, 0x3c,
	0x83, 0xb6, 0xae, 0x73, 0x1b, 0x9d, 0xde, 0xe4, 0xd0, 0x1d, 0xb2, 0xe3, 0x20, 0xae, 0x9f, 0x28,
	0x40, 0x92, 0x66, 0xa9, 0xd0, 0x0b, 0xb9, 0xe2, 0x8a, 0xf0, 0x4f, 0x37, 0x5c, 0x97, 0x78, 0x06,
	0x47, 0x74, 0xad, 0xe7, 0x74, 0xb1, 0xe0, 0x5a, 0xcf, 0x3f, 0xf0, 0xdb, 0xb9, 0x60, 0xd6, 0xe8,
	0x01, 0x5d, 0xeb, 0xa9, 0xc1, 0xaf, 0xf8, 0xed, 0x8c, 0xe1, 0x13, 0x38, 0xae, 0x8e, 0x6a, 0xbe,
	0x50, 0xbc, 0x6c, 0x30, 0x6c, 0x02, 0xa4, 0x6b, 0x9d, 0x99, 0xde, 0x86, 0x94, 0xbc, 0x84, 0xfe,
	0xd6, 0x9d, 0xba, 0x90, 0xb9, 0xe6, 0x38, 0x82, 0xae, 0xb0, 0x19, 0x75, 0xec, 0x0d, 0x83, 0xd3,
	0xde, 0x24, 0x72, 0xbe, 0x5d, 0x78, 0x72, 0x77, 0x24, 0x89, 0xe0, 0x80, 0xa4, 0xd9, 0x2b, 0xa1,
	0x4b, 0x6b, 0x3b, 0x99, 0xc2, 0xe1, 0x06, 0xf9, 0x47, 0xd1, 0xef, 0x9e, 0xd9, 0xd4, 0x94, 0xb1,
	0xff, 0x32, 0x0b, 0x4c, 0xc0, 0x17, 0xcc, 0x3c, 0xed, 0xde, 0x04, 0x77, 0x8d, 0xcd, 0x52, 0xe2,
	0x0b, 0x86, 0x03, 0xe8, 0xdc, 0x68, 0xae, 0x1a, 0xcf, 0x7d, 0x53, 0x57, 0xbd, 0x82, 0x6a, 0xbd,
	0x96, 0x8a, 0x99, 0x17, 0xdc, 0x25, 0x9b, 0xda, 0x0e, 0xc8, 0x44, 0xa9, 0xa7, 0x91, 0x3c, 0x83,
	0x88, 0xa4, 0x19, 0xe1, 0xd7, 0x72, 0xc5, 0x5d, 0xbe, 0xda, 0x81, 0xf7, 0x37, 0x07, 0x49, 0x1f,
	0x8e, 0x1a, 0xbc, 0x5a, 0x6c, 0xf2, 0xcd, 0x87, 0x80, 0xa4, 0x19, 0xbe, 0x85, 0x8e, 0xdb, 0x25,
	0x3e, 0x70, 0x02, 0x3b, 0x2f, 0x6a, 0x10, 0xff, 0xde, 0xb0, 0x9e, 0xe2, 0xcf, 0x3f, 0x7e, 0x7e,
	0xf5, 0x11, 0xa3, 0xf1, 0xea, 0xf1, 0x58, 0x31, 0x3d, 0x66, 0x4e, 0xec, 0x02, 0xc2, 0x6a, 0x97,
	0xd8, 0x77, 0xdc, 0xc6, 0xae, 0x07, 0xf7, 0xb7, 0x41, 0x2b, 0x76, 0x68, 0xc4, 0xba, 0xd8, 0xb6,
	0x62, 0x78, 0x01, 0xc1, 0x94, 0x31, 0xdc, 0x04, 0xbb, 0x5b, 0xec, 0xa0, 0xbf, 0x85, 0x59, 0x01,
	0x34, 0x02, 0x7b, 0x89, 0x13, 0x78, 0xee, 0x9d, 0xe3, 0x15, 0xb4, 0xea, 0xe8, 0x78, 0xec, 0x28,
	0x5b, 0x23, 0x1c, 0x9c, 0xec, 0xc2, 0xdb, 0x62, 0xe7, 0x0d, 0xb1, 0xf7, 0x2d, 0xf3, 0x87, 0x7b,
	0xfa, 0x2b, 0x00, 0x00, 0xff, 0xff, 0xf0, 0xed, 0xf1, 0x29, 0x11, 0x05, 0x00, 0x00,
}
