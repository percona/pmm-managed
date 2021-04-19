// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: managementpb/checks.proto

package managementpb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *SecurityCheckResult) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *SecurityCheck) Validate() error {
	return nil
}
func (this *ChangeSecurityCheckParams) Validate() error {
	return nil
}
func (this *GetSecurityCheckResultsRequest) Validate() error {
	return nil
}
func (this *GetSecurityCheckResultsResponse) Validate() error {
	for _, item := range this.Results {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Results", err)
			}
		}
	}
	return nil
}
func (this *StartSecurityChecksRequest) Validate() error {
	return nil
}
func (this *StartSecurityChecksResponse) Validate() error {
	return nil
}
func (this *ListSecurityChecksRequest) Validate() error {
	return nil
}
func (this *ListSecurityChecksResponse) Validate() error {
	for _, item := range this.Checks {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Checks", err)
			}
		}
	}
	return nil
}
func (this *ChangeSecurityChecksRequest) Validate() error {
	for _, item := range this.Params {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Params", err)
			}
		}
	}
	return nil
}
func (this *ChangeSecurityChecksResponse) Validate() error {
	return nil
}