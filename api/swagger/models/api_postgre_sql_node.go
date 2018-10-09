// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// APIPostgreSQLNode api postgre SQL node
// swagger:model apiPostgreSQLNode
type APIPostgreSQLNode struct {

	// name
	Name string `json:"name,omitempty"`

	// region
	Region string `json:"region,omitempty"`
}

// Validate validates this api postgre SQL node
func (m *APIPostgreSQLNode) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *APIPostgreSQLNode) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *APIPostgreSQLNode) UnmarshalBinary(b []byte) error {
	var res APIPostgreSQLNode
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
