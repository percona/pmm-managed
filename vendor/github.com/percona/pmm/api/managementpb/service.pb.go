// Code generated by protoc-gen-go. DO NOT EDIT.
// source: managementpb/service.proto

package managementpb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "github.com/mwitkow/go-proto-validators"
	inventorypb "github.com/percona/pmm/api/inventorypb"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// AddNodeParams is a params to add new node to inventory while adding new service.
type AddNodeParams struct {
	// Node type to be registered.
	NodeType inventorypb.NodeType `protobuf:"varint,1,opt,name=node_type,json=nodeType,proto3,enum=inventory.NodeType" json:"node_type,omitempty"`
	// Unique across all Nodes user-defined name.
	NodeName string `protobuf:"bytes,2,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	// Linux machine-id.
	MachineId string `protobuf:"bytes,3,opt,name=machine_id,json=machineId,proto3" json:"machine_id,omitempty"`
	// Linux distribution name and version.
	Distro string `protobuf:"bytes,4,opt,name=distro,proto3" json:"distro,omitempty"`
	// Container identifier. If specified, must be a unique Docker container identifier.
	ContainerId string `protobuf:"bytes,5,opt,name=container_id,json=containerId,proto3" json:"container_id,omitempty"`
	// Container name.
	ContainerName string `protobuf:"bytes,6,opt,name=container_name,json=containerName,proto3" json:"container_name,omitempty"`
	// Node model.
	NodeModel string `protobuf:"bytes,7,opt,name=node_model,json=nodeModel,proto3" json:"node_model,omitempty"`
	// Node region.
	Region string `protobuf:"bytes,8,opt,name=region,proto3" json:"region,omitempty"`
	// Node availability zone.
	Az string `protobuf:"bytes,9,opt,name=az,proto3" json:"az,omitempty"`
	// Custom user-assigned labels for Node.
	CustomLabels         map[string]string `protobuf:"bytes,10,rep,name=custom_labels,json=customLabels,proto3" json:"custom_labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *AddNodeParams) Reset()         { *m = AddNodeParams{} }
func (m *AddNodeParams) String() string { return proto.CompactTextString(m) }
func (*AddNodeParams) ProtoMessage()    {}
func (*AddNodeParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_3cb53daa1555e7e7, []int{0}
}

func (m *AddNodeParams) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddNodeParams.Unmarshal(m, b)
}
func (m *AddNodeParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddNodeParams.Marshal(b, m, deterministic)
}
func (m *AddNodeParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddNodeParams.Merge(m, src)
}
func (m *AddNodeParams) XXX_Size() int {
	return xxx_messageInfo_AddNodeParams.Size(m)
}
func (m *AddNodeParams) XXX_DiscardUnknown() {
	xxx_messageInfo_AddNodeParams.DiscardUnknown(m)
}

var xxx_messageInfo_AddNodeParams proto.InternalMessageInfo

func (m *AddNodeParams) GetNodeType() inventorypb.NodeType {
	if m != nil {
		return m.NodeType
	}
	return inventorypb.NodeType_NODE_TYPE_INVALID
}

func (m *AddNodeParams) GetNodeName() string {
	if m != nil {
		return m.NodeName
	}
	return ""
}

func (m *AddNodeParams) GetMachineId() string {
	if m != nil {
		return m.MachineId
	}
	return ""
}

func (m *AddNodeParams) GetDistro() string {
	if m != nil {
		return m.Distro
	}
	return ""
}

func (m *AddNodeParams) GetContainerId() string {
	if m != nil {
		return m.ContainerId
	}
	return ""
}

func (m *AddNodeParams) GetContainerName() string {
	if m != nil {
		return m.ContainerName
	}
	return ""
}

func (m *AddNodeParams) GetNodeModel() string {
	if m != nil {
		return m.NodeModel
	}
	return ""
}

func (m *AddNodeParams) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *AddNodeParams) GetAz() string {
	if m != nil {
		return m.Az
	}
	return ""
}

func (m *AddNodeParams) GetCustomLabels() map[string]string {
	if m != nil {
		return m.CustomLabels
	}
	return nil
}

