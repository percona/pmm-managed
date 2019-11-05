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
	NodeName string `protobuf:"bytes,2,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	// Create a new Node with those parameters.
	// Exactly one of these parameters should be present: node_id, node_name, add_node.
	AddNode *AddNodeParams `protobuf:"bytes,3,opt,name=add_node,json=addNode,proto3" json:"add_node,omitempty"`
	// Unique across all Services user-defined name. Required.
	ServiceName string `protobuf:"bytes,4,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	// Node and Service access address (DNS name or IP). Required.
	Address string `protobuf:"bytes,5,opt,name=address,proto3" json:"address,omitempty"`
	// Service Access port. Required.
	Port uint32 `protobuf:"varint,6,opt,name=port,proto3" json:"port,omitempty"`
	// The "pmm-agent" identifier which should run agents. Required.
	PmmAgentId string `protobuf:"bytes,7,opt,name=pmm_agent_id,json=pmmAgentId,proto3" json:"pmm_agent_id,omitempty"`
	// Environment name.
	Environment string `protobuf:"bytes,8,opt,name=environment,proto3" json:"environment,omitempty"`
	// Cluster name.
	Cluster string `protobuf:"bytes,9,opt,name=cluster,proto3" json:"cluster,omitempty"`
	// Replication set name.
	ReplicationSet string `protobuf:"bytes,10,opt,name=replication_set,json=replicationSet,proto3" json:"replication_set,omitempty"`
	// MySQL username for scraping metrics.
	Username string `protobuf:"bytes,11,opt,name=username,proto3" json:"username,omitempty"`
	// MySQL password for scraping metrics.
	Password string `protobuf:"bytes,12,opt,name=password,proto3" json:"password,omitempty"`
	// If true, adds qan-mysql-perfschema-agent for provided service.
	QanMysqlPerfschema bool `protobuf:"varint,13,opt,name=qan_mysql_perfschema,json=qanMysqlPerfschema,proto3" json:"qan_mysql_perfschema,omitempty"`
	// If true, adds qan-mysql-slowlog-agent for provided service.
	QanMysqlSlowlog bool `protobuf:"varint,14,opt,name=qan_mysql_slowlog,json=qanMysqlSlowlog,proto3" json:"qan_mysql_slowlog,omitempty"`
	// Custom user-assigned labels.
	CustomLabels map[string]string `protobuf:"bytes,15,rep,name=custom_labels,json=customLabels,proto3" json:"custom_labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Skip connection check.
	SkipConnectionCheck bool `protobuf:"varint,16,opt,name=skip_connection_check,json=skipConnectionCheck,proto3" json:"skip_connection_check,omitempty"`
	// Disable query examples.
	DisableQueryExamples bool `protobuf:"varint,17,opt,name=disable_query_examples,json=disableQueryExamples,proto3" json:"disable_query_examples,omitempty"`
	// If qan-mysql-slowlog-agent is added, slowlog file is rotated at this size if > 0.
	// If zero, default value 1GB is used. Use negative value to disable rotation.
	MaxSlowlogFileSize int64 `protobuf:"varint,18,opt,name=max_slowlog_file_size,json=maxSlowlogFileSize,proto3" json:"max_slowlog_file_size,omitempty"`
	// Use TLS for database connections.
	Tls bool `protobuf:"varint,19,opt,name=tls,proto3" json:"tls,omitempty"`
	// Skip TLS certificate and hostname validation.
	TlsSkipVerify bool `protobuf:"varint,20,opt,name=tls_skip_verify,json=tlsSkipVerify,proto3" json:"tls_skip_verify,omitempty"`
	// Max number of tables allowed for a heavy options.
	MaxNumberOfTables    int32    `protobuf:"varint,21,opt,name=max_number_of_tables,json=maxNumberOfTables,proto3" json:"max_number_of_tables,omitempty"`
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

func (m *AddMySQLRequest) GetMaxNumberOfTables() int32 {
	if m != nil {
		return m.MaxNumberOfTables
	}
	return 0
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
	// 917 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x94, 0xdb, 0x6e, 0x1b, 0x37,
	0x13, 0xc7, 0xbf, 0xf5, 0x49, 0x32, 0x65, 0x5b, 0x36, 0x63, 0x7f, 0x61, 0x95, 0x14, 0x59, 0x28,
	0x40, 0xa3, 0xa4, 0xb5, 0xb6, 0x76, 0x83, 0xa2, 0xc8, 0x4d, 0x21, 0x1b, 0x2e, 0x10, 0xd4, 0x71,
	0xe3, 0x55, 0x50, 0x14, 0xbd, 0x59, 0x50, 0xcb, 0xd1, 0x9a, 0xf5, 0x2e, 0xb9, 0x26, 0xa9, 0x53,
	0x2e, 0xfb, 0x08, 0xed, 0x6b, 0xf4, 0xa2, 0xef, 0xd2, 0x07, 0x28, 0x50, 0xf4, 0x41, 0x0a, 0x72,
	0x77, 0x25, 0xd9, 0x71, 0x7b, 0xa5, 0x9d, 0xf9, 0xfd, 0x39, 0x33, 0x1c, 0xcd, 0x10, 0x91, 0x8c,
	0x0a, 0x9a, 0x40, 0x06, 0xc2, 0xe4, 0x83, 0x20, 0x9b, 0xe9, 0x9b, 0xb4, 0x9b, 0x2b, 0x69, 0x24,
	0x46, 0x0b, 0xd2, 0xfa, 0x32, 0xe1, 0xe6, 0x6a, 0x34, 0xe8, 0xc6, 0x32, 0x0b, 0xb2, 0x09, 0x37,
	0xd7, 0x72, 0x12, 0x24, 0xf2, 0xd0, 0x09, 0x0f, 0xc7, 0x34, 0xe5, 0x8c, 0x1a, 0xa9, 0x74, 0x30,
	0xff, 0x2c, 0x62, 0xb4, 0x1e, 0x27, 0x52, 0x26, 0x29, 0x04, 0x34, 0xe7, 0x01, 0x15, 0x42, 0x1a,
	0x6a, 0xb8, 0x14, 0xba, 0xa4, 0x84, 0x8b, 0x31, 0x08, 0x23, 0xd5, 0x2c, 0x1f, 0x04, 0x34, 0x01,
	0x61, 0x2a, 0xd2, 0x5a, 0x26, 0x1a, 0xd4, 0x98, 0xc7, 0x30, 0x67, 0xb7, 0x2a, 0x2e, 0x61, 0xc9,
	0x3e, 0x73, 0x3f, 0xf1, 0x61, 0x02, 0xe2, 0x50, 0x4f, 0x68, 0x92, 0x80, 0x0a, 0x64, 0xee, 0x72,
	0x7e, 0x98, 0xbf, 0xfd, 0x7b, 0x0d, 0x35, 0x7b, 0x8c, 0xbd, 0x99, 0xf5, 0x2f, 0xcf, 0x43, 0xb8,
	0x19, 0x81, 0x36, 0xf8, 0x21, 0xaa, 0x09, 0xc9, 0x20, 0xe2, 0x8c, 0x78, 0xbe, 0xd7, 0xd9, 0x0c,
	0x37, 0xac, 0xf9, 0x9a, 0xe1, 0x47, 0x68, 0xd3, 0x01, 0x41, 0x33, 0x20, 0x2b, 0x0e, 0xd5, 0xad,
	0xe3, 0x82, 0x66, 0x80, 0x5f, 0xa2, 0x3a, 0x65, 0x2c, 0xb2, 0x36, 0x59, 0xf5, 0xbd, 0x4e, 0xe3,
	0xf8, 0xa3, 0xee, 0xa2, 0xcc, 0x6e, 0x8f, 0xb1, 0x0b, 0xc9, 0xe0, 0x2d, 0x55, 0x34, 0xd3, 0x61,
	0x8d, 0x16, 0x26, 0x7e, 0x8e, 0xb6, 0xca, 0xf2, 0x8b, 0xa8, 0x6b, 0x36, 0xea, 0xc9, 0xc6, 0x5f,
	0x7f, 0x3e, 0x59, 0xf9, 0xc1, 0x0b, 0x1b, 0x25, 0x73, 0x09, 0x7c, 0x64, 0x4f, 0x29, 0xd0, 0x9a,
	0xac, 0xdf, 0x52, 0x55, 0x6e, 0xdc, 0x42, 0x6b, 0xb9, 0x54, 0x86, 0x6c, 0xf8, 0x5e, 0x67, 0xbb,
	0xc0, 0xbb, 0xff, 0x0b, 0x9d, 0x0f, 0x77, 0xd0, 0x56, 0x9e, 0x65, 0x91, 0x6b, 0xb1, 0xbd, 0x59,
	0xed, 0x56, 0x08, 0x94, 0x67, 0x59, 0xcf, 0xa2, 0xd7, 0x0c, 0xfb, 0xa8, 0x01, 0x62, 0xcc, 0x95,
	0x14, 0xb6, 0x70, 0x52, 0x77, 0xf7, 0x5c, 0x76, 0x61, 0x82, 0x6a, 0x71, 0x3a, 0xd2, 0x06, 0x14,
	0xd9, 0x74, 0xb4, 0x32, 0xf1, 0x33, 0xd4, 0x54, 0x90, 0xa7, 0x3c, 0x76, 0x4d, 0x8e, 0x34, 0x18,
	0x82, 0x9c, 0x62, 0x67, 0xc9, 0xdd, 0x07, 0x83, 0xdb, 0xa8, 0x3e, 0xd2, 0xa0, 0xdc, 0x9d, 0x1b,
	0xb7, 0x4a, 0x99, 0xfb, 0x71, 0x0b, 0xd5, 0x73, 0xaa, 0xf5, 0x44, 0x2a, 0x46, 0xb6, 0x8a, 0x6e,
	0x57, 0x36, 0xfe, 0x1c, 0xed, 0xdf, 0x50, 0x11, 0xb9, 0x61, 0x8d, 0x72, 0x50, 0x43, 0x1d, 0x5f,
	0x41, 0x46, 0xc9, 0xb6, 0xef, 0x75, 0xea, 0x21, 0xbe, 0xa1, 0xe2, 0x8d, 0x45, 0x6f, 0xe7, 0x04,
	0xbf, 0x40, 0x7b, 0x8b, 0x13, 0x3a, 0x95, 0x93, 0x54, 0x26, 0x64, 0xc7, 0xc9, 0x9b, 0x95, 0xbc,
	0x5f, 0xb8, 0x71, 0x88, 0xb6, 0xe3, 0x91, 0x36, 0x32, 0x8b, 0x52, 0x3a, 0x80, 0x54, 0x93, 0xa6,
	0xbf, 0xda, 0x69, 0x1c, 0x1f, 0xde, 0xf9, 0x43, 0x97, 0xa7, 0xa6, 0x7b, 0xea, 0x0e, 0x9c, 0x3b,
	0xfd, 0x99, 0x30, 0x6a, 0x16, 0x6e, 0xc5, 0x4b, 0x2e, 0x7c, 0x8c, 0x0e, 0xf4, 0x35, 0xcf, 0xa3,
	0x58, 0x0a, 0x01, 0xb1, 0x6b, 0x4f, 0x7c, 0x05, 0xf1, 0x35, 0xd9, 0x75, 0x35, 0x3c, 0xb0, 0xf0,
	0x74, 0xce, 0x4e, 0x2d, 0xc2, 0x2f, 0xd1, 0xff, 0x19, 0xd7, 0x74, 0x90, 0x42, 0x74, 0x33, 0x02,
	0x35, 0x8b, 0x60, 0x4a, 0xb3, 0x3c, 0x05, 0x4d, 0xf6, 0xdc, 0xa1, 0xfd, 0x92, 0x5e, 0x5a, 0x78,
	0x56, 0x32, 0x7c, 0x84, 0x0e, 0x32, 0x3a, 0xad, 0xee, 0x18, 0x0d, 0x79, 0x0a, 0x91, 0xe6, 0xef,
	0x81, 0x60, 0xdf, 0xeb, 0xac, 0x86, 0x38, 0xa3, 0xd3, 0xf2, 0xa2, 0xdf, 0xf0, 0x14, 0xfa, 0xfc,
	0x3d, 0xe0, 0x5d, 0xb4, 0x6a, 0x52, 0x4d, 0x1e, 0xb8, 0xa8, 0xf6, 0x13, 0x7f, 0x82, 0x9a, 0x26,
	0xd5, 0x91, 0x2b, 0x79, 0x0c, 0x8a, 0x0f, 0x67, 0x64, 0xdf, 0xd1, 0x6d, 0x93, 0xea, 0xfe, 0x35,
	0xcf, 0xbf, 0x77, 0x4e, 0x1c, 0xa0, 0x7d, 0x9b, 0x4c, 0x8c, 0xb2, 0x01, 0xa8, 0x48, 0x0e, 0x23,
	0x63, 0x0b, 0xd2, 0xe4, 0xc0, 0xf7, 0x3a, 0xeb, 0xe1, 0x5e, 0x46, 0xa7, 0x17, 0x0e, 0x7d, 0x37,
	0x7c, 0xe7, 0x40, 0xeb, 0x6b, 0xb4, 0xf7, 0x41, 0xab, 0x6c, 0xfe, 0x6b, 0x98, 0x95, 0xeb, 0x66,
	0x3f, 0xf1, 0x3e, 0x5a, 0x1f, 0xd3, 0x74, 0x54, 0xed, 0x59, 0x61, 0xbc, 0x5a, 0xf9, 0xca, 0x6b,
	0xff, 0xb6, 0x82, 0x76, 0x17, 0xcd, 0xd7, 0xb9, 0x14, 0x1a, 0xf0, 0x11, 0xaa, 0x95, 0xbb, 0xe2,
	0x82, 0x34, 0x8e, 0x1f, 0x76, 0xe7, 0xef, 0x47, 0xd7, 0x49, 0xfb, 0x05, 0x0e, 0x2b, 0x1d, 0x3e,
	0x41, 0x4d, 0x37, 0x0c, 0x2c, 0x82, 0xa9, 0x5d, 0x11, 0x50, 0x2e, 0x97, 0xdd, 0xdb, 0x3b, 0x47,
	0xd9, 0x59, 0x29, 0x08, 0x77, 0x8a, 0x13, 0x95, 0x8d, 0xdf, 0xfd, 0xcb, 0x18, 0x16, 0x0f, 0x40,
	0x7b, 0x29, 0xd0, 0x65, 0xef, 0xc2, 0xc5, 0xb2, 0x13, 0xd9, 0x77, 0x22, 0xb7, 0x6f, 0xf7, 0x8e,
	0xea, 0xb7, 0xf7, 0x8d, 0xea, 0x9a, 0x0b, 0xf9, 0xe4, 0x9e, 0x90, 0xe5, 0x9f, 0x59, 0xc4, 0xbb,
	0x3b, 0xcb, 0xc7, 0x1a, 0xad, 0x3b, 0x15, 0xfe, 0x09, 0xd5, 0xab, 0xb6, 0xe1, 0x47, 0xff, 0x31,
	0xc9, 0xad, 0xc7, 0xf7, 0xc3, 0xa2, 0xd3, 0xed, 0xa7, 0x3f, 0xff, 0xf1, 0xf7, 0xaf, 0x2b, 0x1f,
	0xb7, 0x49, 0x30, 0x3e, 0x0a, 0x16, 0xc2, 0xc0, 0xa9, 0x82, 0x1e, 0x63, 0xaf, 0xbc, 0x17, 0x27,
	0xd1, 0x2f, 0xbd, 0x8b, 0xf0, 0x1c, 0xd5, 0x18, 0x0c, 0xe9, 0x28, 0x35, 0xb8, 0x87, 0x70, 0x4f,
	0xf8, 0xa0, 0x94, 0x54, 0xbe, 0x2a, 0x23, 0x75, 0xf1, 0xa7, 0xe8, 0x79, 0xeb, 0xd9, 0xd3, 0x80,
	0xc1, 0x90, 0x0b, 0x5e, 0xbc, 0xd0, 0xcb, 0x0f, 0xfb, 0x99, 0x95, 0x57, 0x79, 0x7f, 0xdc, 0x5a,
	0x46, 0x83, 0x0d, 0xf7, 0x7c, 0x7f, 0xf1, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8f, 0xf3, 0x71,
	0xd3, 0xbc, 0x06, 0x00, 0x00,
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
