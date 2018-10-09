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
	Region               string   `protobuf:"bytes,3,opt,name=region,proto3" json:"region,omitempty"`
	Name                 string   `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MySQLNode) Reset()         { *m = MySQLNode{} }
func (m *MySQLNode) String() string { return proto.CompactTextString(m) }
func (*MySQLNode) ProtoMessage()    {}
func (*MySQLNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{0}
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

func (m *MySQLNode) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

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
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{1}
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

type MySQLInstanceID struct {
	Region               string   `protobuf:"bytes,1,opt,name=region,proto3" json:"region,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MySQLInstanceID) Reset()         { *m = MySQLInstanceID{} }
func (m *MySQLInstanceID) String() string { return proto.CompactTextString(m) }
func (*MySQLInstanceID) ProtoMessage()    {}
func (*MySQLInstanceID) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{2}
}
func (m *MySQLInstanceID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLInstanceID.Unmarshal(m, b)
}
func (m *MySQLInstanceID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLInstanceID.Marshal(b, m, deterministic)
}
func (dst *MySQLInstanceID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLInstanceID.Merge(dst, src)
}
func (m *MySQLInstanceID) XXX_Size() int {
	return xxx_messageInfo_MySQLInstanceID.Size(m)
}
func (m *MySQLInstanceID) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLInstanceID.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLInstanceID proto.InternalMessageInfo

