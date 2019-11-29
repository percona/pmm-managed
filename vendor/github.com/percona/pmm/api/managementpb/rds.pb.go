// Code generated by protoc-gen-go. DO NOT EDIT.
// source: managementpb/rds.proto

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

// DiscoverRDSEngine describes supported RDS instance engines.
type DiscoverRDSEngine int32

const (
	DiscoverRDSEngine_DISCOVER_RDS_ENGINE_INVALID DiscoverRDSEngine = 0
	DiscoverRDSEngine_DISCOVER_RDS_MYSQL          DiscoverRDSEngine = 1
)

var DiscoverRDSEngine_name = map[int32]string{
	0: "DISCOVER_RDS_ENGINE_INVALID",
	1: "DISCOVER_RDS_MYSQL",
}

var DiscoverRDSEngine_value = map[string]int32{
	"DISCOVER_RDS_ENGINE_INVALID": 0,
	"DISCOVER_RDS_MYSQL":          1,
}

func (x DiscoverRDSEngine) String() string {
	return proto.EnumName(DiscoverRDSEngine_name, int32(x))
}

func (DiscoverRDSEngine) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c5c873fbc544be02, []int{0}
}

// DiscoverRDSInstance models an unique RDS instance for the list of instances returned by Discovery.
type DiscoverRDSInstance struct {
	// AWS region.
	Region string `protobuf:"bytes,1,opt,name=region,proto3" json:"region,omitempty"`
	// AWS availability zone.
	Az string `protobuf:"bytes,2,opt,name=az,proto3" json:"az,omitempty"`
	// AWS instance ID.
	InstanceId string `protobuf:"bytes,3,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	// AWS instance class.
	NodeModel string `protobuf:"bytes,4,opt,name=node_model,json=nodeModel,proto3" json:"node_model,omitempty"`
	// Address used to connect to it.
	Address string `protobuf:"bytes,5,opt,name=address,proto3" json:"address,omitempty"`
	// Access port.
	Port uint32 `protobuf:"varint,6,opt,name=port,proto3" json:"port,omitempty"`
	// Instance engine.
	Engine DiscoverRDSEngine `protobuf:"varint,7,opt,name=engine,proto3,enum=management.DiscoverRDSEngine" json:"engine,omitempty"`
	// Engine version.
	EngineVersion        string   `protobuf:"bytes,8,opt,name=engine_version,json=engineVersion,proto3" json:"engine_version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DiscoverRDSInstance) Reset()         { *m = DiscoverRDSInstance{} }
func (m *DiscoverRDSInstance) String() string { return proto.CompactTextString(m) }
func (*DiscoverRDSInstance) ProtoMessage()    {}
func (*DiscoverRDSInstance) Descriptor() ([]byte, []int) {
	return fileDescriptor_c5c873fbc544be02, []int{0}
}

func (m *DiscoverRDSInstance) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiscoverRDSInstance.Unmarshal(m, b)
}
func (m *DiscoverRDSInstance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiscoverRDSInstance.Marshal(b, m, deterministic)
}
func (m *DiscoverRDSInstance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiscoverRDSInstance.Merge(m, src)
}
func (m *DiscoverRDSInstance) XXX_Size() int {
	return xxx_messageInfo_DiscoverRDSInstance.Size(m)
}
func (m *DiscoverRDSInstance) XXX_DiscardUnknown() {
	xxx_messageInfo_DiscoverRDSInstance.DiscardUnknown(m)
}

var xxx_messageInfo_DiscoverRDSInstance proto.InternalMessageInfo

func (m *DiscoverRDSInstance) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *DiscoverRDSInstance) GetAz() string {
	if m != nil {
		return m.Az
	}
	return ""
}

func (m *DiscoverRDSInstance) GetInstanceId() string {
	if m != nil {
		return m.InstanceId
	}
	return ""
}

func (m *DiscoverRDSInstance) GetNodeModel() string {
	if m != nil {
		return m.NodeModel
	}
	return ""
}

func (m *DiscoverRDSInstance) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *DiscoverRDSInstance) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *DiscoverRDSInstance) GetEngine() DiscoverRDSEngine {
	if m != nil {
		return m.Engine
	}
	return DiscoverRDSEngine_DISCOVER_RDS_ENGINE_INVALID
}

func (m *DiscoverRDSInstance) GetEngineVersion() string {
	if m != nil {
		return m.EngineVersion
	}
	return ""
}

