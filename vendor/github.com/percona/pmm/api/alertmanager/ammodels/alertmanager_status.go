// Code generated by go-swagger; DO NOT EDIT.

package ammodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// AlertmanagerStatus alertmanager status
// swagger:model alertmanagerStatus
type AlertmanagerStatus struct {

	// cluster
	// Required: true
	Cluster *ClusterStatus `json:"cluster"`

	// config
	// Required: true
	Config *AlertmanagerConfig `json:"config"`

	// uptime
	// Required: true
	// Format: date-time
	Uptime *strfmt.DateTime `json:"uptime"`

	// version info
	// Required: true
	VersionInfo *VersionInfo `json:"versionInfo"`
}

// Validate validates this alertmanager status
func (m *AlertmanagerStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCluster(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateConfig(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUptime(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVersionInfo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AlertmanagerStatus) validateCluster(formats strfmt.Registry) error {

	if err := validate.Required("cluster", "body", m.Cluster); err != nil {
		return err
	}

	if m.Cluster != nil {
		if err := m.Cluster.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("cluster")
			}
			return err
		}
	}

	return nil
}

func (m *AlertmanagerStatus) validateConfig(formats strfmt.Registry) error {

	if err := validate.Required("config", "body", m.Config); err != nil {
		return err
	}

	if m.Config != nil {
		if err := m.Config.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("config")
			}
			return err
		}
	}

	return nil
}

func (m *AlertmanagerStatus) validateUptime(formats strfmt.Registry) error {

	if err := validate.Required("uptime", "body", m.Uptime); err != nil {
		return err
	}

	if err := validate.FormatOf("uptime", "body", "date-time", m.Uptime.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *AlertmanagerStatus) validateVersionInfo(formats strfmt.Registry) error {

	if err := validate.Required("versionInfo", "body", m.VersionInfo); err != nil {
		return err
	}

	if m.VersionInfo != nil {
		if err := m.VersionInfo.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("versionInfo")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AlertmanagerStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AlertmanagerStatus) UnmarshalBinary(b []byte) error {
	var res AlertmanagerStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
