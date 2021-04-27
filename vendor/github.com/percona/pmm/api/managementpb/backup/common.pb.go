// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: managementpb/backup/common.proto

package backupv1beta1

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// DataModel is a model used for performing a backup.
type DataModel int32

const (
	DataModel_DATA_MODEL_INVALID DataModel = 0
	DataModel_PHYSICAL           DataModel = 1
	DataModel_LOGICAL            DataModel = 2
)

// Enum value maps for DataModel.
var (
	DataModel_name = map[int32]string{
		0: "DATA_MODEL_INVALID",
		1: "PHYSICAL",
		2: "LOGICAL",
	}
	DataModel_value = map[string]int32{
		"DATA_MODEL_INVALID": 0,
		"PHYSICAL":           1,
		"LOGICAL":            2,
	}
)

func (x DataModel) Enum() *DataModel {
	p := new(DataModel)
	*p = x
	return p
}

func (x DataModel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DataModel) Descriptor() protoreflect.EnumDescriptor {
	return file_managementpb_backup_common_proto_enumTypes[0].Descriptor()
}

func (DataModel) Type() protoreflect.EnumType {
	return &file_managementpb_backup_common_proto_enumTypes[0]
}

func (x DataModel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DataModel.Descriptor instead.
func (DataModel) EnumDescriptor() ([]byte, []int) {
	return file_managementpb_backup_common_proto_rawDescGZIP(), []int{0}
}

var File_managementpb_backup_common_proto protoreflect.FileDescriptor

var file_managementpb_backup_common_proto_rawDesc = []byte{
	0x0a, 0x20, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2f, 0x62,
	0x61, 0x63, 0x6b, 0x75, 0x70, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x62, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x2a, 0x3e, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12,
	0x16, 0x0a, 0x12, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x4c, 0x5f, 0x49, 0x4e,
	0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x48, 0x59, 0x53, 0x49,
	0x43, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x4c, 0x4f, 0x47, 0x49, 0x43, 0x41, 0x4c,
	0x10, 0x02, 0x42, 0x27, 0x5a, 0x25, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x3b, 0x62, 0x61,
	0x63, 0x6b, 0x75, 0x70, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_managementpb_backup_common_proto_rawDescOnce sync.Once
	file_managementpb_backup_common_proto_rawDescData = file_managementpb_backup_common_proto_rawDesc
)

func file_managementpb_backup_common_proto_rawDescGZIP() []byte {
	file_managementpb_backup_common_proto_rawDescOnce.Do(func() {
		file_managementpb_backup_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_managementpb_backup_common_proto_rawDescData)
	})
	return file_managementpb_backup_common_proto_rawDescData
}

var file_managementpb_backup_common_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_managementpb_backup_common_proto_goTypes = []interface{}{
	(DataModel)(0), // 0: backup.v1beta1.DataModel
}
var file_managementpb_backup_common_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_managementpb_backup_common_proto_init() }
func file_managementpb_backup_common_proto_init() {
	if File_managementpb_backup_common_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_managementpb_backup_common_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_managementpb_backup_common_proto_goTypes,
		DependencyIndexes: file_managementpb_backup_common_proto_depIdxs,
		EnumInfos:         file_managementpb_backup_common_proto_enumTypes,
	}.Build()
	File_managementpb_backup_common_proto = out.File
	file_managementpb_backup_common_proto_rawDesc = nil
	file_managementpb_backup_common_proto_goTypes = nil
	file_managementpb_backup_common_proto_depIdxs = nil
}