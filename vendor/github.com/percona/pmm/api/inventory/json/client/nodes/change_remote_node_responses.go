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

// ChangeRemoteNodeReader is a Reader for the ChangeRemoteNode structure.
type ChangeRemoteNodeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ChangeRemoteNodeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewChangeRemoteNodeOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewChangeRemoteNodeOK creates a ChangeRemoteNodeOK with default headers values
func NewChangeRemoteNodeOK() *ChangeRemoteNodeOK {
	return &ChangeRemoteNodeOK{}
}

/*ChangeRemoteNodeOK handles this case with default header values.

A successful response.
*/
type ChangeRemoteNodeOK struct {
	Payload *ChangeRemoteNodeOKBody
}

func (o *ChangeRemoteNodeOK) Error() string {
	return fmt.Sprintf("[POST /v1/inventory/Nodes/ChangeRemote][%d] changeRemoteNodeOK  %+v", 200, o.Payload)
}

func (o *ChangeRemoteNodeOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(ChangeRemoteNodeOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*ChangeRemoteNodeBody change remote node body
swagger:model ChangeRemoteNodeBody
*/
type ChangeRemoteNodeBody struct {

	// Unique randomly generated instance identifier.
	NodeID string `json:"node_id,omitempty"`

	// Unique across all Nodes user-defined name.
	NodeName string `json:"node_name,omitempty"`
}

// Validate validates this change remote node body
func (o *ChangeRemoteNodeBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ChangeRemoteNodeBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ChangeRemoteNodeBody) UnmarshalBinary(b []byte) error {
	var res ChangeRemoteNodeBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*ChangeRemoteNodeOKBody change remote node o k body
swagger:model ChangeRemoteNodeOKBody
*/
type ChangeRemoteNodeOKBody struct {

	// remote
	Remote *ChangeRemoteNodeOKBodyRemote `json:"remote,omitempty"`
}

// Validate validates this change remote node o k body
func (o *ChangeRemoteNodeOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateRemote(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ChangeRemoteNodeOKBody) validateRemote(formats strfmt.Registry) error {

	if swag.IsZero(o.Remote) { // not required
		return nil
	}

	if o.Remote != nil {
		if err := o.Remote.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("changeRemoteNodeOK" + "." + "remote")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *ChangeRemoteNodeOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ChangeRemoteNodeOKBody) UnmarshalBinary(b []byte) error {
	var res ChangeRemoteNodeOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*ChangeRemoteNodeOKBodyRemote RemoteNode represents generic remote Node. Agents can't run on Remote Nodes.
swagger:model ChangeRemoteNodeOKBodyRemote
*/
type ChangeRemoteNodeOKBodyRemote struct {

	// Unique randomly generated instance identifier, can't be changed.
	NodeID string `json:"node_id,omitempty"`

	// Unique across all Nodes user-defined name, can be changed.
	NodeName string `json:"node_name,omitempty"`
}

// Validate validates this change remote node o k body remote
func (o *ChangeRemoteNodeOKBodyRemote) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ChangeRemoteNodeOKBodyRemote) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ChangeRemoteNodeOKBodyRemote) UnmarshalBinary(b []byte) error {
	var res ChangeRemoteNodeOKBodyRemote
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
