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
	// Custom user-assigned labels for Node and Service.
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
	TablestatsGroupTableLimit int32 `protobuf:"varint,24,opt,name=tablestats_group_table_limit,json=tablestatsGroupTableLimit,proto3" json:"tablestats_group_table_limit,omitempty"`
	// Disable basic metrics.
	DisableBasicMetrics bool `protobuf:"varint,25,opt,name=disable_basic_metrics,json=disableBasicMetrics,proto3" json:"disable_basic_metrics,omitempty"`
	// Disable enhanced metrics.
	DisableEnhancedMetrics bool     `protobuf:"varint,26,opt,name=disable_enhanced_metrics,json=disableEnhancedMetrics,proto3" json:"disable_enhanced_metrics,omitempty"`
	XXX_NoUnkeyedLiteral   struct{} `json:"-"`
	XXX_unrecognized       []byte   `json:"-"`
	XXX_sizecache          int32    `json:"-"`
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

func (m *AddRDSRequest) GetDisableBasicMetrics() bool {
	if m != nil {
		return m.DisableBasicMetrics
	}
	return false
}

func (m *AddRDSRequest) GetDisableEnhancedMetrics() bool {
	if m != nil {
		return m.DisableEnhancedMetrics
	}
	return false
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
	// 1214 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0x5b, 0x53, 0x1b, 0x37,
	0x14, 0x8e, 0xcd, 0x25, 0x44, 0xc6, 0x86, 0x28, 0xc4, 0x11, 0x26, 0x09, 0x8e, 0xdb, 0x26, 0x34,
	0x09, 0xde, 0x96, 0x5e, 0x26, 0xcd, 0x4b, 0xc6, 0x60, 0x4f, 0x86, 0x89, 0x71, 0xc3, 0x6e, 0x86,
	0xe9, 0xe5, 0x61, 0x2b, 0xef, 0x1e, 0xcc, 0x0e, 0xbb, 0xd2, 0x22, 0xc9, 0x26, 0xce, 0x63, 0xdf,
	0xfa, 0xda, 0xfe, 0xb4, 0xfe, 0x80, 0x4e, 0x2f, 0x0f, 0xfd, 0x19, 0x1d, 0x69, 0x2f, 0x2c, 0x01,
	0x66, 0xda, 0xe9, 0x93, 0xa5, 0xef, 0xfb, 0xce, 0x39, 0xd2, 0xb9, 0xc8, 0x8b, 0xea, 0x11, 0x65,
	0x74, 0x04, 0x11, 0x30, 0x15, 0x0f, 0x2d, 0xe1, 0xcb, 0x76, 0x2c, 0xb8, 0xe2, 0x18, 0x9d, 0xe1,
	0x8d, 0x2f, 0x47, 0x81, 0x3a, 0x1a, 0x0f, 0xdb, 0x1e, 0x8f, 0xac, 0xe8, 0x34, 0x50, 0xc7, 0xfc,
	0xd4, 0x1a, 0xf1, 0x4d, 0x23, 0xdc, 0x9c, 0xd0, 0x30, 0xf0, 0xa9, 0xe2, 0x42, 0x5a, 0xf9, 0x32,
	0xf1, 0xd1, 0xb8, 0x3b, 0xe2, 0x7c, 0x14, 0x82, 0x45, 0xe3, 0xc0, 0xa2, 0x8c, 0x71, 0x45, 0x55,
	0xc0, 0x59, 0x1a, 0xa1, 0x41, 0x02, 0x36, 0x01, 0xa6, 0xb8, 0x98, 0xc6, 0x43, 0x8b, 0x8e, 0x80,
	0xa9, 0x8c, 0xb9, 0x53, 0x64, 0x18, 0xf7, 0x21, 0x23, 0x1a, 0x45, 0x42, 0x82, 0x98, 0x04, 0x5e,
	0xce, 0x3d, 0x35, 0x3f, 0xde, 0xe6, 0x08, 0xd8, 0xa6, 0x3c, 0xa5, 0xa3, 0x11, 0x08, 0x8b, 0xc7,
	0x26, 0xe0, 0xc5, 0xe0, 0xad, 0x9f, 0xca, 0xe8, 0x56, 0x37, 0x90, 0x1e, 0x9f, 0x80, 0xb0, 0xbb,
	0xce, 0x2e, 0x93, 0x8a, 0x32, 0x0f, 0x70, 0x1d, 0xcd, 0x0b, 0x18, 0x05, 0x9c, 0x91, 0x52, 0xb3,
	0xb4, 0x71, 0xc3, 0x4e, 0x77, 0xb8, 0x86, 0xca, 0xf4, 0x1d, 0x29, 0x1b, 0xac, 0x4c, 0xdf, 0xe1,
	0x75, 0x54, 0x09, 0x52, 0x1b, 0x37, 0xf0, 0xc9, 0x8c, 0x21, 0x50, 0x06, 0xed, 0xfa, 0xf8, 0x1e,
	0x42, 0xfa, 0xe4, 0x6e, 0xc4, 0x7d, 0x08, 0xc9, 0xac, 0xe1, 0x6f, 0x68, 0x64, 0x4f, 0x03, 0x98,
	0xa0, 0xeb, 0xd4, 0xf7, 0x05, 0x48, 0x49, 0xe6, 0x0c, 0x97, 0x6d, 0x31, 0x46, 0xb3, 0x31, 0x17,
	0x8a, 0xcc, 0x37, 0x4b, 0x1b, 0x55, 0xdb, 0xac, 0xf1, 0x17, 0x68, 0x1e, 0xd8, 0x28, 0x60, 0x40,
	0xae, 0x37, 0x4b, 0x1b, 0xb5, 0xad, 0x7b, 0xed, 0xb3, 0xea, 0xb4, 0x0b, 0xd7, 0xe8, 0x19, 0x91,
	0x9d, 0x8a, 0xf1, 0x47, 0xa8, 0x96, 0xac, 0xdc, 0x09, 0x08, 0xa9, 0x2f, 0xb5, 0x60, 0x62, 0x55,
	0x13, 0xf4, 0x20, 0x01, 0x5b, 0x3f, 0x20, 0x5c, 0xf0, 0x61, 0xc3, 0xc9, 0x18, 0xa4, 0xc2, 0x1f,
	0xa2, 0x1a, 0x3d, 0x95, 0x2e, 0xf5, 0x3c, 0x90, 0xd2, 0x3d, 0x86, 0x69, 0x9a, 0x91, 0x45, 0x7a,
	0x2a, 0x3b, 0x06, 0x7c, 0x05, 0xd3, 0x4c, 0x25, 0xc1, 0x13, 0xa0, 0x8c, 0xaa, 0x9c, 0xab, 0x1c,
	0x03, 0xbe, 0x82, 0x69, 0xeb, 0xfb, 0x73, 0xc9, 0xb6, 0x41, 0xc6, 0x9c, 0x49, 0xc0, 0x5d, 0x54,
	0x15, 0xbe, 0x74, 0xb3, 0xac, 0x49, 0x52, 0x6a, 0xce, 0x6c, 0x54, 0xb6, 0xd6, 0xaf, 0xb8, 0x5d,
	0x56, 0x24, 0x7b, 0x51, 0xf8, 0x32, 0xdb, 0xc8, 0xd6, 0xdf, 0x0b, 0xa8, 0xda, 0xf1, 0xfd, 0xc2,
	0xd1, 0xef, 0x9f, 0x2f, 0xe2, 0xf6, 0xfc, 0x1f, 0xbf, 0xad, 0x97, 0xbf, 0x29, 0x5d, 0x59, 0xcc,
	0x47, 0x97, 0x14, 0x33, 0x37, 0xfa, 0x0f, 0x45, 0x6d, 0xbe, 0x57, 0xd4, 0xdc, 0x47, 0x5e, 0xdc,
	0x46, 0xb1, 0xb8, 0x09, 0xbd, 0x7c, 0xed, 0xff, 0x15, 0x79, 0x0d, 0x99, 0x13, 0xb8, 0x8c, 0x46,
	0x90, 0xd6, 0x77, 0x41, 0x03, 0x03, 0x1a, 0x01, 0x7e, 0x80, 0x16, 0xd3, 0x31, 0x49, 0xf8, 0x1b,
	0x86, 0xaf, 0xa4, 0x98, 0x91, 0x34, 0x51, 0x05, 0xd8, 0x24, 0x10, 0x9c, 0xe9, 0x40, 0x04, 0x25,
	0x8a, 0x02, 0xa4, 0x7b, 0xd5, 0x0b, 0xc7, 0x52, 0x81, 0x20, 0x95, 0xa4, 0x57, 0xd3, 0x2d, 0x7e,
	0x84, 0x96, 0x04, 0xc4, 0x61, 0xe0, 0x99, 0xd9, 0x72, 0x25, 0x28, 0xb2, 0x68, 0x14, 0xb5, 0x02,
	0xec, 0x80, 0xc2, 0x2d, 0xb4, 0x30, 0x96, 0x20, 0xcc, 0x19, 0xaa, 0xe7, 0x52, 0x93, 0xe3, 0xb8,
	0x81, 0x16, 0x62, 0x2a, 0xe5, 0x29, 0x17, 0x3e, 0xa9, 0x25, 0xf7, 0xc8, 0xf6, 0x97, 0x34, 0xe3,
	0xd2, 0xbf, 0x6a, 0xc6, 0xe5, 0x8b, 0xcd, 0xa8, 0x73, 0xa2, 0xbb, 0x0e, 0xde, 0xea, 0xac, 0x83,
	0x20, 0x37, 0x9b, 0xa5, 0x8d, 0x05, 0xbb, 0x22, 0x7c, 0xd9, 0x4b, 0x21, 0xfc, 0x09, 0x5a, 0x39,
	0xa1, 0xcc, 0x8d, 0xa6, 0xf2, 0x24, 0x74, 0x63, 0x10, 0x87, 0xd2, 0x3b, 0x82, 0x88, 0x12, 0x6c,
	0xa4, 0xf8, 0x84, 0xb2, 0x3d, 0x4d, 0xbd, 0xce, 0x19, 0xfc, 0x1a, 0x55, 0xbd, 0xb1, 0x54, 0x3c,
	0x72, 0x43, 0x3a, 0x84, 0x50, 0x92, 0x5b, 0xa6, 0x95, 0x9f, 0x14, 0x6b, 0x78, 0xae, 0x49, 0xdb,
	0x3b, 0x46, 0xde, 0x37, 0xea, 0x1e, 0x53, 0x62, 0x6a, 0x2f, 0x7a, 0x05, 0x08, 0x6f, 0xa1, 0xdb,
	0xf2, 0x38, 0x88, 0x5d, 0x8f, 0x33, 0x06, 0x9e, 0xc9, 0xaf, 0x77, 0x04, 0xde, 0x31, 0x59, 0x31,
	0x87, 0xb8, 0xa5, 0xc9, 0x9d, 0x9c, 0xdb, 0xd1, 0x14, 0x5e, 0x46, 0x33, 0x2a, 0x94, 0xe4, 0xb6,
	0x51, 0xe8, 0x25, 0x7e, 0x88, 0x96, 0x54, 0x28, 0x5d, 0xe3, 0x69, 0x02, 0x22, 0x38, 0x9c, 0x92,
	0xba, 0x61, 0xab, 0x2a, 0x94, 0xce, 0x71, 0x10, 0x1f, 0x18, 0x10, 0x7f, 0x8e, 0xea, 0x7e, 0x20,
	0xe9, 0x30, 0x04, 0xf7, 0x64, 0x0c, 0x62, 0xea, 0xc2, 0x5b, 0x1a, 0xc5, 0x21, 0x48, 0x72, 0xc7,
	0xc8, 0x57, 0x52, 0x76, 0x5f, 0x93, 0xbd, 0x94, 0xc3, 0x2f, 0xd0, 0x5d, 0xa5, 0x51, 0xa9, 0xa8,
	0x92, 0xee, 0x48, 0xf0, 0x71, 0xec, 0x1a, 0xc0, 0x0d, 0x83, 0x28, 0x50, 0x84, 0x34, 0x4b, 0x1b,
	0x73, 0xf6, 0xea, 0x99, 0xe6, 0xa5, 0x96, 0xbc, 0xd1, 0xdb, 0xbe, 0x16, 0xe8, 0x4b, 0x66, 0x61,
	0x87, 0x54, 0x06, 0x9e, 0x1b, 0x81, 0x12, 0x81, 0x27, 0xc9, 0x6a, 0x72, 0xc9, 0x94, 0xdc, 0xd6,
	0xdc, 0x5e, 0x42, 0xe1, 0x67, 0x88, 0x64, 0x36, 0xc0, 0x8e, 0xf4, 0x64, 0xfa, 0xb9, 0x59, 0xc3,
	0x98, 0x65, 0x57, 0xe9, 0xa5, 0x74, 0x6a, 0xd9, 0x78, 0x81, 0x6e, 0x5e, 0xc8, 0xba, 0xce, 0xd9,
	0xd9, 0xe3, 0xa6, 0x97, 0x78, 0x05, 0xcd, 0x4d, 0x68, 0x38, 0x86, 0xf4, 0x85, 0x48, 0x36, 0xcf,
	0xcb, 0xcf, 0x4a, 0xad, 0x3f, 0xcb, 0xa8, 0x96, 0x55, 0x31, 0x7d, 0xc3, 0x9e, 0xa2, 0x59, 0x3d,
	0x6d, 0xc6, 0xbe, 0xb2, 0x45, 0xda, 0xf9, 0x3f, 0x54, 0xdb, 0x86, 0x88, 0x2b, 0xb0, 0xbb, 0xce,
	0x80, 0xfb, 0x60, 0x1b, 0x15, 0xfe, 0xea, 0xbd, 0xde, 0x2b, 0x1b, 0xab, 0x7a, 0xd1, 0xaa, 0xeb,
	0x64, 0x6d, 0x78, 0xbe, 0x27, 0x37, 0xd1, 0x9c, 0xe9, 0x47, 0xf3, 0x3c, 0x55, 0xb6, 0xee, 0x14,
	0x6c, 0xf6, 0xa6, 0xce, 0x7e, 0xdf, 0x49, 0x66, 0xda, 0x4e, 0x54, 0x78, 0x1b, 0x2d, 0x99, 0x85,
	0x7f, 0x16, 0x6c, 0xd6, 0x18, 0xae, 0xbe, 0x6f, 0xe8, 0xe7, 0xf1, 0x6a, 0x89, 0x45, 0x1e, 0xf2,
	0xcd, 0x15, 0x63, 0x30, 0x67, 0x1c, 0xb5, 0x0a, 0x8e, 0xf6, 0x3b, 0x03, 0xe3, 0x4b, 0x4f, 0x84,
	0x63, 0x44, 0x1d, 0xfd, 0x87, 0x7e, 0xe9, 0xa8, 0xac, 0xa3, 0x4a, 0xd2, 0x23, 0x1e, 0x1f, 0xb3,
	0xe4, 0x29, 0x9c, 0xb3, 0x91, 0x81, 0x76, 0x34, 0xf2, 0xb8, 0x8f, 0x6e, 0x5e, 0x78, 0xee, 0xf0,
	0x3a, 0x5a, 0xeb, 0xee, 0x3a, 0x3b, 0x5f, 0x1f, 0xf4, 0x6c, 0xd7, 0xee, 0x3a, 0x6e, 0x6f, 0xf0,
	0x72, 0x77, 0xd0, 0x73, 0x77, 0x07, 0x07, 0x9d, 0xfe, 0x6e, 0x77, 0xf9, 0x1a, 0xae, 0x23, 0x7c,
	0x4e, 0xb0, 0xf7, 0xad, 0xb3, 0xdf, 0x5f, 0x2e, 0x6d, 0xfd, 0x5e, 0x42, 0x33, 0x76, 0xd7, 0xc1,
	0x13, 0x54, 0x29, 0x78, 0xc5, 0xf7, 0xaf, 0x78, 0x5d, 0xd3, 0xf1, 0x6c, 0xac, 0x5f, 0xc9, 0x27,
	0x85, 0x6f, 0x3d, 0xfc, 0xf1, 0xd7, 0xbf, 0x7e, 0x29, 0x37, 0x5b, 0x6b, 0xd6, 0xe4, 0x53, 0xeb,
	0x4c, 0x6b, 0xd9, 0x5d, 0xc7, 0xca, 0xf4, 0xcf, 0x4b, 0x8f, 0xf1, 0x10, 0xcd, 0x27, 0x2d, 0x83,
	0x57, 0xaf, 0x7c, 0x0c, 0x1a, 0x8d, 0xcb, 0xa8, 0x34, 0xd0, 0x03, 0x13, 0x68, 0xad, 0x55, 0xbf,
	0x24, 0x50, 0xc7, 0xf7, 0x9f, 0x97, 0x1e, 0x6f, 0xbb, 0x3f, 0x77, 0x06, 0x76, 0x1f, 0x5d, 0xf7,
	0xe1, 0x90, 0x8e, 0x43, 0x85, 0x3b, 0x08, 0x77, 0x58, 0x13, 0x84, 0xe0, 0xa2, 0x29, 0x52, 0x3f,
	0x6d, 0xfc, 0x04, 0x7d, 0xdc, 0x78, 0xf4, 0x81, 0xe5, 0xc3, 0x61, 0xc0, 0x82, 0xe4, 0xc3, 0xa8,
	0xf8, 0xf1, 0xd7, 0xd3, 0xf2, 0x2c, 0xea, 0x77, 0x8b, 0x45, 0x6a, 0x38, 0x6f, 0xbe, 0x9a, 0x3e,
	0xfb, 0x27, 0x00, 0x00, 0xff, 0xff, 0x3f, 0xee, 0x55, 0xca, 0x2e, 0x0a, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

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
	cc grpc.ClientConnInterface
}

func NewRDSClient(cc grpc.ClientConnInterface) RDSClient {
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
