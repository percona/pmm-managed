// Code generated by protoc-gen-go. DO NOT EDIT.
// source: qan/qan.proto

package qanpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// ExampleFormat is formant of query example: real or query without values
type ExampleFormat int32

const (
	ExampleFormat_EXAMPLE_FORMAT_INVALID ExampleFormat = 0
	ExampleFormat_EXAMPLE                ExampleFormat = 1
	ExampleFormat_FINGERPRINT            ExampleFormat = 2
)

var ExampleFormat_name = map[int32]string{
	0: "EXAMPLE_FORMAT_INVALID",
	1: "EXAMPLE",
	2: "FINGERPRINT",
}
var ExampleFormat_value = map[string]int32{
	"EXAMPLE_FORMAT_INVALID": 0,
	"EXAMPLE":                1,
	"FINGERPRINT":            2,
}

func (x ExampleFormat) String() string {
	return proto.EnumName(ExampleFormat_name, int32(x))
}
func (ExampleFormat) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_qan_3fc45e23a3421d5b, []int{0}
}

// ExampleType is a type of query example selected for this query class in given period of time.
type ExampleType int32

const (
	ExampleType_EXAMPLE_TYPE_INVALID ExampleType = 0
	ExampleType_RANDOM               ExampleType = 1
	ExampleType_SLOWEST              ExampleType = 2
	ExampleType_FASTEST              ExampleType = 3
	ExampleType_WITH_ERROR           ExampleType = 4
)

var ExampleType_name = map[int32]string{
	0: "EXAMPLE_TYPE_INVALID",
	1: "RANDOM",
	2: "SLOWEST",
	3: "FASTEST",
	4: "WITH_ERROR",
}
var ExampleType_value = map[string]int32{
	"EXAMPLE_TYPE_INVALID": 0,
	"RANDOM":               1,
	"SLOWEST":              2,
	"FASTEST":              3,
	"WITH_ERROR":           4,
}

func (x ExampleType) String() string {
	return proto.EnumName(ExampleType_name, int32(x))
}
func (ExampleType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_qan_3fc45e23a3421d5b, []int{1}
}

func init() {
	proto.RegisterEnum("qan.ExampleFormat", ExampleFormat_name, ExampleFormat_value)
	proto.RegisterEnum("qan.ExampleType", ExampleType_name, ExampleType_value)
}

func init() { proto.RegisterFile("qan/qan.proto", fileDescriptor_qan_3fc45e23a3421d5b) }

var fileDescriptor_qan_3fc45e23a3421d5b = []byte{
	// 197 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0xce, 0xc1, 0x6a, 0x83, 0x30,
	0x1c, 0xc7, 0xf1, 0xa9, 0x9b, 0xc2, 0x5f, 0xdc, 0x42, 0x18, 0x63, 0xec, 0x11, 0x3c, 0x6c, 0x87,
	0x3d, 0x41, 0x86, 0xc9, 0x16, 0xd0, 0x44, 0x62, 0x98, 0x5b, 0xa1, 0x48, 0x04, 0x6f, 0x35, 0x26,
	0xe2, 0xa1, 0x7d, 0xfb, 0x62, 0xb1, 0xd0, 0xe3, 0x97, 0x1f, 0x7c, 0xf8, 0x41, 0xe6, 0x8d, 0xfd,
	0xf0, 0xc6, 0xbe, 0xbb, 0x79, 0x5a, 0x26, 0x1c, 0x79, 0x63, 0x73, 0x0e, 0x19, 0x3d, 0x9a, 0xd1,
	0x1d, 0x06, 0x36, 0xcd, 0xa3, 0x59, 0xf0, 0x1b, 0xbc, 0xd0, 0x3f, 0x52, 0xd5, 0x25, 0xed, 0x98,
	0x54, 0x15, 0xd1, 0x1d, 0x17, 0xbf, 0xa4, 0xe4, 0x05, 0xba, 0xc3, 0x29, 0x24, 0xdb, 0x86, 0x02,
	0xfc, 0x04, 0x29, 0xe3, 0xe2, 0x9b, 0xaa, 0x5a, 0x71, 0xa1, 0x51, 0x98, 0xef, 0x21, 0xdd, 0x28,
	0x7d, 0x72, 0x03, 0x7e, 0x85, 0xe7, 0x2b, 0xa4, 0xff, 0x6b, 0x7a, 0xc3, 0x00, 0xc4, 0x8a, 0x88,
	0x42, 0x56, 0x28, 0x58, 0xc9, 0xa6, 0x94, 0x2d, 0x6d, 0x34, 0x0a, 0xd7, 0x60, 0xa4, 0xd1, 0x6b,
	0x44, 0xf8, 0x11, 0xa0, 0xe5, 0xfa, 0xa7, 0xa3, 0x4a, 0x49, 0x85, 0xee, 0xbf, 0x92, 0xdd, 0x83,
	0x37, 0xd6, 0xf5, 0x7d, 0x7c, 0xb9, 0xff, 0x79, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x26, 0x7e, 0xa1,
	0x89, 0xcf, 0x00, 0x00, 0x00,
}
