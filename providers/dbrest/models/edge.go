// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Edge edge
//
// swagger:model Edge
type Edge struct {

	// dir
	Dir float64 `json:"dir,omitempty"`

	// from location
	FromLocation *Stop `json:"fromLocation,omitempty"`

	// ico crd
	IcoCrd *IcoCrd `json:"icoCrd,omitempty"`

	// icon
	Icon interface{} `json:"icon,omitempty"`

	// to location
	ToLocation *Stop `json:"toLocation,omitempty"`
}

// Validate validates this edge
func (m *Edge) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFromLocation(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIcoCrd(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateToLocation(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Edge) validateFromLocation(formats strfmt.Registry) error {
	if swag.IsZero(m.FromLocation) { // not required
		return nil
	}

	if m.FromLocation != nil {
		if err := m.FromLocation.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("fromLocation")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("fromLocation")
			}
			return err
		}
	}

	return nil
}

func (m *Edge) validateIcoCrd(formats strfmt.Registry) error {
	if swag.IsZero(m.IcoCrd) { // not required
		return nil
	}

	if m.IcoCrd != nil {
		if err := m.IcoCrd.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("icoCrd")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("icoCrd")
			}
			return err
		}
	}

	return nil
}

func (m *Edge) validateToLocation(formats strfmt.Registry) error {
	if swag.IsZero(m.ToLocation) { // not required
		return nil
	}

	if m.ToLocation != nil {
		if err := m.ToLocation.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("toLocation")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("toLocation")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this edge based on the context it is used
func (m *Edge) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateFromLocation(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateIcoCrd(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateToLocation(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Edge) contextValidateFromLocation(ctx context.Context, formats strfmt.Registry) error {

	if m.FromLocation != nil {

		if swag.IsZero(m.FromLocation) { // not required
			return nil
		}

		if err := m.FromLocation.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("fromLocation")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("fromLocation")
			}
			return err
		}
	}

	return nil
}

func (m *Edge) contextValidateIcoCrd(ctx context.Context, formats strfmt.Registry) error {

	if m.IcoCrd != nil {

		if swag.IsZero(m.IcoCrd) { // not required
			return nil
		}

		if err := m.IcoCrd.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("icoCrd")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("icoCrd")
			}
			return err
		}
	}

	return nil
}

func (m *Edge) contextValidateToLocation(ctx context.Context, formats strfmt.Registry) error {

	if m.ToLocation != nil {

		if swag.IsZero(m.ToLocation) { // not required
			return nil
		}

		if err := m.ToLocation.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("toLocation")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("toLocation")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Edge) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Edge) UnmarshalBinary(b []byte) error {
	var res Edge
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
