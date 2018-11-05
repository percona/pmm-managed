// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// InventoryRDSNode RDSNode represents AWS RDS Node.
// swagger:model inventoryRDSNode
type InventoryRDSNode struct {

	// Unique Node identifier.
	ID int64 `json:"id,omitempty"`

	// Unique Node name.
	Name string `json:"name,omitempty"`

	// region
	Region string `json:"region,omitempty"`
}

// Validate validates this inventory RDS node
func (m *InventoryRDSNode) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *InventoryRDSNode) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *InventoryRDSNode) UnmarshalBinary(b []byte) error {
	var res InventoryRDSNode
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
