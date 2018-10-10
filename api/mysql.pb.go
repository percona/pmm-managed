// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mysql.proto

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

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MySQLNode struct {
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MySQLNode) Reset()         { *m = MySQLNode{} }
func (m *MySQLNode) String() string { return proto.CompactTextString(m) }
func (*MySQLNode) ProtoMessage()    {}
func (*MySQLNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_d214ec27dcd8cedc, []int{0}
}
func (m *MySQLNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLNode.Unmarshal(m, b)
}
func (m *MySQLNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLNode.Marshal(b, m, deterministic)
}
func (dst *MySQLNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLNode.Merge(dst, src)
}
func (m *MySQLNode) XXX_Size() int {
	return xxx_messageInfo_MySQLNode.Size(m)
}
func (m *MySQLNode) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLNode.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLNode proto.InternalMessageInfo

func (m *MySQLNode) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type MySQLService struct {
	Address              string   `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty"`
	Port                 uint32   `protobuf:"varint,5,opt,name=port,proto3" json:"port,omitempty"`
	Engine               string   `protobuf:"bytes,6,opt,name=engine,proto3" json:"engine,omitempty"`
	EngineVersion        string   `protobuf:"bytes,7,opt,name=engine_version,json=engineVersion,proto3" json:"engine_version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MySQLService) Reset()         { *m = MySQLService{} }
func (m *MySQLService) String() string { return proto.CompactTextString(m) }
func (*MySQLService) ProtoMessage()    {}
func (*MySQLService) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_d214ec27dcd8cedc, []int{1}
}
func (m *MySQLService) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLService.Unmarshal(m, b)
}
func (m *MySQLService) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLService.Marshal(b, m, deterministic)
}
func (dst *MySQLService) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLService.Merge(dst, src)
}
func (m *MySQLService) XXX_Size() int {
	return xxx_messageInfo_MySQLService.Size(m)
}
func (m *MySQLService) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLService.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLService proto.InternalMessageInfo

func (m *MySQLService) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *MySQLService) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *MySQLService) GetEngine() string {
	if m != nil {
		return m.Engine
	}
	return ""
}

func (m *MySQLService) GetEngineVersion() string {
	if m != nil {
		return m.EngineVersion
	}
	return ""
}

type MySQLInstance struct {
	Node                 *MySQLNode    `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	Service              *MySQLService `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *MySQLInstance) Reset()         { *m = MySQLInstance{} }
func (m *MySQLInstance) String() string { return proto.CompactTextString(m) }
func (*MySQLInstance) ProtoMessage()    {}
func (*MySQLInstance) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_d214ec27dcd8cedc, []int{2}
}
func (m *MySQLInstance) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLInstance.Unmarshal(m, b)
}
func (m *MySQLInstance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLInstance.Marshal(b, m, deterministic)
}
func (dst *MySQLInstance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLInstance.Merge(dst, src)
}
func (m *MySQLInstance) XXX_Size() int {
	return xxx_messageInfo_MySQLInstance.Size(m)
}
func (m *MySQLInstance) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLInstance.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLInstance proto.InternalMessageInfo

func (m *MySQLInstance) GetNode() *MySQLNode {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *MySQLInstance) GetService() *MySQLService {
	if m != nil {
		return m.Service
	}
	return nil
}

type MySQLListRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MySQLListRequest) Reset()         { *m = MySQLListRequest{} }
func (m *MySQLListRequest) String() string { return proto.CompactTextString(m) }
func (*MySQLListRequest) ProtoMessage()    {}
func (*MySQLListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_d214ec27dcd8cedc, []int{3}
}
func (m *MySQLListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLListRequest.Unmarshal(m, b)
}
func (m *MySQLListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLListRequest.Marshal(b, m, deterministic)
}
func (dst *MySQLListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLListRequest.Merge(dst, src)
}
func (m *MySQLListRequest) XXX_Size() int {
	return xxx_messageInfo_MySQLListRequest.Size(m)
}
func (m *MySQLListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLListRequest proto.InternalMessageInfo

type MySQLListResponse struct {
	Instances            []*MySQLInstance `protobuf:"bytes,1,rep,name=instances,proto3" json:"instances,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *MySQLListResponse) Reset()         { *m = MySQLListResponse{} }
func (m *MySQLListResponse) String() string { return proto.CompactTextString(m) }
func (*MySQLListResponse) ProtoMessage()    {}
func (*MySQLListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_d214ec27dcd8cedc, []int{4}
}
func (m *MySQLListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLListResponse.Unmarshal(m, b)
}
func (m *MySQLListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLListResponse.Marshal(b, m, deterministic)
}
func (dst *MySQLListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLListResponse.Merge(dst, src)
}
func (m *MySQLListResponse) XXX_Size() int {
	return xxx_messageInfo_MySQLListResponse.Size(m)
}
func (m *MySQLListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLListResponse proto.InternalMessageInfo

func (m *MySQLListResponse) GetInstances() []*MySQLInstance {
	if m != nil {
		return m.Instances
	}
	return nil
}

type MySQLAddRequest struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Port                 uint32   `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Username             string   `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MySQLAddRequest) Reset()         { *m = MySQLAddRequest{} }
func (m *MySQLAddRequest) String() string { return proto.CompactTextString(m) }
func (*MySQLAddRequest) ProtoMessage()    {}
func (*MySQLAddRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_d214ec27dcd8cedc, []int{5}
}
func (m *MySQLAddRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLAddRequest.Unmarshal(m, b)
}
func (m *MySQLAddRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLAddRequest.Marshal(b, m, deterministic)
}
func (dst *MySQLAddRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLAddRequest.Merge(dst, src)
}
func (m *MySQLAddRequest) XXX_Size() int {
	return xxx_messageInfo_MySQLAddRequest.Size(m)
}
func (m *MySQLAddRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLAddRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLAddRequest proto.InternalMessageInfo

func (m *MySQLAddRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *MySQLAddRequest) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *MySQLAddRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *MySQLAddRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type MySQLAddResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MySQLAddResponse) Reset()         { *m = MySQLAddResponse{} }
func (m *MySQLAddResponse) String() string { return proto.CompactTextString(m) }
func (*MySQLAddResponse) ProtoMessage()    {}
func (*MySQLAddResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_d214ec27dcd8cedc, []int{6}
}
func (m *MySQLAddResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLAddResponse.Unmarshal(m, b)
}
func (m *MySQLAddResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLAddResponse.Marshal(b, m, deterministic)
}
func (dst *MySQLAddResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLAddResponse.Merge(dst, src)
}
func (m *MySQLAddResponse) XXX_Size() int {
	return xxx_messageInfo_MySQLAddResponse.Size(m)
}
func (m *MySQLAddResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLAddResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLAddResponse proto.InternalMessageInfo

type MySQLRemoveRequest struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Port                 uint32   `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MySQLRemoveRequest) Reset()         { *m = MySQLRemoveRequest{} }
func (m *MySQLRemoveRequest) String() string { return proto.CompactTextString(m) }
func (*MySQLRemoveRequest) ProtoMessage()    {}
func (*MySQLRemoveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_d214ec27dcd8cedc, []int{7}
}
func (m *MySQLRemoveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLRemoveRequest.Unmarshal(m, b)
}
func (m *MySQLRemoveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLRemoveRequest.Marshal(b, m, deterministic)
}
func (dst *MySQLRemoveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLRemoveRequest.Merge(dst, src)
}
func (m *MySQLRemoveRequest) XXX_Size() int {
	return xxx_messageInfo_MySQLRemoveRequest.Size(m)
}
func (m *MySQLRemoveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLRemoveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLRemoveRequest proto.InternalMessageInfo

func (m *MySQLRemoveRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *MySQLRemoveRequest) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

type MySQLRemoveResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MySQLRemoveResponse) Reset()         { *m = MySQLRemoveResponse{} }
func (m *MySQLRemoveResponse) String() string { return proto.CompactTextString(m) }
func (*MySQLRemoveResponse) ProtoMessage()    {}
func (*MySQLRemoveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_d214ec27dcd8cedc, []int{8}
}
func (m *MySQLRemoveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLRemoveResponse.Unmarshal(m, b)
}
func (m *MySQLRemoveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLRemoveResponse.Marshal(b, m, deterministic)
}
func (dst *MySQLRemoveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLRemoveResponse.Merge(dst, src)
}
func (m *MySQLRemoveResponse) XXX_Size() int {
	return xxx_messageInfo_MySQLRemoveResponse.Size(m)
}
func (m *MySQLRemoveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLRemoveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLRemoveResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MySQLNode)(nil), "api.MySQLNode")
	proto.RegisterType((*MySQLService)(nil), "api.MySQLService")
	proto.RegisterType((*MySQLInstance)(nil), "api.MySQLInstance")
	proto.RegisterType((*MySQLListRequest)(nil), "api.MySQLListRequest")
	proto.RegisterType((*MySQLListResponse)(nil), "api.MySQLListResponse")
	proto.RegisterType((*MySQLAddRequest)(nil), "api.MySQLAddRequest")
	proto.RegisterType((*MySQLAddResponse)(nil), "api.MySQLAddResponse")
	proto.RegisterType((*MySQLRemoveRequest)(nil), "api.MySQLRemoveRequest")
	proto.RegisterType((*MySQLRemoveResponse)(nil), "api.MySQLRemoveResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MySQLClient is the client API for MySQL service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MySQLClient interface {
	List(ctx context.Context, in *MySQLListRequest, opts ...grpc.CallOption) (*MySQLListResponse, error)
	Add(ctx context.Context, in *MySQLAddRequest, opts ...grpc.CallOption) (*MySQLAddResponse, error)
	Remove(ctx context.Context, in *MySQLRemoveRequest, opts ...grpc.CallOption) (*MySQLRemoveResponse, error)
}

type mySQLClient struct {
	cc *grpc.ClientConn
}

func NewMySQLClient(cc *grpc.ClientConn) MySQLClient {
	return &mySQLClient{cc}
}

func (c *mySQLClient) List(ctx context.Context, in *MySQLListRequest, opts ...grpc.CallOption) (*MySQLListResponse, error) {
	out := new(MySQLListResponse)
	err := c.cc.Invoke(ctx, "/api.MySQL/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mySQLClient) Add(ctx context.Context, in *MySQLAddRequest, opts ...grpc.CallOption) (*MySQLAddResponse, error) {
	out := new(MySQLAddResponse)
	err := c.cc.Invoke(ctx, "/api.MySQL/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mySQLClient) Remove(ctx context.Context, in *MySQLRemoveRequest, opts ...grpc.CallOption) (*MySQLRemoveResponse, error) {
	out := new(MySQLRemoveResponse)
	err := c.cc.Invoke(ctx, "/api.MySQL/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MySQLServer is the server API for MySQL service.
type MySQLServer interface {
	List(context.Context, *MySQLListRequest) (*MySQLListResponse, error)
	Add(context.Context, *MySQLAddRequest) (*MySQLAddResponse, error)
	Remove(context.Context, *MySQLRemoveRequest) (*MySQLRemoveResponse, error)
}

func RegisterMySQLServer(s *grpc.Server, srv MySQLServer) {
	s.RegisterService(&_MySQL_serviceDesc, srv)
}

func _MySQL_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MySQLListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MySQLServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.MySQL/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MySQLServer).List(ctx, req.(*MySQLListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MySQL_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MySQLAddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MySQLServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.MySQL/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MySQLServer).Add(ctx, req.(*MySQLAddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MySQL_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MySQLRemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MySQLServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.MySQL/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MySQLServer).Remove(ctx, req.(*MySQLRemoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MySQL_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.MySQL",
	HandlerType: (*MySQLServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _MySQL_List_Handler,
		},
		{
			MethodName: "Add",
			Handler:    _MySQL_Add_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _MySQL_Remove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mysql.proto",
}

func init() { proto.RegisterFile("mysql.proto", fileDescriptor_mysql_d214ec27dcd8cedc) }

var fileDescriptor_mysql_d214ec27dcd8cedc = []byte{
	// 459 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0x4d, 0x6b, 0xdb, 0x40,
	0x14, 0x64, 0x25, 0xc5, 0x1f, 0xcf, 0xb5, 0x6b, 0xbf, 0xc6, 0xe9, 0x22, 0x7a, 0x30, 0x0b, 0x05,
	0x93, 0x80, 0x1d, 0xdc, 0x5b, 0x6f, 0x29, 0x14, 0x52, 0x93, 0x16, 0xa2, 0x40, 0xaf, 0xed, 0x36,
	0xbb, 0x98, 0x85, 0x78, 0x57, 0xd1, 0x2a, 0x0e, 0xb9, 0xf6, 0xde, 0x53, 0x7f, 0x5a, 0xff, 0x42,
	0xcf, 0xfd, 0x0d, 0xc5, 0x4f, 0xb2, 0x25, 0xa5, 0x3d, 0xe5, 0xb6, 0x6f, 0x66, 0x76, 0x3c, 0x6f,
	0xd6, 0x82, 0xde, 0xfa, 0xc1, 0xdf, 0xde, 0xcc, 0xd2, 0xcc, 0xe5, 0x0e, 0x43, 0x99, 0x9a, 0xf8,
	0xd5, 0xca, 0xb9, 0xd5, 0x8d, 0x9e, 0xcb, 0xd4, 0xcc, 0xa5, 0xb5, 0x2e, 0x97, 0xb9, 0x71, 0xd6,
	0x17, 0x12, 0x71, 0x02, 0xdd, 0x8f, 0x0f, 0x57, 0x97, 0x17, 0x9f, 0x9c, 0xd2, 0x88, 0x10, 0x59,
	0xb9, 0xd6, 0x3c, 0x9c, 0xb0, 0x69, 0x37, 0xa1, 0xf3, 0x32, 0xea, 0xb0, 0x61, 0xb0, 0x8c, 0x3a,
	0xc1, 0x30, 0x14, 0x3f, 0x18, 0x3c, 0x23, 0xf5, 0x95, 0xce, 0x36, 0xe6, 0x5a, 0x23, 0x87, 0xb6,
	0x54, 0x2a, 0xd3, 0xde, 0xf3, 0x88, 0xee, 0xec, 0xc6, 0xad, 0x55, 0xea, 0xb2, 0x9c, 0x1f, 0x4c,
	0xd8, 0xb4, 0x9f, 0xd0, 0x19, 0x8f, 0xa0, 0xa5, 0xed, 0xca, 0x58, 0xcd, 0x5b, 0x24, 0x2e, 0x27,
	0x7c, 0x0d, 0x83, 0xe2, 0xf4, 0x65, 0xa3, 0x33, 0x6f, 0x9c, 0xe5, 0x6d, 0xe2, 0xfb, 0x05, 0xfa,
	0xb9, 0x00, 0xeb, 0x49, 0x96, 0x51, 0x27, 0x1c, 0x46, 0xe2, 0x2b, 0xf4, 0x29, 0xce, 0x07, 0xeb,
	0x73, 0x69, 0xaf, 0x35, 0x0a, 0x88, 0xac, 0x53, 0x9a, 0xb3, 0x09, 0x9b, 0xf6, 0x16, 0x83, 0x99,
	0x4c, 0xcd, 0x6c, 0xbf, 0x5e, 0x42, 0x1c, 0x9e, 0x40, 0xdb, 0x17, 0xf1, 0x79, 0x40, 0xb2, 0x51,
	0x25, 0x2b, 0xf7, 0x4a, 0x76, 0x0a, 0x81, 0x30, 0x24, 0xe2, 0xc2, 0xf8, 0x3c, 0xd1, 0xb7, 0x77,
	0xda, 0xe7, 0xe2, 0x3d, 0x8c, 0x6a, 0x98, 0x4f, 0x9d, 0xf5, 0x1a, 0x4f, 0xa1, 0x6b, 0xca, 0x14,
	0x9e, 0xb3, 0x49, 0x38, 0xed, 0x2d, 0xb0, 0xf2, 0xdd, 0x05, 0x4c, 0x2a, 0x91, 0xb8, 0x87, 0xe7,
	0xc4, 0x9d, 0x29, 0x55, 0x3a, 0xd7, 0xeb, 0x64, 0xff, 0xaf, 0x33, 0xa8, 0xd5, 0x19, 0x43, 0xe7,
	0xce, 0xeb, 0xac, 0xf6, 0x62, 0xfb, 0x79, 0xcb, 0xa5, 0xd2, 0xfb, 0x7b, 0x97, 0xa9, 0xf2, 0x65,
	0xf6, 0xf3, 0x7e, 0x27, 0xfa, 0xe1, 0x22, 0xbe, 0x78, 0x07, 0x48, 0x58, 0xa2, 0xd7, 0x6e, 0xa3,
	0x9f, 0x94, 0x47, 0x8c, 0xe1, 0x45, 0xc3, 0xa3, 0xb0, 0x5e, 0xfc, 0x61, 0x70, 0x40, 0x38, 0x9e,
	0x43, 0xb4, 0xed, 0x0c, 0xc7, 0x55, 0x31, 0xb5, 0x5e, 0xe3, 0xa3, 0xc7, 0x70, 0x99, 0x6d, 0xf4,
	0xfd, 0xd7, 0xef, 0x9f, 0x41, 0x0f, 0xbb, 0xf3, 0xcd, 0xe9, 0x9c, 0xfe, 0xde, 0x78, 0x0e, 0xe1,
	0x99, 0x52, 0x78, 0x58, 0xdd, 0xa8, 0x5a, 0x8c, 0xc7, 0x8f, 0xd0, 0xd2, 0xe6, 0x90, 0x6c, 0x06,
	0xa2, 0xb2, 0x79, 0xcb, 0x8e, 0xf1, 0x12, 0x5a, 0x45, 0x5e, 0x7c, 0x59, 0x5d, 0x6b, 0xb4, 0x10,
	0xf3, 0x7f, 0x89, 0xa6, 0xe5, 0x71, 0xc3, 0xf2, 0x5b, 0x8b, 0xbe, 0xac, 0x37, 0x7f, 0x03, 0x00,
	0x00, 0xff, 0xff, 0x18, 0xcb, 0x93, 0xe5, 0x8b, 0x03, 0x00, 0x00,
}