func (m *MySQLInstanceID) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *MySQLInstanceID) GetName() string {
	if m != nil {
		return m.Name
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
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{3}
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

type MySQLDiscoverRequest struct {
	AwsAccessKeyId       string   `protobuf:"bytes,1,opt,name=aws_access_key_id,json=awsAccessKeyId,proto3" json:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey   string   `protobuf:"bytes,2,opt,name=aws_secret_access_key,json=awsSecretAccessKey,proto3" json:"aws_secret_access_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MySQLDiscoverRequest) Reset()         { *m = MySQLDiscoverRequest{} }
func (m *MySQLDiscoverRequest) String() string { return proto.CompactTextString(m) }
func (*MySQLDiscoverRequest) ProtoMessage()    {}
func (*MySQLDiscoverRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{4}
}
func (m *MySQLDiscoverRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLDiscoverRequest.Unmarshal(m, b)
}
func (m *MySQLDiscoverRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLDiscoverRequest.Marshal(b, m, deterministic)
}
func (dst *MySQLDiscoverRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLDiscoverRequest.Merge(dst, src)
}
func (m *MySQLDiscoverRequest) XXX_Size() int {
	return xxx_messageInfo_MySQLDiscoverRequest.Size(m)
}
func (m *MySQLDiscoverRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLDiscoverRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLDiscoverRequest proto.InternalMessageInfo

func (m *MySQLDiscoverRequest) GetAwsAccessKeyId() string {
	if m != nil {
		return m.AwsAccessKeyId
	}
	return ""
}

func (m *MySQLDiscoverRequest) GetAwsSecretAccessKey() string {
	if m != nil {
		return m.AwsSecretAccessKey
	}
	return ""
}

type MySQLDiscoverResponse struct {
	Instances            []*MySQLInstance `protobuf:"bytes,1,rep,name=instances,proto3" json:"instances,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *MySQLDiscoverResponse) Reset()         { *m = MySQLDiscoverResponse{} }
func (m *MySQLDiscoverResponse) String() string { return proto.CompactTextString(m) }
func (*MySQLDiscoverResponse) ProtoMessage()    {}
func (*MySQLDiscoverResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{5}
}
func (m *MySQLDiscoverResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MySQLDiscoverResponse.Unmarshal(m, b)
}
func (m *MySQLDiscoverResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MySQLDiscoverResponse.Marshal(b, m, deterministic)
}
func (dst *MySQLDiscoverResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLDiscoverResponse.Merge(dst, src)
}
func (m *MySQLDiscoverResponse) XXX_Size() int {
	return xxx_messageInfo_MySQLDiscoverResponse.Size(m)
}
func (m *MySQLDiscoverResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLDiscoverResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLDiscoverResponse proto.InternalMessageInfo

func (m *MySQLDiscoverResponse) GetInstances() []*MySQLInstance {
	if m != nil {
		return m.Instances
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
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{6}
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
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{7}
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
	Id                   *MySQLInstanceID `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username             string           `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Password             string           `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *MySQLAddRequest) Reset()         { *m = MySQLAddRequest{} }
func (m *MySQLAddRequest) String() string { return proto.CompactTextString(m) }
func (*MySQLAddRequest) ProtoMessage()    {}
func (*MySQLAddRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{8}
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

func (m *MySQLAddRequest) GetId() *MySQLInstanceID {
	if m != nil {
		return m.Id
	}
	return nil
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
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{9}
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
	Id                   *MySQLInstanceID `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *MySQLRemoveRequest) Reset()         { *m = MySQLRemoveRequest{} }
func (m *MySQLRemoveRequest) String() string { return proto.CompactTextString(m) }
func (*MySQLRemoveRequest) ProtoMessage()    {}
func (*MySQLRemoveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{10}
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

func (m *MySQLRemoveRequest) GetId() *MySQLInstanceID {
	if m != nil {
		return m.Id
	}
	return nil
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
	return fileDescriptor_mysql_7fb2dce08bcdf289, []int{11}
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
	proto.RegisterType((*MySQLInstanceID)(nil), "api.MySQLInstanceID")
	proto.RegisterType((*MySQLInstance)(nil), "api.MySQLInstance")
	proto.RegisterType((*MySQLDiscoverRequest)(nil), "api.MySQLDiscoverRequest")
	proto.RegisterType((*MySQLDiscoverResponse)(nil), "api.MySQLDiscoverResponse")
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
	Discover(ctx context.Context, in *MySQLDiscoverRequest, opts ...grpc.CallOption) (*MySQLDiscoverResponse, error)
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

func (c *mySQLClient) Discover(ctx context.Context, in *MySQLDiscoverRequest, opts ...grpc.CallOption) (*MySQLDiscoverResponse, error) {
	out := new(MySQLDiscoverResponse)
	err := c.cc.Invoke(ctx, "/api.MySQL/Discover", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
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
	Discover(context.Context, *MySQLDiscoverRequest) (*MySQLDiscoverResponse, error)
	List(context.Context, *MySQLListRequest) (*MySQLListResponse, error)
	Add(context.Context, *MySQLAddRequest) (*MySQLAddResponse, error)
	Remove(context.Context, *MySQLRemoveRequest) (*MySQLRemoveResponse, error)
}

func RegisterMySQLServer(s *grpc.Server, srv MySQLServer) {
	s.RegisterService(&_MySQL_serviceDesc, srv)
}

func _MySQL_Discover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MySQLDiscoverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MySQLServer).Discover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.MySQL/Discover",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MySQLServer).Discover(ctx, req.(*MySQLDiscoverRequest))
	}
	return interceptor(ctx, in, info, handler)
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
			MethodName: "Discover",
			Handler:    _MySQL_Discover_Handler,
		},
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

func init() { proto.RegisterFile("mysql.proto", fileDescriptor_mysql_7fb2dce08bcdf289) }

var fileDescriptor_mysql_7fb2dce08bcdf289 = []byte{
	// 584 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x95, 0x1d, 0x37, 0x4d, 0x26, 0x34, 0x24, 0x43, 0xd3, 0x1a, 0x0b, 0xa4, 0x6a, 0x05, 0x52,
	0x29, 0x52, 0x53, 0xc2, 0xad, 0x12, 0x87, 0x8a, 0x22, 0x35, 0xa5, 0x20, 0xd5, 0x95, 0xb8, 0xa6,
	0x4b, 0x76, 0x14, 0xad, 0x68, 0xbc, 0xae, 0xd7, 0x4d, 0x94, 0x2b, 0x77, 0x4e, 0xfc, 0x18, 0x12,
	0xbf, 0xc0, 0x87, 0x20, 0xaf, 0x37, 0xb6, 0x93, 0xf6, 0x02, 0xb7, 0x9d, 0x37, 0x33, 0x6f, 0xde,
	0x9b, 0x71, 0x02, 0xad, 0xe9, 0x42, 0xdf, 0xde, 0x1c, 0xc6, 0x89, 0x4a, 0x15, 0xd6, 0x78, 0x2c,
	0x83, 0x67, 0x13, 0xa5, 0x26, 0x37, 0xd4, 0xe7, 0xb1, 0xec, 0xf3, 0x28, 0x52, 0x29, 0x4f, 0xa5,
	0x8a, 0x74, 0x5e, 0xc2, 0xde, 0x43, 0xf3, 0xd3, 0xe2, 0xea, 0xf2, 0xe2, 0xb3, 0x12, 0x84, 0x3b,
	0x50, 0x4f, 0x68, 0x22, 0x55, 0xe4, 0xd7, 0xf6, 0x9c, 0xfd, 0x66, 0x68, 0x23, 0x44, 0xf0, 0x22,
	0x3e, 0x25, 0xdf, 0x33, 0xa8, 0x79, 0x9f, 0x7b, 0x0d, 0xa7, 0xe3, 0x9e, 0x7b, 0x0d, 0xb7, 0x53,
	0x63, 0x3f, 0x1c, 0x78, 0x64, 0x58, 0xae, 0x28, 0x99, 0xc9, 0x31, 0xa1, 0x0f, 0x9b, 0x5c, 0x88,
	0x84, 0xb4, 0xb6, 0x3d, 0xcb, 0x30, 0xa3, 0x8a, 0x55, 0x92, 0xfa, 0x1b, 0x7b, 0xce, 0xfe, 0x56,
	0x68, 0xde, 0xd9, 0x58, 0x8a, 0x26, 0x32, 0x22, 0xbf, 0x9e, 0x8f, 0xcd, 0x23, 0x7c, 0x09, 0xed,
	0xfc, 0x35, 0x9a, 0x51, 0xa2, 0x33, 0x59, 0x9b, 0x26, 0xbf, 0x95, 0xa3, 0x5f, 0x72, 0xb0, 0xaa,
	0xe4, 0xdc, 0x6b, 0xd4, 0x3a, 0x1e, 0x7b, 0x07, 0x8f, 0x8d, 0x9c, 0x61, 0xa4, 0x53, 0x1e, 0x8d,
	0x69, 0x78, 0x5a, 0xb1, 0xe6, 0x3c, 0x68, 0xcd, 0x2d, 0xad, 0xb1, 0x6b, 0xd8, 0x5a, 0x69, 0x47,
	0x06, 0x5e, 0xa4, 0x04, 0x99, 0xd6, 0xd6, 0xa0, 0x7d, 0xc8, 0x63, 0x79, 0x58, 0x6c, 0x2d, 0x34,
	0x39, 0x7c, 0x0d, 0x9b, 0x3a, 0x77, 0x6f, 0xb8, 0x5a, 0x83, 0x6e, 0x59, 0x66, 0xd7, 0x12, 0x2e,
	0x2b, 0x58, 0x0a, 0xdb, 0x26, 0x71, 0x2a, 0xf5, 0x58, 0xcd, 0x28, 0x09, 0xe9, 0xf6, 0x8e, 0x74,
	0x8a, 0xaf, 0xa0, 0xcb, 0xe7, 0x7a, 0xc4, 0xc7, 0x63, 0xd2, 0x7a, 0xf4, 0x8d, 0x16, 0x23, 0x29,
	0xac, 0xe0, 0x36, 0x9f, 0xeb, 0x13, 0x83, 0x7f, 0xa4, 0xc5, 0x50, 0xe0, 0x1b, 0xe8, 0x65, 0xa5,
	0x9a, 0xc6, 0x09, 0xa5, 0x95, 0x0e, 0xeb, 0x04, 0xf9, 0x5c, 0x5f, 0x99, 0x5c, 0xd1, 0xc4, 0x86,
	0xd0, 0x5b, 0x9b, 0xaa, 0x63, 0x15, 0x69, 0xc2, 0x23, 0x68, 0x4a, 0xeb, 0x55, 0xfb, 0xce, 0x5e,
	0x6d, 0xbf, 0x35, 0xc0, 0x52, 0xfd, 0x72, 0x0d, 0x61, 0x59, 0xc4, 0x10, 0x3a, 0x26, 0x77, 0x21,
	0x75, 0x6a, 0xc5, 0xb3, 0x0f, 0xd0, 0xad, 0x60, 0xff, 0x4d, 0xad, 0xec, 0xf1, 0x4e, 0x84, 0x58,
	0xae, 0xe5, 0x05, 0xb8, 0x76, 0x0f, 0xad, 0xc1, 0xf6, 0xfd, 0xee, 0xe1, 0x69, 0xe8, 0x4a, 0x81,
	0x01, 0x34, 0xee, 0x34, 0x25, 0x95, 0x73, 0x16, 0x71, 0x96, 0x8b, 0xb9, 0xd6, 0x73, 0x95, 0x08,
	0xfb, 0x6d, 0x17, 0x71, 0xe1, 0xc5, 0x0c, 0xcc, 0x65, 0xb3, 0x63, 0x40, 0x83, 0x85, 0x34, 0x55,
	0x33, 0xfa, 0x27, 0x1d, 0xac, 0x07, 0x4f, 0x56, 0x7a, 0x73, 0xca, 0xc1, 0x2f, 0x17, 0x36, 0x0c,
	0x8e, 0xd7, 0xd0, 0x58, 0x9e, 0x00, 0x9f, 0x96, 0x34, 0x6b, 0x1f, 0x43, 0x10, 0x3c, 0x94, 0xb2,
	0xfa, 0x9e, 0x7f, 0xff, 0xfd, 0xe7, 0xa7, 0xbb, 0xcb, 0xb0, 0x3f, 0x3b, 0xea, 0x9b, 0x9f, 0x7c,
	0x5f, 0xd8, 0x9a, 0x63, 0xe7, 0x00, 0xcf, 0xc0, 0xcb, 0xae, 0x80, 0xbd, 0x92, 0xa2, 0x72, 0xa9,
	0x60, 0x67, 0x1d, 0xb6, 0xac, 0x5d, 0xc3, 0xda, 0xc2, 0x66, 0xc1, 0x8a, 0x67, 0x50, 0x3b, 0x11,
	0x02, 0x2b, 0x6e, 0xcb, 0xbb, 0x04, 0xbd, 0x35, 0xd4, 0xd2, 0x6c, 0x1b, 0x9a, 0x36, 0x2b, 0x69,
	0x32, 0x4d, 0x97, 0x50, 0xcf, 0x37, 0x82, 0xbb, 0x65, 0xdb, 0xca, 0x7e, 0x03, 0xff, 0x7e, 0x62,
	0x95, 0xf2, 0x60, 0x85, 0xf2, 0x6b, 0xdd, 0xfc, 0x87, 0xbd, 0xfd, 0x1b, 0x00, 0x00, 0xff, 0xff,
	0xa2, 0x21, 0x97, 0x86, 0xf5, 0x04, 0x00, 0x00,
}
