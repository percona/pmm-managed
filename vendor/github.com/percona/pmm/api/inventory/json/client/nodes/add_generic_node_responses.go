// Code generated by go-swagger; DO NOT EDIT.

package nodes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// AddGenericNodeReader is a Reader for the AddGenericNode structure.
type AddGenericNodeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddGenericNodeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewAddGenericNodeOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewAddGenericNodeOK creates a AddGenericNodeOK with default headers values
func NewAddGenericNodeOK() *AddGenericNodeOK {
	return &AddGenericNodeOK{}
}

/*AddGenericNodeOK handles this case with default header values.

A successful response.
*/
type AddGenericNodeOK struct {
	Payload *AddGenericNodeOKBody
}

func (o *AddGenericNodeOK) Error() string {
	return fmt.Sprintf("[POST /v1/inventory/Nodes/AddGeneric][%d] addGenericNodeOK  %+v", 200, o.Payload)
}

func (o *AddGenericNodeOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(AddGenericNodeOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*AddGenericNodeBody add generic node body
swagger:model AddGenericNodeBody
*/
type AddGenericNodeBody struct {

	// Custom user-assigned labels. Keys must start with "_".
	CustomLabels map[string]string `json:"custom_labels,omitempty"`

	// Linux distribution (if any).
	Distro string `json:"distro,omitempty"`

	// Linux distribution version (if any).
	DistroVersion string `json:"distro_version,omitempty"`

	// Linux machine-id. Must be unique across all Generic Nodes if specified.
	MachineID string `json:"machine_id,omitempty"`

	// Unique across all Nodes user-defined name.
	NodeName string `json:"node_name,omitempty"`
}

// Validate validates this add generic node body
func (o *AddGenericNodeBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *AddGenericNodeBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AddGenericNodeBody) UnmarshalBinary(b []byte) error {
	var res AddGenericNodeBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*AddGenericNodeOKBody add generic node o k body
swagger:model AddGenericNodeOKBody
*/
type AddGenericNodeOKBody struct {

	// generic
	Generic *AddGenericNodeOKBodyGeneric `json:"generic,omitempty"`
}

// Validate validates this add generic node o k body
func (o *AddGenericNodeOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateGeneric(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *AddGenericNodeOKBody) validateGeneric(formats strfmt.Registry) error {

	if swag.IsZero(o.Generic) { // not required
		return nil
	}

	if o.Generic != nil {
		if err := o.Generic.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("addGenericNodeOK" + "." + "generic")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *AddGenericNodeOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AddGenericNodeOKBody) UnmarshalBinary(b []byte) error {
	var res AddGenericNodeOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*AddGenericNodeOKBodyGeneric GenericNode represents a bare metal server or virtual machine.
swagger:model AddGenericNodeOKBodyGeneric
*/
type AddGenericNodeOKBodyGeneric struct {

	// Custom user-assigned labels. Keys must start with "_". Can be changed.
	CustomLabels map[string]string `json:"custom_labels,omitempty"`

	// Linux distribution (if any). Can be changed.
	Distro string `json:"distro,omitempty"`

	// Linux distribution version (if any). Can be changed.
	DistroVersion string `json:"distro_version,omitempty"`

	// Linux machine-id. Can't be changed. Must be unique across all Generic Nodes if specified.
	MachineID string `json:"machine_id,omitempty"`

	// Unique randomly generated instance identifier, can't be changed.
	NodeID string `json:"node_id,omitempty"`

	// Unique across all Nodes user-defined name, can be changed.
	NodeName string `json:"node_name,omitempty"`
}

// Validate validates this add generic node o k body generic
func (o *AddGenericNodeOKBodyGeneric) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *AddGenericNodeOKBodyGeneric) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AddGenericNodeOKBodyGeneric) UnmarshalBinary(b []byte) error {
	var res AddGenericNodeOKBodyGeneric
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
