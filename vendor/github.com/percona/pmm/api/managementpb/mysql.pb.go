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
	// Node identifier on which a service is been running.
	// Exactly one of these parameters should be present: node_id, node_name, add_node.
	NodeId string `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	// Node name on which a service is been running.
	// Exactly one of these parameters should be present: node_id, node_name, add_node.
	NodeName string `protobuf:"bytes,21,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	// Create a new Node with those parameters.
	// Exactly one of these parameters should be present: node_id, node_name, add_node.
	AddNode *AddNodeParams `protobuf:"bytes,22,opt,name=add_node,json=addNode,proto3" json:"add_node,omitempty"`
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
	SkipConnectionCheck bool `protobuf:"varint,30,opt,name=skip_connection_check,json=skipConnectionCheck,proto3" json:"skip_connection_check,omitempty"`
	// Disable query examples.
	DisableQueryExamples bool `protobuf:"varint,31,opt,name=disable_query_examples,json=disableQueryExamples,proto3" json:"disable_query_examples,omitempty"`
	// If qan-mysql-slowlog-agent is added, slowlog file is rotated at this size if > 0.
	// If zero, default value 1GB is used. Use negative value to disable rotation.
	MaxSlowlogFileSize int64 `protobuf:"varint,32,opt,name=max_slowlog_file_size,json=maxSlowlogFileSize,proto3" json:"max_slowlog_file_size,omitempty"`
	// Use TLS for database connections.
	Tls bool `protobuf:"varint,41,opt,name=tls,proto3" json:"tls,omitempty"`
	// Skip TLS certificate and hostname validation.
	TlsSkipVerify        bool     `protobuf:"varint,42,opt,name=tls_skip_verify,json=tlsSkipVerify,proto3" json:"tls_skip_verify,omitempty"`
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

func (m *AddMySQLRequest) GetDisableQueryExamples() bool {
	if m != nil {
		return m.DisableQueryExamples
	}
	return false
}

func (m *AddMySQLRequest) GetMaxSlowlogFileSize() int64 {
	if m != nil {
		return m.MaxSlowlogFileSize
	}
	return 0
}

func (m *AddMySQLRequest) GetTls() bool {
	if m != nil {
		return m.Tls
	}
	return false
}

func (m *AddMySQLRequest) GetTlsSkipVerify() bool {
	if m != nil {
		return m.TlsSkipVerify
	}
	return false
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
	return fileDescriptor_ab81470951176953, []int{1}
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
	proto.RegisterType((*AddMySQLResponse)(nil), "management.AddMySQLResponse")
}

func init() { proto.RegisterFile("managementpb/mysql.proto", fileDescriptor_ab81470951176953) }

