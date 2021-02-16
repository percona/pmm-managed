// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: managementpb/backup/locations.proto

package backupv1beta1

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *PMMServerLocationConfig) Validate() error {
	if this.Path == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Path", fmt.Errorf(`value '%v' must not be an empty string`, this.Path))
	}
	return nil
}
func (this *PMMClientLocationConfig) Validate() error {
	if this.Path == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Path", fmt.Errorf(`value '%v' must not be an empty string`, this.Path))
	}
	return nil
}
func (this *S3LocationConfig) Validate() error {
	if this.Endpoint == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Endpoint", fmt.Errorf(`value '%v' must not be an empty string`, this.Endpoint))
	}
	if this.AccessKey == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("AccessKey", fmt.Errorf(`value '%v' must not be an empty string`, this.AccessKey))
	}
	if this.SecretKey == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("SecretKey", fmt.Errorf(`value '%v' must not be an empty string`, this.SecretKey))
	}
	return nil
}
func (this *Location) Validate() error {
	if oneOfNester, ok := this.GetConfig().(*Location_PmmClientConfig); ok {
		if oneOfNester.PmmClientConfig != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.PmmClientConfig); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("PmmClientConfig", err)
			}
		}
	}
	if oneOfNester, ok := this.GetConfig().(*Location_PmmServerConfig); ok {
		if oneOfNester.PmmServerConfig != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.PmmServerConfig); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("PmmServerConfig", err)
			}
		}
	}
	if oneOfNester, ok := this.GetConfig().(*Location_S3Config); ok {
		if oneOfNester.S3Config != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(oneOfNester.S3Config); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("S3Config", err)
			}
		}
	}
	return nil
}
func (this *ListLocationsRequest) Validate() error {
	return nil
}
func (this *ListLocationsResponse) Validate() error {
	for _, item := range this.Locations {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Locations", err)
			}
		}
	}
	return nil
}
func (this *AddLocationRequest) Validate() error {
	if this.Name == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Name", fmt.Errorf(`value '%v' must not be an empty string`, this.Name))
	}
	if this.PmmClientConfig != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PmmClientConfig); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PmmClientConfig", err)
		}
	}
	if this.PmmServerConfig != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.PmmServerConfig); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("PmmServerConfig", err)
		}
	}
	if this.S3Config != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.S3Config); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("S3Config", err)
		}
	}
	return nil
}
func (this *AddLocationResponse) Validate() error {
	return nil
}
func (this *RemoveLocationRequest) Validate() error {
	return nil
}
func (this *RemoveLocationResponse) Validate() error {
	return nil
}
