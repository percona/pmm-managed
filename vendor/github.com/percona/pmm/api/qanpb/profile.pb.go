// Code generated by protoc-gen-go. DO NOT EDIT.
// source: qanpb/profile.proto

package qanpb

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

// ReportRequest defines filtering of metrics report for db server or other dimentions.
type ReportRequest struct {
	PeriodStartFrom      string                 `protobuf:"bytes,1,opt,name=period_start_from,json=periodStartFrom,proto3" json:"period_start_from,omitempty"`
	PeriodStartTo        string                 `protobuf:"bytes,2,opt,name=period_start_to,json=periodStartTo,proto3" json:"period_start_to,omitempty"`
	Keyword              string                 `protobuf:"bytes,3,opt,name=keyword,proto3" json:"keyword,omitempty"`
	FirstSeen            bool                   `protobuf:"varint,4,opt,name=first_seen,json=firstSeen,proto3" json:"first_seen,omitempty"`
	GroupBy              string                 `protobuf:"bytes,5,opt,name=group_by,json=groupBy,proto3" json:"group_by,omitempty"`
	Labels               []*ReportMapFieldEntry `protobuf:"bytes,6,rep,name=labels,proto3" json:"labels,omitempty"`
	IncludeOnlyFields    []string               `protobuf:"bytes,7,rep,name=include_only_fields,json=includeOnlyFields,proto3" json:"include_only_fields,omitempty"`
	OrderBy              string                 `protobuf:"bytes,8,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Offset               uint32                 `protobuf:"varint,9,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                uint32                 `protobuf:"varint,10,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *ReportRequest) Reset()         { *m = ReportRequest{} }
func (m *ReportRequest) String() string { return proto.CompactTextString(m) }
func (*ReportRequest) ProtoMessage()    {}
func (*ReportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_a87bd6affa0ecd70, []int{0}
}
func (m *ReportRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportRequest.Unmarshal(m, b)
}
func (m *ReportRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportRequest.Marshal(b, m, deterministic)
}
func (dst *ReportRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportRequest.Merge(dst, src)
}
func (m *ReportRequest) XXX_Size() int {
	return xxx_messageInfo_ReportRequest.Size(m)
}
func (m *ReportRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReportRequest proto.InternalMessageInfo

func (m *ReportRequest) GetPeriodStartFrom() string {
	if m != nil {
		return m.PeriodStartFrom
	}
	return ""
}

func (m *ReportRequest) GetPeriodStartTo() string {
	if m != nil {
		return m.PeriodStartTo
	}
	return ""
}

func (m *ReportRequest) GetKeyword() string {
	if m != nil {
		return m.Keyword
	}
	return ""
}

func (m *ReportRequest) GetFirstSeen() bool {
	if m != nil {
		return m.FirstSeen
	}
	return false
}

func (m *ReportRequest) GetGroupBy() string {
	if m != nil {
		return m.GroupBy
	}
	return ""
}

func (m *ReportRequest) GetLabels() []*ReportMapFieldEntry {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *ReportRequest) GetIncludeOnlyFields() []string {
	if m != nil {
		return m.IncludeOnlyFields
	}
	return nil
}

func (m *ReportRequest) GetOrderBy() string {
	if m != nil {
		return m.OrderBy
	}
	return ""
}

func (m *ReportRequest) GetOffset() uint32 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *ReportRequest) GetLimit() uint32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

