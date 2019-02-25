// Code generated by go-swagger; DO NOT EDIT.

package services

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

// ChangeMySQLServiceReader is a Reader for the ChangeMySQLService structure.
type ChangeMySQLServiceReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ChangeMySQLServiceReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewChangeMySQLServiceOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewChangeMySQLServiceOK creates a ChangeMySQLServiceOK with default headers values
func NewChangeMySQLServiceOK() *ChangeMySQLServiceOK {
	return &ChangeMySQLServiceOK{}
}

/*ChangeMySQLServiceOK handles this case with default header values.

A successful response.
*/
type ChangeMySQLServiceOK struct {
	Payload *ChangeMySQLServiceOKBody
}

func (o *ChangeMySQLServiceOK) Error() string {
	return fmt.Sprintf("[POST /v1/inventory/Services/ChangeMySQL][%d] changeMySqlServiceOK  %+v", 200, o.Payload)
}

func (o *ChangeMySQLServiceOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(ChangeMySQLServiceOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*ChangeMySQLServiceBody change my SQL service body
swagger:model ChangeMySQLServiceBody
*/
type ChangeMySQLServiceBody struct {

	// Access address (DNS name or IP). Required if unix_socket is absent.
	Address string `json:"address,omitempty"`

	// Custom user-assigned labels. Keys must start with "_".
	CustomLabels map[string]string `json:"custom_labels,omitempty"`

	// Node identifier where this instance runs.
	NodeID string `json:"node_id,omitempty"`

	// Access port. Required if unix_socket is absent.
	Port int64 `json:"port,omitempty"`

	// Unique randomly generated instance identifier.
	ServiceID string `json:"service_id,omitempty"`

	// Unique across all Services user-defined name.
	ServiceName string `json:"service_name,omitempty"`

	// Access Unix socket. Required if address and port are absent.
	UnixSocket string `json:"unix_socket,omitempty"`
}

// Validate validates this change my SQL service body
func (o *ChangeMySQLServiceBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ChangeMySQLServiceBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ChangeMySQLServiceBody) UnmarshalBinary(b []byte) error {
	var res ChangeMySQLServiceBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*ChangeMySQLServiceOKBody change my SQL service o k body
swagger:model ChangeMySQLServiceOKBody
*/
type ChangeMySQLServiceOKBody struct {

	// mysql
	Mysql *ChangeMySQLServiceOKBodyMysql `json:"mysql,omitempty"`
}

// Validate validates this change my SQL service o k body
func (o *ChangeMySQLServiceOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMysql(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ChangeMySQLServiceOKBody) validateMysql(formats strfmt.Registry) error {

	if swag.IsZero(o.Mysql) { // not required
		return nil
	}

	if o.Mysql != nil {
		if err := o.Mysql.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("changeMySqlServiceOK" + "." + "mysql")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *ChangeMySQLServiceOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ChangeMySQLServiceOKBody) UnmarshalBinary(b []byte) error {
	var res ChangeMySQLServiceOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

/*ChangeMySQLServiceOKBodyMysql MySQLService represents a generic MySQL instance.
swagger:model ChangeMySQLServiceOKBodyMysql
*/
type ChangeMySQLServiceOKBodyMysql struct {

	// Access address (DNS name or IP). Required if unix_socket is absent.
	Address string `json:"address,omitempty"`

	// Custom user-assigned labels. Keys must start with "_".
	CustomLabels map[string]string `json:"custom_labels,omitempty"`

	// Node identifier where this instance runs.
	NodeID string `json:"node_id,omitempty"`

	// Access port. Required if unix_socket is absent.
	Port int64 `json:"port,omitempty"`

	// Unique randomly generated instance identifier.
	ServiceID string `json:"service_id,omitempty"`

	// Unique across all Services user-defined name.
	ServiceName string `json:"service_name,omitempty"`

	// Access Unix socket. Required if address and port are absent.
	UnixSocket string `json:"unix_socket,omitempty"`
}

// Validate validates this change my SQL service o k body mysql
func (o *ChangeMySQLServiceOKBodyMysql) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ChangeMySQLServiceOKBodyMysql) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ChangeMySQLServiceOKBodyMysql) UnmarshalBinary(b []byte) error {
	var res ChangeMySQLServiceOKBodyMysql
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