type DiscoverRDSRequest struct {
	// AWS Access key. Optional.
	AwsAccessKey string `protobuf:"bytes,1,opt,name=aws_access_key,json=awsAccessKey,proto3" json:"aws_access_key,omitempty"`
	// AWS Secret key. Optional.
	AwsSecretKey         string   `protobuf:"bytes,2,opt,name=aws_secret_key,json=awsSecretKey,proto3" json:"aws_secret_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DiscoverRDSRequest) Reset()         { *m = DiscoverRDSRequest{} }
func (m *DiscoverRDSRequest) String() string { return proto.CompactTextString(m) }
func (*DiscoverRDSRequest) ProtoMessage()    {}
func (*DiscoverRDSRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c5c873fbc544be02, []int{1}
}

func (m *DiscoverRDSRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiscoverRDSRequest.Unmarshal(m, b)
}
func (m *DiscoverRDSRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiscoverRDSRequest.Marshal(b, m, deterministic)
}
func (m *DiscoverRDSRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiscoverRDSRequest.Merge(m, src)
}
func (m *DiscoverRDSRequest) XXX_Size() int {
	return xxx_messageInfo_DiscoverRDSRequest.Size(m)
}
func (m *DiscoverRDSRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DiscoverRDSRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DiscoverRDSRequest proto.InternalMessageInfo

func (m *DiscoverRDSRequest) GetAwsAccessKey() string {
	if m != nil {
		return m.AwsAccessKey
	}
	return ""
}

func (m *DiscoverRDSRequest) GetAwsSecretKey() string {
	if m != nil {
		return m.AwsSecretKey
	}
	return ""
}

type DiscoverRDSResponse struct {
	RdsInstances         []*DiscoverRDSInstance `protobuf:"bytes,1,rep,name=rds_instances,json=rdsInstances,proto3" json:"rds_instances,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *DiscoverRDSResponse) Reset()         { *m = DiscoverRDSResponse{} }
func (m *DiscoverRDSResponse) String() string { return proto.CompactTextString(m) }
func (*DiscoverRDSResponse) ProtoMessage()    {}
func (*DiscoverRDSResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c5c873fbc544be02, []int{2}
}

func (m *DiscoverRDSResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiscoverRDSResponse.Unmarshal(m, b)
}
func (m *DiscoverRDSResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiscoverRDSResponse.Marshal(b, m, deterministic)
}
func (m *DiscoverRDSResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiscoverRDSResponse.Merge(m, src)
}
func (m *DiscoverRDSResponse) XXX_Size() int {
	return xxx_messageInfo_DiscoverRDSResponse.Size(m)
}
func (m *DiscoverRDSResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DiscoverRDSResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DiscoverRDSResponse proto.InternalMessageInfo

func (m *DiscoverRDSResponse) GetRdsInstances() []*DiscoverRDSInstance {
	if m != nil {
		return m.RdsInstances
	}
	return nil
}

type AddRDSRequest struct {
	// AWS region.
	Region string `protobuf:"bytes,1,opt,name=region,proto3" json:"region,omitempty"`
	// AWS availability zone.
	Az string `protobuf:"bytes,2,opt,name=az,proto3" json:"az,omitempty"`
	// AWS instance ID.
	InstanceId string `protobuf:"bytes,3,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
	// AWS instance class.
	NodeModel string `protobuf:"bytes,4,opt,name=node_model,json=nodeModel,proto3" json:"node_model,omitempty"`
	// Address used to connect to it.
	Address string `protobuf:"bytes,5,opt,name=address,proto3" json:"address,omitempty"`
	// Access port.
	Port uint32 `protobuf:"varint,6,opt,name=port,proto3" json:"port,omitempty"`
	// Instance engine.
	Engine DiscoverRDSEngine `protobuf:"varint,7,opt,name=engine,proto3,enum=management.DiscoverRDSEngine" json:"engine,omitempty"`
	// Unique across all Nodes user-defined name. Defaults to AWS instance ID.
	NodeName string `protobuf:"bytes,8,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	// Unique across all Services user-defined name. Defaults to AWS instance ID.
	ServiceName string `protobuf:"bytes,9,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	// Environment name.
	Environment string `protobuf:"bytes,10,opt,name=environment,proto3" json:"environment,omitempty"`
	// Cluster name.
	Cluster string `protobuf:"bytes,11,opt,name=cluster,proto3" json:"cluster,omitempty"`
	// Replication set name.
	ReplicationSet string `protobuf:"bytes,12,opt,name=replication_set,json=replicationSet,proto3" json:"replication_set,omitempty"`
	// Username for scraping metrics.
	Username string `protobuf:"bytes,13,opt,name=username,proto3" json:"username,omitempty"`
	// Password for scraping metrics.
	Password string `protobuf:"bytes,14,opt,name=password,proto3" json:"password,omitempty"`
	// AWS Access key.
	AwsAccessKey string `protobuf:"bytes,15,opt,name=aws_access_key,json=awsAccessKey,proto3" json:"aws_access_key,omitempty"`
	// AWS Secret key.
	AwsSecretKey string `protobuf:"bytes,16,opt,name=aws_secret_key,json=awsSecretKey,proto3" json:"aws_secret_key,omitempty"`
	// If true, adds rds_exporter.
	RdsExporter bool `protobuf:"varint,17,opt,name=rds_exporter,json=rdsExporter,proto3" json:"rds_exporter,omitempty"`
	// If true, adds qan-mysql-perfschema-agent.
	QanMysqlPerfschema bool `protobuf:"varint,18,opt,name=qan_mysql_perfschema,json=qanMysqlPerfschema,proto3" json:"qan_mysql_perfschema,omitempty"`
	// Custom user-assigned labels for service.
	CustomLabels map[string]string `protobuf:"bytes,19,rep,name=custom_labels,json=customLabels,proto3" json:"custom_labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Skip connection check.
	SkipConnectionCheck bool `protobuf:"varint,20,opt,name=skip_connection_check,json=skipConnectionCheck,proto3" json:"skip_connection_check,omitempty"`
	// Use TLS for database connections.
	Tls bool `protobuf:"varint,21,opt,name=tls,proto3" json:"tls,omitempty"`
	// Skip TLS certificate and hostname validation.
	TlsSkipVerify bool `protobuf:"varint,22,opt,name=tls_skip_verify,json=tlsSkipVerify,proto3" json:"tls_skip_verify,omitempty"`
	// Disable query examples.
	DisableQueryExamples bool `protobuf:"varint,23,opt,name=disable_query_examples,json=disableQueryExamples,proto3" json:"disable_query_examples,omitempty"`
	// Tablestats group collectors will be disabled if there are more than that number of tables.
	// If zero, server's default value is used.
	// Use negative value to disable them.
	TablestatsGroupTableLimit int32    `protobuf:"varint,24,opt,name=tablestats_group_table_limit,json=tablestatsGroupTableLimit,proto3" json:"tablestats_group_table_limit,omitempty"`
	XXX_NoUnkeyedLiteral      struct{} `json:"-"`
	XXX_unrecognized          []byte   `json:"-"`
	XXX_sizecache             int32    `json:"-"`
}

func (m *AddRDSRequest) Reset()         { *m = AddRDSRequest{} }
func (m *AddRDSRequest) String() string { return proto.CompactTextString(m) }
func (*AddRDSRequest) ProtoMessage()    {}
func (*AddRDSRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c5c873fbc544be02, []int{3}
}

func (m *AddRDSRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddRDSRequest.Unmarshal(m, b)
}
func (m *AddRDSRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddRDSRequest.Marshal(b, m, deterministic)
}
func (m *AddRDSRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddRDSRequest.Merge(m, src)
}
func (m *AddRDSRequest) XXX_Size() int {
	return xxx_messageInfo_AddRDSRequest.Size(m)
}
func (m *AddRDSRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddRDSRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddRDSRequest proto.InternalMessageInfo

func (m *AddRDSRequest) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *AddRDSRequest) GetAz() string {
	if m != nil {
		return m.Az
	}
	return ""
}

func (m *AddRDSRequest) GetInstanceId() string {
	if m != nil {
		return m.InstanceId
	}
	return ""
}

func (m *AddRDSRequest) GetNodeModel() string {
	if m != nil {
		return m.NodeModel
	}
	return ""
}

func (m *AddRDSRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *AddRDSRequest) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *AddRDSRequest) GetEngine() DiscoverRDSEngine {
	if m != nil {
		return m.Engine
	}
	return DiscoverRDSEngine_DISCOVER_RDS_ENGINE_INVALID
}

func (m *AddRDSRequest) GetNodeName() string {
	if m != nil {
		return m.NodeName
	}
	return ""
}

func (m *AddRDSRequest) GetServiceName() string {
	if m != nil {
		return m.ServiceName
	}
	return ""
}

func (m *AddRDSRequest) GetEnvironment() string {
	if m != nil {
		return m.Environment
	}
	return ""
}

func (m *AddRDSRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *AddRDSRequest) GetReplicationSet() string {
	if m != nil {
		return m.ReplicationSet
	}
	return ""
}

func (m *AddRDSRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *AddRDSRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *AddRDSRequest) GetAwsAccessKey() string {
	if m != nil {
		return m.AwsAccessKey
	}
	return ""
}

func (m *AddRDSRequest) GetAwsSecretKey() string {
	if m != nil {
		return m.AwsSecretKey
	}
	return ""
}

func (m *AddRDSRequest) GetRdsExporter() bool {
	if m != nil {
		return m.RdsExporter
	}
	return false
}

func (m *AddRDSRequest) GetQanMysqlPerfschema() bool {
	if m != nil {
		return m.QanMysqlPerfschema
	}
	return false
}

func (m *AddRDSRequest) GetCustomLabels() map[string]string {
	if m != nil {
		return m.CustomLabels
	}
	return nil
}

func (m *AddRDSRequest) GetSkipConnectionCheck() bool {
	if m != nil {
		return m.SkipConnectionCheck
	}
	return false
}

func (m *AddRDSRequest) GetTls() bool {
	if m != nil {
		return m.Tls
	}
	return false
}

func (m *AddRDSRequest) GetTlsSkipVerify() bool {
	if m != nil {
		return m.TlsSkipVerify
	}
	return false
}

func (m *AddRDSRequest) GetDisableQueryExamples() bool {
	if m != nil {
		return m.DisableQueryExamples
	}
	return false
}

func (m *AddRDSRequest) GetTablestatsGroupTableLimit() int32 {
	if m != nil {
		return m.TablestatsGroupTableLimit
	}
	return 0
}

type AddRDSResponse struct {
	Node               *inventorypb.RemoteRDSNode           `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	RdsExporter        *inventorypb.RDSExporter             `protobuf:"bytes,2,opt,name=rds_exporter,json=rdsExporter,proto3" json:"rds_exporter,omitempty"`
	Mysql              *inventorypb.MySQLService            `protobuf:"bytes,3,opt,name=mysql,proto3" json:"mysql,omitempty"`
	MysqldExporter     *inventorypb.MySQLdExporter          `protobuf:"bytes,4,opt,name=mysqld_exporter,json=mysqldExporter,proto3" json:"mysqld_exporter,omitempty"`
	QanMysqlPerfschema *inventorypb.QANMySQLPerfSchemaAgent `protobuf:"bytes,5,opt,name=qan_mysql_perfschema,json=qanMysqlPerfschema,proto3" json:"qan_mysql_perfschema,omitempty"`
	// Actual table count at the moment of adding.
	TableCount           int32    `protobuf:"varint,6,opt,name=table_count,json=tableCount,proto3" json:"table_count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddRDSResponse) Reset()         { *m = AddRDSResponse{} }
func (m *AddRDSResponse) String() string { return proto.CompactTextString(m) }
func (*AddRDSResponse) ProtoMessage()    {}
func (*AddRDSResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c5c873fbc544be02, []int{4}
}

func (m *AddRDSResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddRDSResponse.Unmarshal(m, b)
}
func (m *AddRDSResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddRDSResponse.Marshal(b, m, deterministic)
}
func (m *AddRDSResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddRDSResponse.Merge(m, src)
}
func (m *AddRDSResponse) XXX_Size() int {
	return xxx_messageInfo_AddRDSResponse.Size(m)
}
func (m *AddRDSResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddRDSResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddRDSResponse proto.InternalMessageInfo

func (m *AddRDSResponse) GetNode() *inventorypb.RemoteRDSNode {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *AddRDSResponse) GetRdsExporter() *inventorypb.RDSExporter {
	if m != nil {
		return m.RdsExporter
	}
	return nil
}

func (m *AddRDSResponse) GetMysql() *inventorypb.MySQLService {
	if m != nil {
		return m.Mysql
	}
	return nil
}

func (m *AddRDSResponse) GetMysqldExporter() *inventorypb.MySQLdExporter {
	if m != nil {
		return m.MysqldExporter
	}
	return nil
}

func (m *AddRDSResponse) GetQanMysqlPerfschema() *inventorypb.QANMySQLPerfSchemaAgent {
	if m != nil {
		return m.QanMysqlPerfschema
	}
	return nil
}

func (m *AddRDSResponse) GetTableCount() int32 {
	if m != nil {
		return m.TableCount
	}
	return 0
}

func init() {
	proto.RegisterEnum("management.DiscoverRDSEngine", DiscoverRDSEngine_name, DiscoverRDSEngine_value)
	proto.RegisterType((*DiscoverRDSInstance)(nil), "management.DiscoverRDSInstance")
	proto.RegisterType((*DiscoverRDSRequest)(nil), "management.DiscoverRDSRequest")
	proto.RegisterType((*DiscoverRDSResponse)(nil), "management.DiscoverRDSResponse")
	proto.RegisterType((*AddRDSRequest)(nil), "management.AddRDSRequest")
	proto.RegisterMapType((map[string]string)(nil), "management.AddRDSRequest.CustomLabelsEntry")
	proto.RegisterType((*AddRDSResponse)(nil), "management.AddRDSResponse")
}

func init() { proto.RegisterFile("managementpb/rds.proto", fileDescriptor_c5c873fbc544be02) }

var fileDescriptor_c5c873fbc544be02 = []byte{
	// 1171 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0x4b, 0x73, 0xdb, 0x36,
	0x17, 0x0d, 0xe5, 0x67, 0x20, 0x4b, 0x76, 0x10, 0x47, 0x41, 0xe4, 0x24, 0x56, 0xf4, 0x7d, 0x4d,
	0xdc, 0x24, 0x16, 0x5b, 0xf7, 0x31, 0x6d, 0x36, 0x19, 0xc5, 0xd2, 0x64, 0x3c, 0x91, 0xd5, 0x98,
	0xcc, 0x78, 0xfa, 0x58, 0xb0, 0x10, 0x79, 0xad, 0x70, 0x4c, 0x02, 0x34, 0x00, 0xc9, 0x51, 0x96,
	0xdd, 0x75, 0xdb, 0xfe, 0xb4, 0xfe, 0x80, 0x4e, 0x1f, 0xbf, 0xa2, 0xab, 0x0e, 0x40, 0x52, 0xa6,
	0x1f, 0x9a, 0x69, 0xa7, 0x2b, 0x01, 0xe7, 0x1c, 0xdc, 0x0b, 0xdc, 0x7b, 0x00, 0x11, 0xd5, 0x62,
	0xca, 0xe8, 0x10, 0x62, 0x60, 0x2a, 0x19, 0xd8, 0x22, 0x90, 0xad, 0x44, 0x70, 0xc5, 0x31, 0x3a,
	0xc3, 0xeb, 0x9f, 0x0f, 0x43, 0xf5, 0x76, 0x34, 0x68, 0xf9, 0x3c, 0xb6, 0xe3, 0xd3, 0x50, 0x1d,
	0xf3, 0x53, 0x7b, 0xc8, 0xb7, 0x8d, 0x70, 0x7b, 0x4c, 0xa3, 0x30, 0xa0, 0x8a, 0x0b, 0x69, 0x4f,
	0x87, 0x69, 0x8c, 0xfa, 0xdd, 0x21, 0xe7, 0xc3, 0x08, 0x6c, 0x9a, 0x84, 0x36, 0x65, 0x8c, 0x2b,
	0xaa, 0x42, 0xce, 0xb2, 0x0c, 0x75, 0x12, 0xb2, 0x31, 0x30, 0xc5, 0xc5, 0x24, 0x19, 0xd8, 0x74,
	0x08, 0x4c, 0xe5, 0xcc, 0xed, 0x22, 0xc3, 0x78, 0x00, 0x39, 0x51, 0x2f, 0x12, 0x12, 0xc4, 0x38,
	0xf4, 0xa7, 0xdc, 0x53, 0xf3, 0xe3, 0x6f, 0x0f, 0x81, 0x6d, 0xcb, 0x53, 0x3a, 0x1c, 0x82, 0xb0,
	0x79, 0x62, 0x12, 0x5e, 0x4e, 0xde, 0xfc, 0xb1, 0x84, 0x6e, 0x76, 0x42, 0xe9, 0xf3, 0x31, 0x08,
	0xa7, 0xe3, 0xee, 0x31, 0xa9, 0x28, 0xf3, 0x01, 0xd7, 0xd0, 0xa2, 0x80, 0x61, 0xc8, 0x19, 0xb1,
	0x1a, 0xd6, 0xd6, 0x75, 0x27, 0x9b, 0xe1, 0x2a, 0x2a, 0xd1, 0xf7, 0xa4, 0x64, 0xb0, 0x12, 0x7d,
	0x8f, 0x37, 0x51, 0x39, 0xcc, 0xd6, 0x78, 0x61, 0x40, 0xe6, 0x0c, 0x81, 0x72, 0x68, 0x2f, 0xc0,
	0xf7, 0x10, 0xd2, 0x3b, 0xf7, 0x62, 0x1e, 0x40, 0x44, 0xe6, 0x0d, 0x7f, 0x5d, 0x23, 0xfb, 0x1a,
	0xc0, 0x04, 0x2d, 0xd1, 0x20, 0x10, 0x20, 0x25, 0x59, 0x30, 0x5c, 0x3e, 0xc5, 0x18, 0xcd, 0x27,
	0x5c, 0x28, 0xb2, 0xd8, 0xb0, 0xb6, 0x2a, 0x8e, 0x19, 0xe3, 0xcf, 0xd0, 0x22, 0xb0, 0x61, 0xc8,
	0x80, 0x2c, 0x35, 0xac, 0xad, 0xea, 0xce, 0xbd, 0xd6, 0x59, 0x77, 0x5a, 0x85, 0x63, 0x74, 0x8d,
	0xc8, 0xc9, 0xc4, 0xf8, 0x03, 0x54, 0x4d, 0x47, 0xde, 0x18, 0x84, 0xd4, 0x87, 0x5a, 0x36, 0xb9,
	0x2a, 0x29, 0x7a, 0x98, 0x82, 0xcd, 0xef, 0x11, 0x2e, 0xc4, 0x70, 0xe0, 0x64, 0x04, 0x52, 0xe1,
	0xff, 0xa3, 0x2a, 0x3d, 0x95, 0x1e, 0xf5, 0x7d, 0x90, 0xd2, 0x3b, 0x86, 0x49, 0x56, 0x91, 0x15,
	0x7a, 0x2a, 0xdb, 0x06, 0x7c, 0x05, 0x93, 0x5c, 0x25, 0xc1, 0x17, 0xa0, 0x8c, 0xaa, 0x34, 0x55,
	0xb9, 0x06, 0x7c, 0x05, 0x93, 0xe6, 0x77, 0xe7, 0x8a, 0xed, 0x80, 0x4c, 0x38, 0x93, 0x80, 0x3b,
	0xa8, 0x22, 0x02, 0xe9, 0xe5, 0x55, 0x93, 0xc4, 0x6a, 0xcc, 0x6d, 0x95, 0x77, 0x36, 0x67, 0x9c,
	0x2e, 0x6f, 0x92, 0xb3, 0x22, 0x02, 0x99, 0x4f, 0x64, 0xf3, 0xaf, 0x25, 0x54, 0x69, 0x07, 0x41,
	0x61, 0xeb, 0xf7, 0xcf, 0x37, 0xf1, 0xc5, 0xe2, 0xef, 0xbf, 0x6e, 0x96, 0xbe, 0xb6, 0x66, 0x36,
	0xf3, 0xd1, 0x15, 0xcd, 0x9c, 0x2e, 0xfa, 0x17, 0x4d, 0x6d, 0x5c, 0x68, 0xea, 0x34, 0xc6, 0xb4,
	0xb9, 0xf5, 0x62, 0x73, 0x53, 0x7a, 0xed, 0xda, 0x7f, 0x6b, 0xf2, 0x06, 0x32, 0x3b, 0xf0, 0x18,
	0x8d, 0x21, 0xeb, 0xef, 0xb2, 0x06, 0xfa, 0x34, 0x06, 0xfc, 0x00, 0xad, 0x64, 0xd7, 0x24, 0xe5,
	0xaf, 0x1b, 0xbe, 0x9c, 0x61, 0x46, 0xd2, 0x40, 0x65, 0x60, 0xe3, 0x50, 0x70, 0xa6, 0x13, 0x11,
	0x94, 0x2a, 0x0a, 0x90, 0xf6, 0xaa, 0x1f, 0x8d, 0xa4, 0x02, 0x41, 0xca, 0xa9, 0x57, 0xb3, 0x29,
	0x7e, 0x84, 0x56, 0x05, 0x24, 0x51, 0xe8, 0x9b, 0xbb, 0xe5, 0x49, 0x50, 0x64, 0xc5, 0x28, 0xaa,
	0x05, 0xd8, 0x05, 0x85, 0x9b, 0x68, 0x79, 0x24, 0x41, 0x98, 0x3d, 0x54, 0xce, 0x95, 0x66, 0x8a,
	0xe3, 0x3a, 0x5a, 0x4e, 0xa8, 0x94, 0xa7, 0x5c, 0x04, 0xa4, 0x9a, 0x9e, 0x23, 0x9f, 0x5f, 0x61,
	0xc6, 0xd5, 0x7f, 0x64, 0xc6, 0xb5, 0xcb, 0x66, 0xd4, 0x35, 0xd1, 0xae, 0x83, 0x77, 0xba, 0xea,
	0x20, 0xc8, 0x8d, 0x86, 0xb5, 0xb5, 0xec, 0x94, 0x45, 0x20, 0xbb, 0x19, 0x84, 0x3f, 0x42, 0xeb,
	0x27, 0x94, 0x79, 0xf1, 0x44, 0x9e, 0x44, 0x5e, 0x02, 0xe2, 0x48, 0xfa, 0x6f, 0x21, 0xa6, 0x04,
	0x1b, 0x29, 0x3e, 0xa1, 0x6c, 0x5f, 0x53, 0xaf, 0xa7, 0x0c, 0x7e, 0x8d, 0x2a, 0xfe, 0x48, 0x2a,
	0x1e, 0x7b, 0x11, 0x1d, 0x40, 0x24, 0xc9, 0x4d, 0x63, 0xe5, 0x27, 0xc5, 0x1e, 0x9e, 0x33, 0x69,
	0x6b, 0xd7, 0xc8, 0x7b, 0x46, 0xdd, 0x65, 0x4a, 0x4c, 0x9c, 0x15, 0xbf, 0x00, 0xe1, 0x1d, 0x74,
	0x4b, 0x1e, 0x87, 0x89, 0xe7, 0x73, 0xc6, 0xc0, 0x37, 0xf5, 0xf5, 0xdf, 0x82, 0x7f, 0x4c, 0xd6,
	0xcd, 0x26, 0x6e, 0x6a, 0x72, 0x77, 0xca, 0xed, 0x6a, 0x0a, 0xaf, 0xa1, 0x39, 0x15, 0x49, 0x72,
	0xcb, 0x28, 0xf4, 0x10, 0x3f, 0x44, 0xab, 0x2a, 0x92, 0x9e, 0x89, 0x34, 0x06, 0x11, 0x1e, 0x4d,
	0x48, 0xcd, 0xb0, 0x15, 0x15, 0x49, 0xf7, 0x38, 0x4c, 0x0e, 0x0d, 0x88, 0x3f, 0x45, 0xb5, 0x20,
	0x94, 0x74, 0x10, 0x81, 0x77, 0x32, 0x02, 0x31, 0xf1, 0xe0, 0x1d, 0x8d, 0x93, 0x08, 0x24, 0xb9,
	0x6d, 0xe4, 0xeb, 0x19, 0x7b, 0xa0, 0xc9, 0x6e, 0xc6, 0xe1, 0xe7, 0xe8, 0xae, 0xd2, 0xa8, 0x54,
	0x54, 0x49, 0x6f, 0x28, 0xf8, 0x28, 0xf1, 0x0c, 0xe0, 0x45, 0x61, 0x1c, 0x2a, 0x42, 0x1a, 0xd6,
	0xd6, 0x82, 0x73, 0xe7, 0x4c, 0xf3, 0x52, 0x4b, 0xde, 0xe8, 0x69, 0x4f, 0x0b, 0xea, 0xcf, 0xd1,
	0x8d, 0x4b, 0x75, 0xd0, 0xa7, 0x38, 0x7b, 0x6e, 0xf4, 0x10, 0xaf, 0xa3, 0x85, 0x31, 0x8d, 0x46,
	0x90, 0xdd, 0xd9, 0x74, 0xf2, 0xac, 0xf4, 0x85, 0xd5, 0xfc, 0xa3, 0x84, 0xaa, 0x79, 0x5d, 0xb3,
	0x57, 0xe5, 0x29, 0x9a, 0xd7, 0xfe, 0x37, 0xeb, 0xcb, 0x3b, 0xa4, 0x35, 0xfd, 0xcf, 0x68, 0x39,
	0x10, 0x73, 0x05, 0x4e, 0xc7, 0xed, 0xf3, 0x00, 0x1c, 0xa3, 0xc2, 0x5f, 0x5e, 0x70, 0x43, 0xc9,
	0xac, 0xaa, 0x15, 0x57, 0x75, 0xdc, 0xdc, 0x18, 0xe7, 0x5d, 0xb2, 0x8d, 0x16, 0x8c, 0x43, 0xcc,
	0x83, 0x51, 0xde, 0xb9, 0x5d, 0x58, 0xb3, 0x3f, 0x71, 0x0f, 0x7a, 0x6e, 0x7a, 0xcb, 0x9c, 0x54,
	0x85, 0x5f, 0xa0, 0x55, 0x33, 0x08, 0xce, 0x92, 0xcd, 0x9b, 0x85, 0x77, 0x2e, 0x2e, 0x0c, 0xa6,
	0xf9, 0xaa, 0xe9, 0x8a, 0x69, 0xca, 0x37, 0x33, 0x8c, 0xb9, 0x60, 0x02, 0x35, 0x0b, 0x81, 0x0e,
	0xda, 0x7d, 0x13, 0x4b, 0x7b, 0xd4, 0x35, 0xa2, 0xb6, 0xfe, 0x8b, 0xbd, 0xd2, 0xbc, 0x9b, 0xa8,
	0x9c, 0x76, 0xcd, 0xe7, 0x23, 0x96, 0x3e, 0x4e, 0x0b, 0x0e, 0x32, 0xd0, 0xae, 0x46, 0x1e, 0xf7,
	0xd0, 0x8d, 0x4b, 0x0f, 0x10, 0xde, 0x44, 0x1b, 0x9d, 0x3d, 0x77, 0xf7, 0xab, 0xc3, 0xae, 0xe3,
	0x39, 0x1d, 0xd7, 0xeb, 0xf6, 0x5f, 0xee, 0xf5, 0xbb, 0xde, 0x5e, 0xff, 0xb0, 0xdd, 0xdb, 0xeb,
	0xac, 0x5d, 0xc3, 0x35, 0x84, 0xcf, 0x09, 0xf6, 0xbf, 0x71, 0x0f, 0x7a, 0x6b, 0xd6, 0xce, 0x6f,
	0x16, 0x9a, 0x73, 0x3a, 0x2e, 0x1e, 0xa3, 0x72, 0x21, 0x2a, 0xbe, 0x3f, 0xe3, 0xbd, 0xcb, 0x2e,
	0x4c, 0x7d, 0x73, 0x26, 0x9f, 0x36, 0xbe, 0xf9, 0xf0, 0x87, 0x5f, 0xfe, 0xfc, 0xb9, 0xd4, 0x68,
	0x6e, 0xd8, 0xe3, 0x8f, 0xed, 0x33, 0xad, 0xed, 0x74, 0x5c, 0x3b, 0xd7, 0x3f, 0xb3, 0x1e, 0xe3,
	0x01, 0x5a, 0x4c, 0x2d, 0x83, 0xef, 0xcc, 0xbc, 0x9e, 0xf5, 0xfa, 0x55, 0x54, 0x96, 0xe8, 0x81,
	0x49, 0xb4, 0xd1, 0xac, 0x5d, 0x91, 0xa8, 0x1d, 0x04, 0xcf, 0xac, 0xc7, 0x2f, 0xbc, 0x9f, 0xda,
	0x7d, 0xa7, 0x87, 0x96, 0x02, 0x38, 0xa2, 0xa3, 0x48, 0xe1, 0x36, 0xc2, 0x6d, 0xd6, 0x00, 0x21,
	0xb8, 0x68, 0x88, 0x2c, 0x4e, 0x0b, 0x3f, 0x41, 0x1f, 0xd6, 0x1f, 0xfd, 0xcf, 0x0e, 0xe0, 0x28,
	0x64, 0x61, 0xfa, 0xa9, 0x52, 0xfc, 0x1c, 0xeb, 0x6a, 0x79, 0x9e, 0xf5, 0xdb, 0x95, 0x22, 0x35,
	0x58, 0x34, 0xdf, 0x31, 0x9f, 0xfc, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x66, 0xaa, 0xe4, 0x65, 0xc0,
	0x09, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RDSClient is the client API for RDS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RDSClient interface {
	// DiscoverRDS discovers RDS instances.
	DiscoverRDS(ctx context.Context, in *DiscoverRDSRequest, opts ...grpc.CallOption) (*DiscoverRDSResponse, error)
	// AddRDS adds RDS instance.
	AddRDS(ctx context.Context, in *AddRDSRequest, opts ...grpc.CallOption) (*AddRDSResponse, error)
}

type rDSClient struct {
	cc *grpc.ClientConn
}

func NewRDSClient(cc *grpc.ClientConn) RDSClient {
	return &rDSClient{cc}
}

func (c *rDSClient) DiscoverRDS(ctx context.Context, in *DiscoverRDSRequest, opts ...grpc.CallOption) (*DiscoverRDSResponse, error) {
	out := new(DiscoverRDSResponse)
	err := c.cc.Invoke(ctx, "/management.RDS/DiscoverRDS", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rDSClient) AddRDS(ctx context.Context, in *AddRDSRequest, opts ...grpc.CallOption) (*AddRDSResponse, error) {
	out := new(AddRDSResponse)
	err := c.cc.Invoke(ctx, "/management.RDS/AddRDS", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RDSServer is the server API for RDS service.
type RDSServer interface {
	// DiscoverRDS discovers RDS instances.
	DiscoverRDS(context.Context, *DiscoverRDSRequest) (*DiscoverRDSResponse, error)
	// AddRDS adds RDS instance.
	AddRDS(context.Context, *AddRDSRequest) (*AddRDSResponse, error)
}

// UnimplementedRDSServer can be embedded to have forward compatible implementations.
type UnimplementedRDSServer struct {
}

func (*UnimplementedRDSServer) DiscoverRDS(ctx context.Context, req *DiscoverRDSRequest) (*DiscoverRDSResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiscoverRDS not implemented")
}
func (*UnimplementedRDSServer) AddRDS(ctx context.Context, req *AddRDSRequest) (*AddRDSResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRDS not implemented")
}

func RegisterRDSServer(s *grpc.Server, srv RDSServer) {
	s.RegisterService(&_RDS_serviceDesc, srv)
}

func _RDS_DiscoverRDS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiscoverRDSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDSServer).DiscoverRDS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/management.RDS/DiscoverRDS",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDSServer).DiscoverRDS(ctx, req.(*DiscoverRDSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RDS_AddRDS_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRDSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RDSServer).AddRDS(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/management.RDS/AddRDS",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RDSServer).AddRDS(ctx, req.(*AddRDSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RDS_serviceDesc = grpc.ServiceDesc{
	ServiceName: "management.RDS",
	HandlerType: (*RDSServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DiscoverRDS",
			Handler:    _RDS_DiscoverRDS_Handler,
		},
		{
			MethodName: "AddRDS",
			Handler:    _RDS_AddRDS_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "managementpb/rds.proto",
}
