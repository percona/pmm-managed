/* 
 * A Bit of Everything
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * OpenAPI spec version: 1.0
 * Contact: none@example.com
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 */

package abe

import (
	"time"
)

// Intentionaly complicated message type to cover many features of Protobuf.
type ExamplepbABitOfEverything struct {

	SingleNested ABitOfEverythingNested `json:"single_nested,omitempty"`

	Uuid string `json:"uuid"`

	Nested []ABitOfEverythingNested `json:"nested,omitempty"`

	// Float value field
	FloatValue float32 `json:"float_value,omitempty"`

	DoubleValue float64 `json:"double_value"`

	Int64Value string `json:"int64_value"`

	Uint64Value string `json:"uint64_value,omitempty"`

	Int32Value int32 `json:"int32_value,omitempty"`

	Fixed64Value string `json:"fixed64_value,omitempty"`

	Fixed32Value int64 `json:"fixed32_value,omitempty"`

	BoolValue bool `json:"bool_value,omitempty"`

	StringValue string `json:"string_value,omitempty"`

	BytesValue string `json:"bytes_value,omitempty"`

	Uint32Value int64 `json:"uint32_value,omitempty"`

	EnumValue ExamplepbNumericEnum `json:"enum_value,omitempty"`

	PathEnumValue PathenumPathEnum `json:"path_enum_value,omitempty"`

	NestedPathEnumValue MessagePathEnumNestedPathEnum `json:"nested_path_enum_value,omitempty"`

	Sfixed32Value int32 `json:"sfixed32_value,omitempty"`

	Sfixed64Value string `json:"sfixed64_value,omitempty"`

	Sint32Value int32 `json:"sint32_value,omitempty"`

	Sint64Value string `json:"sint64_value,omitempty"`

	RepeatedStringValue []string `json:"repeated_string_value,omitempty"`

	OneofEmpty interface{} `json:"oneof_empty,omitempty"`

	OneofString string `json:"oneof_string,omitempty"`

	MapValue map[string]ExamplepbNumericEnum `json:"map_value,omitempty"`

	MappedStringValue map[string]string `json:"mapped_string_value,omitempty"`

	MappedNestedValue map[string]ABitOfEverythingNested `json:"mapped_nested_value,omitempty"`

	NonConventionalNameValue string `json:"nonConventionalNameValue,omitempty"`

	TimestampValue time.Time `json:"timestamp_value,omitempty"`

	RepeatedEnumValue []ExamplepbNumericEnum `json:"repeated_enum_value,omitempty"`
}