var fileDescriptor_ab81470951176953 = []byte{
	// 885 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x94, 0xd1, 0x6e, 0xdb, 0x36,
	0x14, 0x86, 0xe7, 0x38, 0xb1, 0x1d, 0x3a, 0x89, 0x53, 0x2e, 0x69, 0x39, 0xb7, 0x5b, 0x04, 0x17,
	0x58, 0x9d, 0x6c, 0xb6, 0x5a, 0xaf, 0x18, 0x86, 0xde, 0x0c, 0x4e, 0x90, 0x01, 0xc5, 0xd2, 0xa0,
	0x91, 0x87, 0x61, 0xd8, 0x8d, 0x40, 0x8b, 0xc7, 0x0a, 0x67, 0x89, 0x94, 0x49, 0xda, 0x8e, 0x7b,
	0xb9, 0x47, 0xd8, 0x5e, 0x63, 0x6f, 0xb3, 0x07, 0x18, 0x30, 0xec, 0x3d, 0x36, 0x88, 0x92, 0x6c,
	0x27, 0xcd, 0x7a, 0x25, 0x9d, 0xf3, 0xfd, 0xfc, 0x79, 0x78, 0x74, 0x44, 0x44, 0x62, 0x2a, 0x68,
	0x08, 0x31, 0x08, 0x93, 0x0c, 0xdd, 0x78, 0xa1, 0x27, 0x51, 0x37, 0x51, 0xd2, 0x48, 0x8c, 0x56,
	0xa4, 0xf9, 0x75, 0xc8, 0xcd, 0xf5, 0x74, 0xd8, 0x0d, 0x64, 0xec, 0xc6, 0x73, 0x6e, 0xc6, 0x72,
	0xee, 0x86, 0xb2, 0x63, 0x85, 0x9d, 0x19, 0x8d, 0x38, 0xa3, 0x46, 0x2a, 0xed, 0x2e, 0x5f, 0x33,
	0x8f, 0xe6, 0x93, 0x50, 0xca, 0x30, 0x02, 0x97, 0x26, 0xdc, 0xa5, 0x42, 0x48, 0x43, 0x0d, 0x97,
	0x42, 0xe7, 0x94, 0x70, 0x31, 0x03, 0x61, 0xa4, 0x5a, 0x24, 0x43, 0x97, 0x86, 0x20, 0x4c, 0x41,
	0x9a, 0xeb, 0x44, 0x83, 0x9a, 0xf1, 0x00, 0x96, 0xec, 0x56, 0xc5, 0x39, 0xcc, 0xd9, 0x97, 0xf6,
	0x11, 0x74, 0x42, 0x10, 0x1d, 0x3d, 0xa7, 0x61, 0x08, 0xca, 0x95, 0x89, 0xdd, 0xf3, 0xfd, 0xfd,
	0x5b, 0xff, 0x56, 0x50, 0xa3, 0xcf, 0xd8, 0x9b, 0xc5, 0xe0, 0xea, 0xc2, 0x83, 0xc9, 0x14, 0xb4,
	0xc1, 0x8f, 0x50, 0x55, 0x48, 0x06, 0x3e, 0x67, 0xa4, 0xe4, 0x94, 0xda, 0xdb, 0x5e, 0x25, 0x0d,
	0x5f, 0x33, 0xfc, 0x18, 0x6d, 0x5b, 0x20, 0x68, 0x0c, 0xe4, 0xd0, 0xa2, 0x5a, 0x9a, 0xb8, 0xa4,
	0x31, 0xe0, 0x97, 0xa8, 0x46, 0x19, 0xf3, 0xd3, 0x98, 0x3c, 0x74, 0x4a, 0xed, 0x7a, 0xef, 0x93,
	0xee, 0xaa, 0xcc, 0x6e, 0x9f, 0xb1, 0x4b, 0xc9, 0xe0, 0x2d, 0x55, 0x34, 0xd6, 0x5e, 0x95, 0x66,
	0x21, 0x3e, 0x46, 0x3b, 0x79, 0xf9, 0x99, 0xeb, 0x46, 0xea, 0x7a, 0x5a, 0xf9, 0xfb, 0xaf, 0xa3,
	0x8d, 0x9f, 0x4a, 0x5e, 0x3d, 0x67, 0x76, 0x03, 0x07, 0xa5, 0xab, 0x14, 0x68, 0x4d, 0xca, 0xb7,
	0x54, 0x45, 0x1a, 0x37, 0xd1, 0x66, 0x22, 0x95, 0x21, 0x9b, 0x4e, 0xa9, 0xbd, 0x9b, 0xe1, 0xfd,
	0x8f, 0x3c, 0x9b, 0xc3, 0x6d, 0xb4, 0x93, 0xc4, 0xb1, 0x6f, 0x5b, 0x9c, 0x9e, 0x6c, 0xeb, 0x96,
	0x05, 0x4a, 0xe2, 0xb8, 0x9f, 0xa2, 0xd7, 0x0c, 0x3b, 0xa8, 0x0e, 0x62, 0xc6, 0x95, 0x14, 0x69,
	0xe1, 0xa4, 0x62, 0xcf, 0xb9, 0x9e, 0xc2, 0x04, 0x55, 0x83, 0x68, 0xaa, 0x0d, 0x28, 0x52, 0xb5,
	0xb4, 0x08, 0xf1, 0x33, 0xd4, 0x50, 0x90, 0x44, 0x3c, 0xb0, 0x4d, 0xf6, 0x35, 0x18, 0x52, 0xb3,
	0x8a, 0xbd, 0xb5, 0xf4, 0x00, 0x0c, 0x6e, 0xa1, 0xda, 0x54, 0x83, 0xb2, 0x67, 0xde, 0xbe, 0x55,
	0xca, 0x32, 0x8f, 0x9b, 0xa8, 0x96, 0x50, 0xad, 0xe7, 0x52, 0x31, 0x82, 0xb2, 0x6e, 0x17, 0x31,
	0x7e, 0x8e, 0x0e, 0x26, 0x54, 0xf8, 0x76, 0x58, 0xfd, 0x04, 0xd4, 0x48, 0x07, 0xd7, 0x10, 0x53,
	0xb2, 0xe7, 0x94, 0xda, 0x35, 0x0f, 0x4f, 0xa8, 0x78, 0x93, 0xa2, 0xb7, 0x4b, 0x82, 0x4f, 0xd0,
	0x83, 0xd5, 0x0a, 0x1d, 0xc9, 0x79, 0x24, 0x43, 0xd2, 0xb0, 0xf2, 0x46, 0x21, 0x1f, 0x64, 0x69,
	0xec, 0xa1, 0xdd, 0x60, 0xaa, 0x8d, 0x8c, 0xfd, 0x88, 0x0e, 0x21, 0xd2, 0xe4, 0xc0, 0x29, 0xb7,
	0xeb, 0xbd, 0xce, 0x9d, 0x0f, 0xba, 0x3e, 0x35, 0xdd, 0x33, 0xbb, 0xe0, 0xc2, 0xea, 0xcf, 0x85,
	0x51, 0x0b, 0x6f, 0x27, 0x58, 0x4b, 0xe1, 0x1e, 0x3a, 0xd4, 0x63, 0x9e, 0xf8, 0x81, 0x14, 0x02,
	0x02, 0xdb, 0x9e, 0xe0, 0x1a, 0x82, 0x31, 0xf9, 0xcc, 0xd6, 0xf0, 0x71, 0x0a, 0xcf, 0x96, 0xec,
	0x2c, 0x45, 0xf8, 0x25, 0x7a, 0xc8, 0xb8, 0xa6, 0xc3, 0x08, 0xfc, 0xc9, 0x14, 0xd4, 0xc2, 0x87,
	0x1b, 0x1a, 0x27, 0x11, 0x68, 0x72, 0x64, 0x17, 0x1d, 0xe4, 0xf4, 0x2a, 0x85, 0xe7, 0x39, 0xc3,
	0x2f, 0xd0, 0x61, 0x4c, 0x6f, 0x8a, 0x33, 0xfa, 0x23, 0x1e, 0x81, 0xaf, 0xf9, 0x3b, 0x20, 0x8e,
	0x53, 0x6a, 0x97, 0x3d, 0x1c, 0xd3, 0x9b, 0xfc, 0xa0, 0xdf, 0xf1, 0x08, 0x06, 0xfc, 0x1d, 0xe0,
	0x7d, 0x54, 0x36, 0x91, 0x26, 0xc7, 0xd6, 0x35, 0x7d, 0xc5, 0x9f, 0xa3, 0x86, 0x89, 0xb4, 0x6f,
	0x4b, 0x9e, 0x81, 0xe2, 0xa3, 0x05, 0x39, 0xb1, 0x74, 0xd7, 0x44, 0x7a, 0x30, 0xe6, 0xc9, 0x8f,
	0x36, 0xd9, 0xfc, 0x16, 0x3d, 0x78, 0xef, 0xe4, 0xa9, 0xdd, 0x18, 0x16, 0xf9, 0xdf, 0x93, 0xbe,
	0xe2, 0x03, 0xb4, 0x35, 0xa3, 0xd1, 0x34, 0x1f, 0x70, 0x2f, 0x0b, 0x5e, 0x6d, 0x7c, 0x53, 0x6a,
	0xfd, 0xb1, 0x81, 0xf6, 0x57, 0xbd, 0xd4, 0x89, 0x14, 0x1a, 0xf0, 0x0b, 0x54, 0xcd, 0x47, 0xdf,
	0x9a, 0xd4, 0x7b, 0x8f, 0xba, 0xcb, 0xeb, 0xa0, 0x6b, 0xa5, 0x83, 0x0c, 0x7b, 0x85, 0x0e, 0x9f,
	0xa2, 0x86, 0xfd, 0xb6, 0xcc, 0x87, 0x9b, 0x74, 0xe2, 0x41, 0xd9, 0xbd, 0xd2, 0xdf, 0xf0, 0xce,
	0x52, 0x76, 0x9e, 0x0b, 0xbc, 0xbd, 0x6c, 0x45, 0x11, 0xe3, 0x1f, 0xfe, 0x67, 0xaa, 0xca, 0xd6,
	0xa8, 0xb5, 0x66, 0x74, 0xd5, 0xbf, 0xb4, 0x5e, 0xe9, 0x80, 0x0d, 0xac, 0xc8, 0xfe, 0x3e, 0xf7,
	0x4e, 0xde, 0xf7, 0xf7, 0x4d, 0xde, 0xa6, 0xb5, 0x3c, 0xba, 0xc7, 0x32, 0xff, 0x36, 0x99, 0xdf,
	0xdd, 0xd1, 0xec, 0x69, 0xb4, 0x65, 0x55, 0xf8, 0x17, 0x54, 0x2b, 0xda, 0x86, 0x1f, 0x7f, 0x60,
	0x30, 0x9b, 0x4f, 0xee, 0x87, 0x59, 0xa7, 0x5b, 0x4f, 0x7f, 0xfd, 0xf3, 0x9f, 0xdf, 0x37, 0x3e,
	0x6d, 0x11, 0x77, 0xf6, 0xdc, 0x5d, 0x09, 0x5d, 0xab, 0x72, 0xfb, 0x8c, 0xbd, 0x2a, 0x9d, 0x9c,
	0xfa, 0xbf, 0xf5, 0x2f, 0xbd, 0x0b, 0x54, 0x65, 0x30, 0xa2, 0xd3, 0xc8, 0xe0, 0x3e, 0xc2, 0x7d,
	0xe1, 0x80, 0x52, 0x52, 0x39, 0x2a, 0x77, 0xea, 0xe2, 0x2f, 0xd0, 0x71, 0xf3, 0xd9, 0x53, 0x97,
	0xc1, 0x88, 0x0b, 0x9e, 0x5d, 0xb8, 0xeb, 0xf7, 0xf4, 0x79, 0x2a, 0x2f, 0xf6, 0xfd, 0x79, 0x67,
	0x1d, 0x0d, 0x2b, 0xf6, 0x36, 0xfe, 0xea, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4f, 0xb5, 0x73,
	0x39, 0x8b, 0x06, 0x00, 0x00,
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
