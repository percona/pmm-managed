// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// InventoryGetNodeResponse inventory get node response
// swagger:model inventoryGetNodeResponse
type InventoryGetNodeResponse struct {

	// bare metal
	BareMetal *InventoryBareMetalNode `json:"bare_metal,omitempty"`

	// container
	Container *InventoryContainerNode `json:"container,omitempty"`

	// rds
	RDS *InventoryRDSNode `json:"rds,omitempty"`

	// remote
	Remote *InventoryRemoteNode `json:"remote,omitempty"`

	// virtual machine
	VirtualMachine *InventoryVirtualMachineNode `json:"virtual_machine,omitempty"`
}

// Validate validates this inventory get node response
func (m *InventoryGetNodeResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBareMetal(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateContainer(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRDS(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRemote(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVirtualMachine(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *InventoryGetNodeResponse) validateBareMetal(formats strfmt.Registry) error {

	if swag.IsZero(m.BareMetal) { // not required
		return nil
	}

	if m.BareMetal != nil {
		if err := m.BareMetal.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("bare_metal")
			}
			return err
		}
	}

	return nil
}

func (m *InventoryGetNodeResponse) validateContainer(formats strfmt.Registry) error {

	if swag.IsZero(m.Container) { // not required
		return nil
	}

	if m.Container != nil {
		if err := m.Container.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("container")
			}
			return err
		}
	}

	return nil
}

func (m *InventoryGetNodeResponse) validateRDS(formats strfmt.Registry) error {

	if swag.IsZero(m.RDS) { // not required
		return nil
	}

	if m.RDS != nil {
		if err := m.RDS.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("rds")
			}
			return err
		}
	}

	return nil
}

func (m *InventoryGetNodeResponse) validateRemote(formats strfmt.Registry) error {

	if swag.IsZero(m.Remote) { // not required
		return nil
	}

	if m.Remote != nil {
		if err := m.Remote.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("remote")
			}
			return err
		}
	}

	return nil
}

func (m *InventoryGetNodeResponse) validateVirtualMachine(formats strfmt.Registry) error {

	if swag.IsZero(m.VirtualMachine) { // not required
		return nil
	}

	if m.VirtualMachine != nil {
		if err := m.VirtualMachine.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("virtual_machine")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *InventoryGetNodeResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *InventoryGetNodeResponse) UnmarshalBinary(b []byte) error {
	var res InventoryGetNodeResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
