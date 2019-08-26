// Code generated by protoc-gen-go. DO NOT EDIT.
// source: managementpb/mysql.proto

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

type AddMySQLRequest struct {
	// Node identifier on which a service is been running. Required.
	// Use only one of these paramse (node_id, node_name or add_node)
	NodeId   string         `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	NodeName string         `protobuf:"bytes,21,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	AddNode  *AddNodeParams `protobuf:"bytes,22,opt,name=add_node,json=addNode,proto3" json:"add_node,omitempty"`
	// Unique across all Services user-defined name. Required.
	ServiceName string `protobuf:"bytes,2,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	// Node and Service access address (DNS name or IP). Required.
	Address string `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	// Service Access port. Required.
	Port uint32 `protobuf:"varint,4,opt,name=port,proto3" json:"port,omitempty"`
	// The "pmm-agent" identifier which should run agents. Required.
	PmmAgentId string `protobuf:"bytes,5,opt,name=pmm_agent_id,json=pmmAgentId,proto3" json:"pmm_agent_id,omitempty"`
	// Environment name.
	Environment string `protobuf:"bytes,6,opt,name=environment,proto3" json:"environment,omitempty"`
	// Cluster name.
	Cluster string `protobuf:"bytes,7,opt,name=cluster,proto3" json:"cluster,omitempty"`
	// Replication set name.
	ReplicationSet string `protobuf:"bytes,8,opt,name=replication_set,json=replicationSet,proto3" json:"replication_set,omitempty"`
	// MySQL username for scraping metrics.
	Username string `protobuf:"bytes,9,opt,name=username,proto3" json:"username,omitempty"`
	// MySQL password for scraping metrics.
	Password string `protobuf:"bytes,10,opt,name=password,proto3" json:"password,omitempty"`
	// If true, adds qan-mysql-perfschema-agent for provided service.
	QanMysqlPerfschema bool `protobuf:"varint,14,opt,name=qan_mysql_perfschema,json=qanMysqlPerfschema,proto3" json:"qan_mysql_perfschema,omitempty"`
	// If true, adds qan-mysql-slowlog-agent for provided service.
	QanMysqlSlowlog bool `protobuf:"varint,15,opt,name=qan_mysql_slowlog,json=qanMysqlSlowlog,proto3" json:"qan_mysql_slowlog,omitempty"`
	// Custom user-assigned labels.
	CustomLabels map[string]string `protobuf:"bytes,20,rep,name=custom_labels,json=customLabels,proto3" json:"custom_labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Skip connection check.
	SkipConnectionCheck  bool     `protobuf:"varint,30,opt,name=skip_connection_check,json=skipConnectionCheck,proto3" json:"skip_connection_check,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddMySQLRequest) Reset()         { *m = AddMySQLRequest{} }
func (m *AddMySQLRequest) String() string { return proto.CompactTextString(m) }
func (*AddMySQLRequest) ProtoMessage()    {}
func (*AddMySQLRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab81470951176953, []int{0}
}

func (m *AddMySQLRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddMySQLRequest.Unmarshal(m, b)
}
func (m *AddMySQLRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddMySQLRequest.Marshal(b, m, deterministic)
}
func (m *AddMySQLRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddMySQLRequest.Merge(m, src)
}
func (m *AddMySQLRequest) XXX_Size() int {
	return xxx_messageInfo_AddMySQLRequest.Size(m)
}
func (m *AddMySQLRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddMySQLRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddMySQLRequest proto.InternalMessageInfo

func (m *AddMySQLRequest) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *AddMySQLRequest) GetNodeName() string {
	if m != nil {
		return m.NodeName
	}
	return ""
}

func (m *AddMySQLRequest) GetAddNode() *AddNodeParams {
	if m != nil {
		return m.AddNode
	}
	return nil
}

func (m *AddMySQLRequest) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *AddMySQLRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *AddMySQLRequest) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *AddMySQLRequest) GetPmmAgentId() string {
	if m != nil {
		return m.PmmAgentId
	}
	return ""
}

func (m *AddMySQLRequest) GetEnvironment() string {
	if m != nil {
		return m.Environment
	}
	return ""
}

func (m *AddMySQLRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *AddMySQLRequest) GetReplicationSet() string {
	if m != nil {
		return m.ReplicationSet
	}
	return ""
}