// ReportMapFieldEntry allows to pass labels/dimentions in form like {"d_server": ["db1", "db2"...]}.
type ReportMapFieldEntry struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []string `protobuf:"bytes,2,rep,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportMapFieldEntry) Reset()         { *m = ReportMapFieldEntry{} }
func (m *ReportMapFieldEntry) String() string { return proto.CompactTextString(m) }
func (*ReportMapFieldEntry) ProtoMessage()    {}
func (*ReportMapFieldEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_a87bd6affa0ecd70, []int{1}
}
func (m *ReportMapFieldEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportMapFieldEntry.Unmarshal(m, b)
}
func (m *ReportMapFieldEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportMapFieldEntry.Marshal(b, m, deterministic)
}
func (dst *ReportMapFieldEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportMapFieldEntry.Merge(dst, src)
}
func (m *ReportMapFieldEntry) XXX_Size() int {
	return xxx_messageInfo_ReportMapFieldEntry.Size(m)
}
func (m *ReportMapFieldEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportMapFieldEntry.DiscardUnknown(m)
}

var xxx_messageInfo_ReportMapFieldEntry proto.InternalMessageInfo

func (m *ReportMapFieldEntry) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *ReportMapFieldEntry) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

// ReportReply is list of reports per quieryids, hosts etc.
type ReportReply struct {
	Rows                 []*ProfileRow `protobuf:"bytes,1,rep,name=rows,proto3" json:"rows,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ReportReply) Reset()         { *m = ReportReply{} }
func (m *ReportReply) String() string { return proto.CompactTextString(m) }
func (*ReportReply) ProtoMessage()    {}
func (*ReportReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_a87bd6affa0ecd70, []int{2}
}
func (m *ReportReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportReply.Unmarshal(m, b)
}
func (m *ReportReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportReply.Marshal(b, m, deterministic)
}
func (dst *ReportReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportReply.Merge(dst, src)
}
func (m *ReportReply) XXX_Size() int {
	return xxx_messageInfo_ReportReply.Size(m)
}
func (m *ReportReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportReply.DiscardUnknown(m)
}

var xxx_messageInfo_ReportReply proto.InternalMessageInfo

func (m *ReportReply) GetRows() []*ProfileRow {
	if m != nil {
		return m.Rows
	}
	return nil
}

