// Code generated by go-swagger; DO NOT EDIT.

package nodes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// ListNodesReader is a Reader for the ListNodes structure.
type ListNodesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListNodesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewListNodesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewListNodesOK creates a ListNodesOK with default headers values
func NewListNodesOK() *ListNodesOK {
	return &ListNodesOK{}
}

/*ListNodesOK handles this case with default header values.

(empty)
*/
type ListNodesOK struct {
	Payload *ListNodesOKBody
}

func (o *ListNodesOK) Error() string {
	return fmt.Sprintf("[POST /v0/inventory/Nodes/List][%d] listNodesOK  %+v", 200, o.Payload)
}

func (o *ListNodesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(ListNodesOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*AmazonRDSRemoteItems0 AmazonRDSRemoteNode represents Amazon (AWS) RDS remote Node.
swagger:model AmazonRDSRemoteItems0
*/
type AmazonRDSRemoteItems0 struct {

	// Hostname. Unique in combination with region.
	Hostname string `json:"hostname,omitempty"`

	// Unique Node identifier.
	ID string `json:"id,omitempty"`

	// Unique user-defined Node name.
	Name string `json:"name,omitempty"`

	// AWS region. Unique in combination with hostname.
	Region string `json:"region,omitempty"`
}

// Validate validates this amazon RDS remote items0
func (o *AmazonRDSRemoteItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *AmazonRDSRemoteItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AmazonRDSRemoteItems0) UnmarshalBinary(b []byte) error {
	var res AmazonRDSRemoteItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*GenericItems0 GenericNode represents Node without more specialized type.
swagger:model GenericItems0
*/
type GenericItems0 struct {

	// Hostname. Is not unique. May be empty.
	Hostname string `json:"hostname,omitempty"`

	// Unique Node identifier.
	ID string `json:"id,omitempty"`

	// Unique user-defined Node name.
	Name string `json:"name,omitempty"`
}

// Validate validates this generic items0
func (o *GenericItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GenericItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GenericItems0) UnmarshalBinary(b []byte) error {
	var res GenericItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*ListNodesOKBody list nodes o k body
swagger:model ListNodesOKBody
*/
type ListNodesOKBody struct {

	// amazon rds remote
	AmazonRDSRemote []*AmazonRDSRemoteItems0 `json:"amazon_rds_remote"`

	// generic
	Generic []*GenericItems0 `json:"generic"`

	// remote
	Remote []*RemoteItems0 `json:"remote"`
}

// Validate validates this list nodes o k body
func (o *ListNodesOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAmazonRDSRemote(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateGeneric(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateRemote(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ListNodesOKBody) validateAmazonRDSRemote(formats strfmt.Registry) error {

	if swag.IsZero(o.AmazonRDSRemote) { // not required
		return nil
	}

	for i := 0; i < len(o.AmazonRDSRemote); i++ {
		if swag.IsZero(o.AmazonRDSRemote[i]) { // not required
			continue
		}

		if o.AmazonRDSRemote[i] != nil {
			if err := o.AmazonRDSRemote[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listNodesOK" + "." + "amazon_rds_remote" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *ListNodesOKBody) validateGeneric(formats strfmt.Registry) error {

	if swag.IsZero(o.Generic) { // not required
		return nil
	}

	for i := 0; i < len(o.Generic); i++ {
		if swag.IsZero(o.Generic[i]) { // not required
			continue
		}

		if o.Generic[i] != nil {
			if err := o.Generic[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listNodesOK" + "." + "generic" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (o *ListNodesOKBody) validateRemote(formats strfmt.Registry) error {

	if swag.IsZero(o.Remote) { // not required
		return nil
	}

	for i := 0; i < len(o.Remote); i++ {
		if swag.IsZero(o.Remote[i]) { // not required
			continue
		}

		if o.Remote[i] != nil {
			if err := o.Remote[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("listNodesOK" + "." + "remote" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *ListNodesOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ListNodesOKBody) UnmarshalBinary(b []byte) error {
	var res ListNodesOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*RemoteItems0 RemoteNode represents a generic remote Node.
// Agents can't be run on remote Nodes.
swagger:model RemoteItems0
*/
type RemoteItems0 struct {

	// Unique Node identifier.
	ID string `json:"id,omitempty"`

	// Unique user-defined Node name.
	Name string `json:"name,omitempty"`
}

// Validate validates this remote items0
func (o *RemoteItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *RemoteItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RemoteItems0) UnmarshalBinary(b []byte) error {
	var res RemoteItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