func (m *AddMySQLRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *AddMySQLRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *AddMySQLRequest) GetQanMysqlPerfschema() bool {
	if m != nil {
		return m.QanMysqlPerfschema
	}
	return false
}

func (m *AddMySQLRequest) GetQanMysqlSlowlog() bool {
	if m != nil {
		return m.QanMysqlSlowlog
	}
	return false
}

func (m *AddMySQLRequest) GetCustomLabels() map[string]string {
	if m != nil {
		return m.CustomLabels
	}
	return nil
}

func (m *AddMySQLRequest) GetSkipConnectionCheck() bool {
	if m != nil {
		return m.SkipConnectionCheck
	}
	return false
}

// AddNodeParams is a params to add new node to inventory while adding new service.
type AddNodeParams struct {
	// Node type to be registered.
	NodeType inventorypb.NodeType `protobuf:"varint,1,opt,name=node_type,json=nodeType,proto3,enum=inventory.NodeType" json:"node_type,omitempty"`
	// Unique across all Nodes user-defined name. Can't be changed.
	NodeName string `protobuf:"bytes,2,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	// Linux machine-id.
	// Must be unique across all Generic Nodes if specified.
	MachineId string `protobuf:"bytes,3,opt,name=machine_id,json=machineId,proto3" json:"machine_id,omitempty"`
	// Linux distribution name and version.
	Distro string `protobuf:"bytes,4,opt,name=distro,proto3" json:"distro,omitempty"`
	// Container identifier. If specified, must be a unique Docker container identifier.
	ContainerId string `protobuf:"bytes,6,opt,name=container_id,json=containerId,proto3" json:"container_id,omitempty"`
	// Container name.
	ContainerName string `protobuf:"bytes,7,opt,name=container_name,json=containerName,proto3" json:"container_name,omitempty"`
	// Node model.
	NodeModel string `protobuf:"bytes,8,opt,name=node_model,json=nodeModel,proto3" json:"node_model,omitempty"`
	// Node region.
	Region string `protobuf:"bytes,9,opt,name=region,proto3" json:"region,omitempty"`
	// Node availability zone.
	Az string `protobuf:"bytes,10,opt,name=az,proto3" json:"az,omitempty"`
	// Custom user-assigned labels.
	CustomLabels         map[string]string `protobuf:"bytes,11,rep,name=custom_labels,json=customLabels,proto3" json:"custom_labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *AddNodeParams) Reset()         { *m = AddNodeParams{} }
func (m *AddNodeParams) String() string { return proto.CompactTextString(m) }
func (*AddNodeParams) ProtoMessage()    {}
func (*AddNodeParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab81470951176953, []int{1}
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

type AddMySQLResponse struct {
	Service              *inventorypb.MySQLService            `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	MysqldExporter       *inventorypb.MySQLdExporter          `protobuf:"bytes,2,opt,name=mysqld_exporter,json=mysqldExporter,proto3" json:"mysqld_exporter,omitempty"`
	QanMysqlPerfschema   *inventorypb.QANMySQLPerfSchemaAgent `protobuf:"bytes,3,opt,name=qan_mysql_perfschema,json=qanMysqlPerfschema,proto3" json:"qan_mysql_perfschema,omitempty"`
	QanMysqlSlowlog      *inventorypb.QANMySQLSlowlogAgent    `protobuf:"bytes,4,opt,name=qan_mysql_slowlog,json=qanMysqlSlowlog,proto3" json:"qan_mysql_slowlog,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                             `json:"-"`
	XXX_unrecognized     []byte                               `json:"-"`
	XXX_sizecache        int32                                `json:"-"`
}

func (m *AddMySQLResponse) Reset()         { *m = AddMySQLResponse{} }
func (m *AddMySQLResponse) String() string { return proto.CompactTextString(m) }
func (*AddMySQLResponse) ProtoMessage()    {}
func (*AddMySQLResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ab81470951176953, []int{2}
}

func (m *AddMySQLResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddMySQLResponse.Unmarshal(m, b)
}
func (m *AddMySQLResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddMySQLResponse.Marshal(b, m, deterministic)
}
func (m *AddMySQLResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddMySQLResponse.Merge(m, src)
}
func (m *AddMySQLResponse) XXX_Size() int {
	return xxx_messageInfo_AddMySQLResponse.Size(m)
}
func (m *AddMySQLResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddMySQLResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddMySQLResponse proto.InternalMessageInfo

func (m *AddMySQLResponse) GetService() *inventorypb.MySQLService {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *AddMySQLResponse) GetMysqldExporter() *inventorypb.MySQLdExporter {
	if m != nil {
		return m.MysqldExporter
	}
	return nil
}

func (m *AddMySQLResponse) GetQanMysqlPerfschema() *inventorypb.QANMySQLPerfSchemaAgent {
	if m != nil {
		return m.QanMysqlPerfschema
	}
	return nil
}

func (m *AddMySQLResponse) GetQanMysqlSlowlog() *inventorypb.QANMySQLSlowlogAgent {
	if m != nil {
		return m.QanMysqlSlowlog
	}
	return nil
}

func init() {
	proto.RegisterType((*AddMySQLRequest)(nil), "management.AddMySQLRequest")
	proto.RegisterMapType((map[string]string)(nil), "management.AddMySQLRequest.CustomLabelsEntry")
	proto.RegisterType((*AddNodeParams)(nil), "management.AddNodeParams")
	proto.RegisterMapType((map[string]string)(nil), "management.AddNodeParams.CustomLabelsEntry")
	proto.RegisterType((*AddMySQLResponse)(nil), "management.AddMySQLResponse")
}

func init() { proto.RegisterFile("managementpb/mysql.proto", fileDescriptor_ab81470951176953) }

var fileDescriptor_ab81470951176953 = []byte{
	// 940 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x55, 0xcd, 0x6e, 0xdb, 0x46,
	0x10, 0xae, 0xe4, 0x1f, 0xc9, 0x23, 0x5b, 0x76, 0x36, 0x4e, 0xcc, 0x2a, 0x49, 0xa3, 0xca, 0x28,
	0xa2, 0x24, 0xb5, 0xe8, 0xaa, 0x45, 0x51, 0xe4, 0x52, 0xc8, 0x86, 0x0f, 0x46, 0x6d, 0xc3, 0xa1,
	0x73, 0x28, 0x7a, 0x21, 0x56, 0xdc, 0x35, 0xcd, 0x9a, 0xdc, 0xa5, 0x77, 0x57, 0x52, 0x95, 0x63,
	0x1f, 0xa1, 0x7d, 0x8d, 0xbe, 0x44, 0x9f, 0xa1, 0xd7, 0x02, 0x05, 0x8a, 0x3e, 0x48, 0xb1, 0x43,
	0x52, 0xa2, 0x1c, 0xa7, 0x3d, 0xe4, 0x44, 0xce, 0x7c, 0xdf, 0x7e, 0x3b, 0x9c, 0xfd, 0x66, 0x09,
	0x4e, 0x42, 0x05, 0x0d, 0x79, 0xc2, 0x85, 0x49, 0x87, 0x6e, 0x32, 0xd5, 0x37, 0x71, 0x2f, 0x55,
	0xd2, 0x48, 0x02, 0x73, 0xa4, 0xf5, 0x75, 0x18, 0x99, 0xab, 0xd1, 0xb0, 0x17, 0xc8, 0xc4, 0x4d,
	0x26, 0x91, 0xb9, 0x96, 0x13, 0x37, 0x94, 0x7b, 0x48, 0xdc, 0x1b, 0xd3, 0x38, 0x62, 0xd4, 0x48,
	0xa5, 0xdd, 0xd9, 0x6b, 0xa6, 0xd1, 0x7a, 0x1c, 0x4a, 0x19, 0xc6, 0xdc, 0xa5, 0x69, 0xe4, 0x52,
	0x21, 0xa4, 0xa1, 0x26, 0x92, 0x42, 0xe7, 0xa8, 0x13, 0x89, 0x31, 0x17, 0x46, 0xaa, 0x69, 0x3a,
	0x74, 0x69, 0xc8, 0x85, 0x29, 0x90, 0x9d, 0x32, 0x22, 0x24, 0xe3, 0x05, 0xd0, 0x2a, 0x03, 0x9a,
	0xab, 0x71, 0x14, 0xcc, 0xb0, 0xcf, 0xf1, 0x11, 0xec, 0x85, 0x5c, 0xec, 0xe9, 0x09, 0x0d, 0x43,
	0xae, 0x5c, 0x99, 0xe2, 0x86, 0xef, 0x6e, 0xde, 0xf9, 0x73, 0x05, 0x36, 0x07, 0x8c, 0x9d, 0x4e,
	0x2f, 0x5e, 0x9f, 0x78, 0xfc, 0x66, 0xc4, 0xb5, 0x21, 0x3b, 0x50, 0xb3, 0x9b, 0xf9, 0x11, 0x73,
	0x2a, 0xed, 0x4a, 0x77, 0xcd, 0x5b, 0xb5, 0xe1, 0x31, 0x23, 0x8f, 0x60, 0x0d, 0x01, 0x41, 0x13,
	0xee, 0x3c, 0x40, 0xa8, 0x6e, 0x13, 0x67, 0x34, 0xe1, 0xe4, 0x2b, 0xa8, 0x53, 0xc6, 0x7c, 0x1b,
	0x3b, 0x0f, 0xdb, 0x95, 0x6e, 0xa3, 0xff, 0x71, 0x6f, 0xde, 0xbb, 0xde, 0x80, 0xb1, 0x33, 0xc9,
	0xf8, 0x39, 0x55, 0x34, 0xd1, 0x5e, 0x8d, 0x66, 0x21, 0x79, 0x0e, 0xeb, 0x79, 0xfd, 0x99, 0x6a,
	0xd5, 0xaa, 0x1e, 0xac, 0xfe, 0xfd, 0xd7, 0xd3, 0xea, 0xf7, 0x15, 0xaf, 0x91, 0x63, 0xb8, 0x41,
	0x1b, 0xec, 0x2a, 0xc5, 0xb5, 0x76, 0x96, 0x16, 0x58, 0x45, 0x9a, 0xb4, 0x60, 0x39, 0x95, 0xca,
	0x38, 0xcb, 0xed, 0x4a, 0x77, 0x23, 0x83, 0xb7, 0x3e, 0xf2, 0x30, 0x47, 0xba, 0xb0, 0x9e, 0x26,
	0x89, 0x8f, 0xfd, 0xb5, 0x5f, 0xb6, 0xb2, 0x20, 0x01, 0x69, 0x92, 0x0c, 0x2c, 0x74, 0xcc, 0x48,
	0x1b, 0x1a, 0x5c, 0x8c, 0x23, 0x25, 0x85, 0x2d, 0xdc, 0x59, 0xc5, 0xef, 0x2c, 0xa7, 0x88, 0x03,
	0xb5, 0x20, 0x1e, 0x69, 0xc3, 0x95, 0x53, 0x43, 0xb4, 0x08, 0xc9, 0x33, 0xd8, 0x54, 0x3c, 0x8d,
	0xa3, 0x00, 0x9b, 0xec, 0x6b, 0x6e, 0x9c, 0x3a, 0x32, 0x9a, 0xa5, 0xf4, 0x05, 0x37, 0xa4, 0x03,
	0xf5, 0x91, 0xe6, 0x0a, 0xbf, 0x79, 0x6d, 0xa1, 0x94, 0x59, 0x9e, 0xb4, 0xa0, 0x9e, 0x52, 0xad,
	0x27, 0x52, 0x31, 0x07, 0xb2, 0x6e, 0x17, 0x31, 0xd9, 0x87, 0xed, 0x1b, 0x2a, 0x7c, 0x74, 0xaa,
	0x9f, 0x72, 0x75, 0xa9, 0x83, 0x2b, 0x9e, 0x50, 0xa7, 0xd9, 0xae, 0x74, 0xeb, 0x1e, 0xb9, 0xa1,
	0xe2, 0xd4, 0x42, 0xe7, 0x33, 0x84, 0xbc, 0x80, 0x7b, 0xf3, 0x15, 0x3a, 0x96, 0x93, 0x58, 0x86,
	0xce, 0x26, 0xd2, 0x37, 0x0b, 0xfa, 0x45, 0x96, 0x26, 0x1e, 0x6c, 0x04, 0x23, 0x6d, 0x64, 0xe2,
	0xc7, 0x74, 0xc8, 0x63, 0xed, 0x6c, 0xb7, 0x97, 0xba, 0x8d, 0xfe, 0xde, 0xad, 0x03, 0x2d, 0xbb,
	0xa6, 0x77, 0x88, 0x0b, 0x4e, 0x90, 0x7f, 0x24, 0x8c, 0x9a, 0x7a, 0xeb, 0x41, 0x29, 0x45, 0xfa,
	0xf0, 0x40, 0x5f, 0x47, 0xa9, 0x1f, 0x48, 0x21, 0x78, 0x80, 0xed, 0x09, 0xae, 0x78, 0x70, 0xed,
	0x7c, 0x82, 0x35, 0xdc, 0xb7, 0xe0, 0xe1, 0x0c, 0x3b, 0xb4, 0x50, 0xeb, 0x5b, 0xb8, 0xf7, 0x8e,
	0x2c, 0xd9, 0x82, 0xa5, 0x6b, 0x3e, 0xcd, 0xad, 0x69, 0x5f, 0xc9, 0x36, 0xac, 0x8c, 0x69, 0x3c,
	0xca, 0xdd, 0xe3, 0x65, 0xc1, 0xab, 0xea, 0x37, 0x95, 0xce, 0xef, 0x4b, 0xb0, 0xb1, 0xe0, 0x3c,
	0xb2, 0x9f, 0x7b, 0xd8, 0x4c, 0x53, 0x8e, 0x1a, 0xcd, 0xfe, 0xfd, 0xde, 0x6c, 0x9c, 0x7a, 0x96,
	0xf9, 0x66, 0x9a, 0xf2, 0xcc, 0xd8, 0xf6, 0x8d, 0xec, 0x96, 0x5d, 0xbf, 0xe8, 0xcf, 0xb9, 0xfb,
	0x9f, 0x00, 0x24, 0x34, 0xb8, 0x8a, 0x04, 0x8e, 0x0d, 0xfa, 0xd3, 0x5b, 0xcb, 0x33, 0xc7, 0x8c,
	0x3c, 0x84, 0x55, 0x16, 0x69, 0xa3, 0x24, 0x7a, 0x73, 0xcd, 0xcb, 0x23, 0xf2, 0x29, 0xac, 0x07,
	0x52, 0x18, 0x1a, 0x09, 0xae, 0xec, 0xc2, 0xdc, 0x6c, 0xb3, 0xdc, 0x31, 0x23, 0x9f, 0x41, 0x73,
	0x4e, 0xc1, 0x1a, 0x32, 0xcf, 0x6d, 0xcc, 0xb2, 0x45, 0x01, 0x58, 0x65, 0x22, 0x19, 0x8f, 0x73,
	0xd3, 0x61, 0xdd, 0xa7, 0x36, 0x61, 0x0b, 0x50, 0x3c, 0x8c, 0xa4, 0xc8, 0xdc, 0xe6, 0xe5, 0x11,
	0x69, 0x42, 0x95, 0xbe, 0xcd, 0xdd, 0x55, 0xa5, 0x6f, 0xc9, 0xf9, 0xed, 0x93, 0x6f, 0xe0, 0xc9,
	0xbf, 0x7c, 0xef, 0x28, 0xff, 0xdf, 0xb9, 0x7f, 0xf8, 0x19, 0xfe, 0x56, 0x85, 0xad, 0xb9, 0xd9,
	0x74, 0x2a, 0x85, 0xe6, 0xe4, 0x0b, 0xa8, 0xe5, 0x77, 0x03, 0x8a, 0x34, 0xfa, 0x3b, 0xa5, 0x43,
	0x44, 0xea, 0x45, 0x06, 0x7b, 0x05, 0x8f, 0x1c, 0xc0, 0x26, 0x9a, 0x9f, 0xf9, 0xfc, 0x27, 0x7b,
	0x25, 0x70, 0x85, 0x7b, 0xd9, 0x7b, 0xea, 0xd6, 0x52, 0x76, 0x94, 0x13, 0xbc, 0x66, 0xb6, 0xa2,
	0x88, 0xc9, 0x9b, 0xf7, 0x8c, 0xdd, 0x12, 0x0a, 0x75, 0x4a, 0x42, 0xaf, 0x07, 0x67, 0xa8, 0x65,
	0x27, 0xf0, 0x02, 0x49, 0x78, 0xbf, 0xdc, 0x39, 0x9a, 0xdf, 0xdd, 0x35, 0x9a, 0xcb, 0x28, 0xf9,
	0xf4, 0x0e, 0xc9, 0x7c, 0x4a, 0x33, 0xbd, 0xdb, 0xb3, 0xdb, 0xd7, 0xb0, 0x82, 0x2c, 0xf2, 0x23,
	0xd4, 0x8b, 0xb6, 0x91, 0x47, 0xff, 0x31, 0xb9, 0xad, 0xc7, 0x77, 0x83, 0x59, 0xa7, 0x3b, 0xbb,
	0x3f, 0xff, 0xf1, 0xcf, 0xaf, 0xd5, 0x27, 0x1d, 0xc7, 0x1d, 0xef, 0xbb, 0x73, 0xa2, 0x8b, 0x2c,
	0x77, 0xc0, 0xd8, 0xab, 0xca, 0x8b, 0x03, 0xff, 0x97, 0xc1, 0x99, 0x77, 0x02, 0x35, 0xc6, 0x2f,
	0xe9, 0x28, 0x36, 0x64, 0x00, 0x64, 0x20, 0xda, 0x5c, 0x29, 0xa9, 0xda, 0x2a, 0x57, 0xea, 0x91,
	0x97, 0xf0, 0xbc, 0xf5, 0x6c, 0xd7, 0x65, 0xfc, 0x32, 0x12, 0x51, 0xf6, 0x47, 0x2a, 0xff, 0x77,
	0x8f, 0x2c, 0xbd, 0xd8, 0xf7, 0x87, 0xf5, 0x32, 0x34, 0x5c, 0xc5, 0xdf, 0xd5, 0x97, 0xff, 0x06,
	0x00, 0x00, 0xff, 0xff, 0xa7, 0x43, 0x52, 0xb4, 0xa9, 0x07, 0x00, 0x00,
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
	// AddMySQL adds MySQL Service and starts several Agents.
	// It automatically adds a service to inventory, which is running on provided "node_id",
	// then adds "mysqld_exporter", and "qan_mysql_perfschema" agents
	// with provided "pmm_agent_id" and other parameters.
	AddMySQL(ctx context.Context, in *AddMySQLRequest, opts ...grpc.CallOption) (*AddMySQLResponse, error)
}

type mySQLClient struct {
	cc *grpc.ClientConn
}

func NewMySQLClient(cc *grpc.ClientConn) MySQLClient {
	return &mySQLClient{cc}
}

func (c *mySQLClient) AddMySQL(ctx context.Context, in *AddMySQLRequest, opts ...grpc.CallOption) (*AddMySQLResponse, error) {
	out := new(AddMySQLResponse)
	err := c.cc.Invoke(ctx, "/management.MySQL/AddMySQL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MySQLServer is the server API for MySQL service.
type MySQLServer interface {
	// AddMySQL adds MySQL Service and starts several Agents.
	// It automatically adds a service to inventory, which is running on provided "node_id",
	// then adds "mysqld_exporter", and "qan_mysql_perfschema" agents
	// with provided "pmm_agent_id" and other parameters.
	AddMySQL(context.Context, *AddMySQLRequest) (*AddMySQLResponse, error)
}

// UnimplementedMySQLServer can be embedded to have forward compatible implementations.
type UnimplementedMySQLServer struct {
}

func (*UnimplementedMySQLServer) AddMySQL(ctx context.Context, req *AddMySQLRequest) (*AddMySQLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMySQL not implemented")
}

func RegisterMySQLServer(s *grpc.Server, srv MySQLServer) {
	s.RegisterService(&_MySQL_serviceDesc, srv)
}

func _MySQL_AddMySQL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMySQLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MySQLServer).AddMySQL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/management.MySQL/AddMySQL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MySQLServer).AddMySQL(ctx, req.(*AddMySQLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MySQL_serviceDesc = grpc.ServiceDesc{
	ServiceName: "management.MySQL",
	HandlerType: (*MySQLServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddMySQL",
			Handler:    _MySQL_AddMySQL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "managementpb/mysql.proto",
}