// ProfileRow define metrics for selected dimention.
type ProfileRow struct {
	Rank                 uint32   `protobuf:"varint,1,opt,name=rank,proto3" json:"rank,omitempty"`
	Percentage           float32  `protobuf:"fixed32,2,opt,name=percentage,proto3" json:"percentage,omitempty"`
	Dimension            string   `protobuf:"bytes,3,opt,name=dimension,proto3" json:"dimension,omitempty"`
	RowNumber            float32  `protobuf:"fixed32,4,opt,name=row_number,json=rowNumber,proto3" json:"row_number,omitempty"`
	DServers             string   `protobuf:"bytes,5,opt,name=d_servers,json=dServers,proto3" json:"d_servers,omitempty"`
	DDatabases           string   `protobuf:"bytes,6,opt,name=d_databases,json=dDatabases,proto3" json:"d_databases,omitempty"`
	DSchemas             string   `protobuf:"bytes,7,opt,name=d_schemas,json=dSchemas,proto3" json:"d_schemas,omitempty"`
	DUsernames           string   `protobuf:"bytes,8,opt,name=d_usernames,json=dUsernames,proto3" json:"d_usernames,omitempty"`
	DClientHosts         string   `protobuf:"bytes,9,opt,name=d_client_hosts,json=dClientHosts,proto3" json:"d_client_hosts,omitempty"`
	FirstSeen            string   `protobuf:"bytes,10,opt,name=first_seen,json=firstSeen,proto3" json:"first_seen,omitempty"`
	Qps                  float32  `protobuf:"fixed32,11,opt,name=qps,proto3" json:"qps,omitempty"`
	Load                 float32  `protobuf:"fixed32,12,opt,name=load,proto3" json:"load,omitempty"`
	Fingerprint          string   `protobuf:"bytes,13,opt,name=fingerprint,proto3" json:"fingerprint,omitempty"`
	Stats                *Stats   `protobuf:"bytes,14,opt,name=stats,proto3" json:"stats,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProfileRow) Reset()         { *m = ProfileRow{} }
func (m *ProfileRow) String() string { return proto.CompactTextString(m) }
func (*ProfileRow) ProtoMessage()    {}
func (*ProfileRow) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_a87bd6affa0ecd70, []int{3}
}
func (m *ProfileRow) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProfileRow.Unmarshal(m, b)
}
func (m *ProfileRow) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProfileRow.Marshal(b, m, deterministic)
}
func (dst *ProfileRow) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProfileRow.Merge(dst, src)
}
func (m *ProfileRow) XXX_Size() int {
	return xxx_messageInfo_ProfileRow.Size(m)
}
func (m *ProfileRow) XXX_DiscardUnknown() {
	xxx_messageInfo_ProfileRow.DiscardUnknown(m)
}

var xxx_messageInfo_ProfileRow proto.InternalMessageInfo

func (m *ProfileRow) GetRank() uint32 {
	if m != nil {
		return m.Rank
	}
	return 0
}

func (m *ProfileRow) GetPercentage() float32 {
	if m != nil {
		return m.Percentage
	}
	return 0
}

func (m *ProfileRow) GetDimension() string {
	if m != nil {
		return m.Dimension
	}
	return ""
}

func (m *ProfileRow) GetRowNumber() float32 {
	if m != nil {
		return m.RowNumber
	}
	return 0
}

func (m *ProfileRow) GetDServers() string {
	if m != nil {
		return m.DServers
	}
	return ""
}

func (m *ProfileRow) GetDDatabases() string {
	if m != nil {
		return m.DDatabases
	}
	return ""
}

func (m *ProfileRow) GetDSchemas() string {
	if m != nil {
		return m.DSchemas
	}
	return ""
}

func (m *ProfileRow) GetDUsernames() string {
	if m != nil {
		return m.DUsernames
	}
	return ""
}

func (m *ProfileRow) GetDClientHosts() string {
	if m != nil {
		return m.DClientHosts
	}
	return ""
}

func (m *ProfileRow) GetFirstSeen() string {
	if m != nil {
		return m.FirstSeen
	}
	return ""
}

func (m *ProfileRow) GetQps() float32 {
	if m != nil {
		return m.Qps
	}
	return 0
}

func (m *ProfileRow) GetLoad() float32 {
	if m != nil {
		return m.Load
	}
	return 0
}

func (m *ProfileRow) GetFingerprint() string {
	if m != nil {
		return m.Fingerprint
	}
	return ""
}

func (m *ProfileRow) GetStats() *Stats {
	if m != nil {
		return m.Stats
	}
	return nil
}

// Stats metrics.
type Stats struct {
	NumQueries           float32  `protobuf:"fixed32,1,opt,name=num_queries,json=numQueries,proto3" json:"num_queries,omitempty"`
	MQueryTimeSum        float32  `protobuf:"fixed32,2,opt,name=m_query_time_sum,json=mQueryTimeSum,proto3" json:"m_query_time_sum,omitempty"`
	MQueryTimeMin        float32  `protobuf:"fixed32,3,opt,name=m_query_time_min,json=mQueryTimeMin,proto3" json:"m_query_time_min,omitempty"`
	MQueryTimeMax        float32  `protobuf:"fixed32,4,opt,name=m_query_time_max,json=mQueryTimeMax,proto3" json:"m_query_time_max,omitempty"`
	MQueryTimeP99        float32  `protobuf:"fixed32,5,opt,name=m_query_time_p99,json=mQueryTimeP99,proto3" json:"m_query_time_p99,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stats) Reset()         { *m = Stats{} }
func (m *Stats) String() string { return proto.CompactTextString(m) }
func (*Stats) ProtoMessage()    {}
func (*Stats) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_a87bd6affa0ecd70, []int{4}
}
func (m *Stats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stats.Unmarshal(m, b)
}
func (m *Stats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stats.Marshal(b, m, deterministic)
}
func (dst *Stats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stats.Merge(dst, src)
}
func (m *Stats) XXX_Size() int {
	return xxx_messageInfo_Stats.Size(m)
}
func (m *Stats) XXX_DiscardUnknown() {
	xxx_messageInfo_Stats.DiscardUnknown(m)
}

var xxx_messageInfo_Stats proto.InternalMessageInfo

func (m *Stats) GetNumQueries() float32 {
	if m != nil {
		return m.NumQueries
	}
	return 0
}

func (m *Stats) GetMQueryTimeSum() float32 {
	if m != nil {
		return m.MQueryTimeSum
	}
	return 0
}

func (m *Stats) GetMQueryTimeMin() float32 {
	if m != nil {
		return m.MQueryTimeMin
	}
	return 0
}