type RemoveServiceRequest struct {
	// Service type.
	ServiceType inventorypb.ServiceType `protobuf:"varint,1,opt,name=service_type,json=serviceType,proto3,enum=inventory.ServiceType" json:"service_type,omitempty"`
	// Service ID or Service Name is required.
	// Unique randomly generated instance identifier.
	ServiceId string `protobuf:"bytes,2,opt,name=service_id,json=serviceId,proto3" json:"service_id,omitempty"`
	// Unique across all Services user-defined name.
	ServiceName          string   `protobuf:"bytes,3,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveServiceRequest) Reset()         { *m = RemoveServiceRequest{} }
func (m *RemoveServiceRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveServiceRequest) ProtoMessage()    {}
func (*RemoveServiceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_3cb53daa1555e7e7, []int{1}
}

func (m *RemoveServiceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveServiceRequest.Unmarshal(m, b)
}
func (m *RemoveServiceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveServiceRequest.Marshal(b, m, deterministic)
}
func (m *RemoveServiceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveServiceRequest.Merge(m, src)
}
func (m *RemoveServiceRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveServiceRequest.Size(m)
}
func (m *RemoveServiceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveServiceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveServiceRequest proto.InternalMessageInfo

func (m *RemoveServiceRequest) GetServiceType() inventorypb.ServiceType {
	if m != nil {
		return m.ServiceType
	}
	return inventorypb.ServiceType_SERVICE_TYPE_INVALID
}

func (m *RemoveServiceRequest) GetServiceId() string {
	if m != nil {
		return m.ServiceId
	}
	return ""
}

func (m *RemoveServiceRequest) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

type RemoveServiceResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveServiceResponse) Reset()         { *m = RemoveServiceResponse{} }
func (m *RemoveServiceResponse) String() string { return proto.CompactTextString(m) }
func (*RemoveServiceResponse) ProtoMessage()    {}
func (*RemoveServiceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_3cb53daa1555e7e7, []int{2}
}

func (m *RemoveServiceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveServiceResponse.Unmarshal(m, b)
}
func (m *RemoveServiceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveServiceResponse.Marshal(b, m, deterministic)
}
func (m *RemoveServiceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveServiceResponse.Merge(m, src)
}
func (m *RemoveServiceResponse) XXX_Size() int {
	return xxx_messageInfo_RemoveServiceResponse.Size(m)
}
func (m *RemoveServiceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveServiceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveServiceResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*AddNodeParams)(nil), "management.AddNodeParams")
	proto.RegisterMapType((map[string]string)(nil), "management.AddNodeParams.CustomLabelsEntry")
	proto.RegisterType((*RemoveServiceRequest)(nil), "management.RemoveServiceRequest")
	proto.RegisterType((*RemoveServiceResponse)(nil), "management.RemoveServiceResponse")
}

func init() { proto.RegisterFile("managementpb/service.proto", fileDescriptor_3cb53daa1555e7e7) }

var fileDescriptor_3cb53daa1555e7e7 = []byte{
	// 608 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0xc1, 0x6e, 0xd3, 0x4c,
	0x10, 0x96, 0x93, 0xbf, 0x69, 0xb3, 0x4d, 0xaa, 0x9f, 0xa5, 0xb4, 0x96, 0x45, 0x45, 0x9a, 0x0a,
	0x11, 0x28, 0xf1, 0x42, 0x91, 0x10, 0xf4, 0x82, 0x52, 0xd4, 0x43, 0xa5, 0x52, 0x55, 0x86, 0x03,
	0xe2, 0x12, 0x6d, 0xbc, 0x53, 0x77, 0x55, 0x7b, 0xd7, 0xac, 0x37, 0x89, 0xd2, 0x03, 0x07, 0x6e,
	0x5c, 0xe9, 0x13, 0xf1, 0x0c, 0x3c, 0x00, 0x12, 0xe2, 0x41, 0x90, 0xd7, 0xeb, 0xc4, 0xa5, 0x88,
	0x53, 0x76, 0xbe, 0xef, 0x9b, 0x99, 0x2f, 0x33, 0xbb, 0x46, 0x5e, 0x42, 0x05, 0x8d, 0x20, 0x01,
	0xa1, 0xd3, 0x11, 0xc9, 0x40, 0x4d, 0x78, 0x08, 0x7e, 0xaa, 0xa4, 0x96, 0x18, 0x2d, 0x38, 0xef,
	0x79, 0xc4, 0xf5, 0xf9, 0x78, 0xe4, 0x87, 0x32, 0x21, 0xc9, 0x94, 0xeb, 0x0b, 0x39, 0x25, 0x91,
	0xec, 0x1b, 0x61, 0x7f, 0x42, 0x63, 0xce, 0xa8, 0x96, 0x2a, 0x23, 0xf3, 0x63, 0x51, 0xc3, 0xbb,
	0x1b, 0x49, 0x19, 0xc5, 0x40, 0x68, 0xca, 0x09, 0x15, 0x42, 0x6a, 0xaa, 0xb9, 0x14, 0x99, 0x65,
	0x37, 0xb9, 0x98, 0x80, 0xd0, 0x52, 0xcd, 0xd2, 0x11, 0x11, 0x92, 0x41, 0x49, 0x78, 0x55, 0xc2,
	0xba, 0x2a, 0xb9, 0xc7, 0xe6, 0x27, 0xec, 0x47, 0x20, 0xfa, 0xd9, 0x94, 0x46, 0x11, 0x28, 0x22,
	0x53, 0x53, 0xf6, 0x66, 0x8b, 0xee, 0xb7, 0x3a, 0x6a, 0x0f, 0x18, 0x3b, 0x91, 0x0c, 0x4e, 0xa9,
	0xa2, 0x49, 0x86, 0x9f, 0xa0, 0x66, 0xde, 0x6a, 0xa8, 0x67, 0x29, 0xb8, 0x4e, 0xc7, 0xe9, 0xad,
	0xed, 0xdd, 0xf6, 0xe7, 0xfd, 0xfc, 0x5c, 0xf9, 0x6e, 0x96, 0x42, 0xb0, 0x22, 0xec, 0x09, 0xef,
	0xd8, 0x0c, 0x41, 0x13, 0x70, 0x6b, 0x1d, 0xa7, 0xd7, 0x3c, 0x68, 0xfc, 0xfc, 0x71, 0xaf, 0xf6,
	0xde, 0x29, 0x44, 0x27, 0x34, 0x01, 0xbc, 0x85, 0x50, 0x42, 0xc3, 0x73, 0x2e, 0x60, 0xc8, 0x99,
	0x5b, 0xcf, 0x55, 0x41, 0xd3, 0x22, 0x47, 0x0c, 0x6f, 0xa0, 0x06, 0xe3, 0x99, 0x56, 0xd2, 0xfd,
	0xcf, 0x50, 0x36, 0xc2, 0xdb, 0xa8, 0x15, 0x4a, 0xa1, 0x29, 0x17, 0xa0, 0xf2, 0xc4, 0x25, 0xc3,
	0xae, 0xce, 0xb1, 0x23, 0x86, 0xef, 0xa3, 0xb5, 0x85, 0xc4, 0x78, 0x68, 0x18, 0x51, 0x7b, 0x8e,
	0x96, 0x06, 0x8c, 0xcb, 0x44, 0x32, 0x88, 0xdd, 0xe5, 0xc2, 0x40, 0x8e, 0xbc, 0xc9, 0x81, 0xdc,
	0x80, 0x82, 0x88, 0x4b, 0xe1, 0xae, 0x14, 0x06, 0x8a, 0x08, 0xaf, 0xa1, 0x1a, 0xbd, 0x74, 0x9b,
	0x06, 0xab, 0xd1, 0x4b, 0x7c, 0x8a, 0xda, 0xe1, 0x38, 0xd3, 0x32, 0x19, 0xc6, 0x74, 0x04, 0x71,
	0xe6, 0xa2, 0x4e, 0xbd, 0xb7, 0xba, 0xb7, 0xeb, 0x2f, 0x6e, 0x83, 0x7f, 0x6d, 0xa0, 0xfe, 0x6b,
	0x23, 0x3f, 0x36, 0xea, 0x43, 0xa1, 0xd5, 0x2c, 0x68, 0x85, 0x15, 0xc8, 0x7b, 0x85, 0x6e, 0xdd,
	0x90, 0xe0, 0xff, 0x51, 0xfd, 0x02, 0x66, 0x66, 0xfe, 0xcd, 0x20, 0x3f, 0xe2, 0x75, 0xb4, 0x34,
	0xa1, 0xf1, 0xd8, 0x4e, 0x38, 0x28, 0x82, 0xfd, 0xda, 0x0b, 0xa7, 0x7b, 0xe5, 0xa0, 0xf5, 0x00,
	0x12, 0x39, 0x81, 0xb7, 0xc5, 0x55, 0x08, 0xe0, 0xe3, 0x18, 0x32, 0x8d, 0x5f, 0xa2, 0x96, 0xbd,
	0x1c, 0xd5, 0x6d, 0x6e, 0x54, 0xb6, 0x69, 0x13, 0xcc, 0x42, 0x57, 0xb3, 0x45, 0x90, 0x4f, 0xab,
	0x4c, 0xe5, 0xcc, 0xb6, 0x6c, 0x5a, 0xe4, 0x88, 0xe5, 0x6b, 0x29, 0x69, 0x33, 0xf1, 0x62, 0x9f,
	0x65, 0x85, 0x7c, 0xde, 0xdd, 0x4d, 0x74, 0xe7, 0x0f, 0x53, 0x59, 0x2a, 0x45, 0x06, 0x7b, 0x5f,
	0x1c, 0xb4, 0x6c, 0x31, 0xfc, 0x09, 0xb5, 0xaf, 0x89, 0x70, 0xa7, 0x3a, 0xc7, 0xbf, 0xfd, 0x29,
	0x6f, 0xfb, 0x1f, 0x8a, 0xa2, 0x43, 0xb7, 0xf7, 0xf9, 0xfb, 0xaf, 0xab, 0x5a, 0xb7, 0xbb, 0x45,
	0x26, 0x4f, 0xc9, 0x42, 0x4d, 0xac, 0x8e, 0x14, 0x59, 0xfb, 0xce, 0xa3, 0x83, 0xe1, 0xd7, 0xc1,
	0x49, 0x70, 0x8c, 0x96, 0x19, 0x9c, 0xd1, 0x71, 0xac, 0xf1, 0x00, 0xe1, 0x81, 0xe8, 0x80, 0x52,
	0x52, 0x75, 0x94, 0x2d, 0xe7, 0xe3, 0x5d, 0xf4, 0xd0, 0x7b, 0xb0, 0x43, 0x18, 0x9c, 0x71, 0xc1,
	0x8b, 0x97, 0x54, 0xfd, 0x2e, 0x1c, 0xe6, 0xf2, 0xb2, 0xf9, 0x87, 0x56, 0x95, 0x1a, 0x35, 0xcc,
	0x33, 0x7b, 0xf6, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xf6, 0xc4, 0xfa, 0x3c, 0x49, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServiceClient interface {
	// RemoveService removes Service with Agents.
	RemoveService(ctx context.Context, in *RemoveServiceRequest, opts ...grpc.CallOption) (*RemoveServiceResponse, error)
}

type serviceClient struct {
	cc *grpc.ClientConn
}

func NewServiceClient(cc *grpc.ClientConn) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) RemoveService(ctx context.Context, in *RemoveServiceRequest, opts ...grpc.CallOption) (*RemoveServiceResponse, error) {
	out := new(RemoveServiceResponse)
	err := c.cc.Invoke(ctx, "/management.Service/RemoveService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
type ServiceServer interface {
	// RemoveService removes Service with Agents.
	RemoveService(context.Context, *RemoveServiceRequest) (*RemoveServiceResponse, error)
}

// UnimplementedServiceServer can be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (*UnimplementedServiceServer) RemoveService(ctx context.Context, req *RemoveServiceRequest) (*RemoveServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveService not implemented")
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_RemoveService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).RemoveService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/management.Service/RemoveService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).RemoveService(ctx, req.(*RemoveServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "management.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RemoveService",
			Handler:    _Service_RemoveService_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "managementpb/service.proto",
}