func (m *Stats) GetMQueryTimeMax() float32 {
	if m != nil {
		return m.MQueryTimeMax
	}
	return 0
}

func (m *Stats) GetMQueryTimeP99() float32 {
	if m != nil {
		return m.MQueryTimeP99
	}
	return 0
}

func init() {
	proto.RegisterType((*ReportRequest)(nil), "qan.ReportRequest")
	proto.RegisterType((*ReportMapFieldEntry)(nil), "qan.ReportMapFieldEntry")
	proto.RegisterType((*ReportReply)(nil), "qan.ReportReply")
	proto.RegisterType((*ProfileRow)(nil), "qan.ProfileRow")
	proto.RegisterType((*Stats)(nil), "qan.Stats")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProfileClient is the client API for Profile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProfileClient interface {
	// GetReport returns list of metrics group by queryid or other dimentions.
	GetReport(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportReply, error)
}

type profileClient struct {
	cc *grpc.ClientConn
}

func NewProfileClient(cc *grpc.ClientConn) ProfileClient {
	return &profileClient{cc}
}

func (c *profileClient) GetReport(ctx context.Context, in *ReportRequest, opts ...grpc.CallOption) (*ReportReply, error) {
	out := new(ReportReply)
	err := c.cc.Invoke(ctx, "/qan.Profile/GetReport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileServer is the server API for Profile service.
type ProfileServer interface {
	// GetReport returns list of metrics group by queryid or other dimentions.
	GetReport(context.Context, *ReportRequest) (*ReportReply, error)
}

func RegisterProfileServer(s *grpc.Server, srv ProfileServer) {
	s.RegisterService(&_Profile_serviceDesc, srv)
}

func _Profile_GetReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServer).GetReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/qan.Profile/GetReport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServer).GetReport(ctx, req.(*ReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Profile_serviceDesc = grpc.ServiceDesc{
	ServiceName: "qan.Profile",
	HandlerType: (*ProfileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetReport",
			Handler:    _Profile_GetReport_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "qanpb/profile.proto",
}

func init() { proto.RegisterFile("qanpb/profile.proto", fileDescriptor_profile_a87bd6affa0ecd70) }

var fileDescriptor_profile_a87bd6affa0ecd70 = []byte{
	// 726 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x94, 0x4f, 0x6e, 0xe3, 0x36,
	0x14, 0xc6, 0x21, 0x3b, 0xfe, 0xa3, 0xe7, 0x38, 0x71, 0x98, 0x22, 0x60, 0xd3, 0xb4, 0x15, 0xdc,
	0xa2, 0x35, 0xb2, 0xb0, 0x5b, 0x77, 0xe5, 0x02, 0xdd, 0xa4, 0x6d, 0xda, 0x4d, 0x9a, 0x54, 0x4e,
	0x37, 0xd9, 0x08, 0xb4, 0xf5, 0xec, 0x10, 0x91, 0x48, 0x99, 0xa4, 0xe2, 0x68, 0xdb, 0x2b, 0xcc,
	0x25, 0xe6, 0x30, 0xb3, 0x9b, 0x1b, 0x0c, 0xe6, 0x20, 0x03, 0x91, 0xf2, 0xc4, 0x9e, 0xcc, 0x8e,
	0xdf, 0xf7, 0xbe, 0x27, 0xea, 0xe9, 0x47, 0x11, 0x8e, 0x57, 0x4c, 0x64, 0xb3, 0x51, 0xa6, 0xe4,
	0x82, 0x27, 0x38, 0xcc, 0x94, 0x34, 0x92, 0xd4, 0x57, 0x4c, 0x9c, 0x9e, 0x2d, 0xa5, 0x5c, 0x26,
	0x38, 0x62, 0x19, 0x1f, 0x31, 0x21, 0xa4, 0x61, 0x86, 0x4b, 0xa1, 0x5d, 0xa4, 0xff, 0xae, 0x06,
	0xdd, 0x10, 0x33, 0xa9, 0x4c, 0x88, 0xab, 0x1c, 0xb5, 0x21, 0xe7, 0x70, 0x94, 0xa1, 0xe2, 0x32,
	0x8e, 0xb4, 0x61, 0xca, 0x44, 0x0b, 0x25, 0x53, 0xea, 0x05, 0xde, 0xc0, 0x0f, 0x0f, 0x5d, 0x61,
	0x5a, 0xfa, 0x97, 0x4a, 0xa6, 0xe4, 0x07, 0x38, 0xdc, 0xc9, 0x1a, 0x49, 0x6b, 0x36, 0xd9, 0xdd,
	0x4a, 0xde, 0x4a, 0x42, 0xa1, 0xf5, 0x80, 0xc5, 0x5a, 0xaa, 0x98, 0xd6, 0x6d, 0x7d, 0x23, 0xc9,
	0xd7, 0x00, 0x0b, 0xae, 0xb4, 0x89, 0x34, 0xa2, 0xa0, 0x7b, 0x81, 0x37, 0x68, 0x87, 0xbe, 0x75,
	0xa6, 0x88, 0x82, 0x7c, 0x09, 0xed, 0xa5, 0x92, 0x79, 0x16, 0xcd, 0x0a, 0xda, 0x70, 0x9d, 0x56,
	0x5f, 0x14, 0xe4, 0x27, 0x68, 0x26, 0x6c, 0x86, 0x89, 0xa6, 0xcd, 0xa0, 0x3e, 0xe8, 0x8c, 0xe9,
	0x70, 0xc5, 0xc4, 0xd0, 0xcd, 0x72, 0xc5, 0xb2, 0x4b, 0x8e, 0x49, 0xfc, 0xa7, 0x30, 0xaa, 0x08,
	0xab, 0x1c, 0x19, 0xc2, 0x31, 0x17, 0xf3, 0x24, 0x8f, 0x31, 0x92, 0x22, 0x29, 0xa2, 0x45, 0x19,
	0xd1, 0xb4, 0x15, 0xd4, 0x07, 0x7e, 0x78, 0x54, 0x95, 0xae, 0x45, 0x52, 0xd8, 0x5e, 0x5d, 0x6e,
	0x2e, 0x55, 0x8c, 0xaa, 0xdc, 0xbc, 0xed, 0x36, 0xb7, 0xfa, 0xa2, 0x20, 0x27, 0xd0, 0x94, 0x8b,
	0x85, 0x46, 0x43, 0xfd, 0xc0, 0x1b, 0x74, 0xc3, 0x4a, 0x91, 0x2f, 0xa0, 0x91, 0xf0, 0x94, 0x1b,
	0x0a, 0xd6, 0x76, 0xa2, 0xff, 0x1b, 0x1c, 0x7f, 0xe6, 0xbd, 0x48, 0x0f, 0xea, 0x0f, 0x58, 0x54,
	0xdf, 0xb6, 0x5c, 0x96, 0xed, 0x8f, 0x2c, 0xc9, 0x91, 0xd6, 0xec, 0x3b, 0x39, 0xd1, 0x1f, 0x43,
	0x67, 0x83, 0x28, 0x4b, 0x0a, 0xf2, 0x1d, 0xec, 0x29, 0xb9, 0xd6, 0xd4, 0xb3, 0x63, 0x1f, 0xda,
	0xb1, 0x6f, 0x1c, 0xf7, 0x50, 0xae, 0x43, 0x5b, 0xec, 0xbf, 0xae, 0x03, 0x3c, 0x9b, 0x84, 0xc0,
	0x9e, 0x62, 0xe2, 0xc1, 0xee, 0xd5, 0x0d, 0xed, 0x9a, 0x7c, 0x03, 0x90, 0xa1, 0x9a, 0xa3, 0x30,
	0x6c, 0x89, 0x96, 0x5b, 0x2d, 0xdc, 0x72, 0xc8, 0x19, 0xf8, 0x31, 0x4f, 0x51, 0x68, 0x2e, 0x45,
	0x85, 0xed, 0xd9, 0x28, 0xc1, 0x29, 0xb9, 0x8e, 0x44, 0x9e, 0xce, 0x50, 0x59, 0x70, 0xb5, 0xd0,
	0x57, 0x72, 0xfd, 0x8f, 0x35, 0xc8, 0x57, 0xe0, 0xc7, 0x91, 0x46, 0xf5, 0x88, 0x4a, 0x57, 0xe4,
	0xda, 0xf1, 0xd4, 0x69, 0xf2, 0x2d, 0x74, 0xe2, 0x28, 0x66, 0x86, 0xcd, 0x98, 0xc6, 0x92, 0x5f,
	0x59, 0x86, 0xf8, 0x8f, 0x8d, 0x53, 0x75, 0xcf, 0xef, 0x31, 0x65, 0x25, 0x9f, 0xaa, 0xdb, 0x69,
	0xd7, 0x9d, 0x6b, 0x54, 0x82, 0xa5, 0xa8, 0x2b, 0x32, 0x10, 0xff, 0xb7, 0x71, 0xc8, 0xf7, 0x70,
	0x10, 0x47, 0xf3, 0x84, 0xa3, 0x30, 0xd1, 0xbd, 0xd4, 0x46, 0x5b, 0x48, 0x7e, 0xb8, 0x1f, 0xff,
	0x6e, 0xcd, 0xbf, 0x4b, 0xef, 0x93, 0x93, 0x07, 0x6e, 0xbe, 0xe7, 0x93, 0xd7, 0x83, 0xfa, 0x2a,
	0xd3, 0xb4, 0x63, 0x07, 0x2b, 0x97, 0xe5, 0x37, 0x4c, 0x24, 0x8b, 0xe9, 0xbe, 0xb5, 0xec, 0x9a,
	0x04, 0xd0, 0x59, 0x70, 0xb1, 0x44, 0x95, 0x29, 0x2e, 0x0c, 0xed, 0xda, 0xa7, 0x6c, 0x5b, 0x24,
	0x80, 0x86, 0x36, 0xcc, 0x68, 0x7a, 0x10, 0x78, 0x83, 0xce, 0x18, 0x2c, 0xae, 0x69, 0xe9, 0x84,
	0xae, 0xd0, 0x7f, 0xe3, 0x41, 0xc3, 0x1a, 0xe5, 0x64, 0x22, 0x4f, 0xa3, 0x55, 0x8e, 0x8a, 0xa3,
	0xb6, 0xb0, 0x6a, 0x21, 0x88, 0x3c, 0xfd, 0xd7, 0x39, 0xe4, 0x47, 0xe8, 0xb9, 0x72, 0x11, 0x19,
	0x9e, 0x62, 0xa4, 0xf3, 0xb4, 0x02, 0xd7, 0xb5, 0x99, 0xe2, 0x96, 0xa7, 0x38, 0xcd, 0xd3, 0x17,
	0xc1, 0x94, 0x3b, 0x84, 0x3b, 0xc1, 0x2b, 0x2e, 0x5e, 0x06, 0xd9, 0x53, 0x05, 0x73, 0x3b, 0xc8,
	0x9e, 0x5e, 0x04, 0xb3, 0xc9, 0xc4, 0x72, 0xdd, 0x09, 0xde, 0x4c, 0x26, 0xe3, 0x3b, 0x68, 0x55,
	0x07, 0x8f, 0x5c, 0x83, 0xff, 0x17, 0x1a, 0x77, 0x76, 0x09, 0xd9, 0xfa, 0x3f, 0xab, 0xbb, 0xe6,
	0xb4, 0xb7, 0xe3, 0x65, 0x49, 0xd1, 0x3f, 0xfb, 0xff, 0xed, 0xfb, 0x57, 0xb5, 0x93, 0xfe, 0xd1,
	0xe8, 0xf1, 0xe7, 0xd1, 0x8a, 0x89, 0xd1, 0xc7, 0x07, 0xfc, 0xea, 0x9d, 0x5f, 0xb4, 0xee, 0x1a,
	0xf6, 0x9e, 0x9b, 0x35, 0xed, 0xed, 0xf5, 0xcb, 0x87, 0x00, 0x00, 0x00, 0xff, 0xff, 0x65, 0x2d,
	0xbd, 0x4b, 0xf7, 0x04, 0x00, 0x00,
}
