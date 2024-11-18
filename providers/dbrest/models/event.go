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

// Event event
//
// swagger:model Event
type Event struct {

	// end
	End string `json:"end,omitempty"`

	// from location
	FromLocation *Stop `json:"fromLocation,omitempty"`

	// sections
	Sections []string `json:"sections"`

	// start
	Start string `json:"start,omitempty"`

	// to location
	ToLocation *Stop `json:"toLocation,omitempty"`
}

// Validate validates this event
func (m *Event) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFromLocation(formats); err != nil {
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

func (m *Event) validateFromLocation(formats strfmt.Registry) error {
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

func (m *Event) validateToLocation(formats strfmt.Registry) error {
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

// ContextValidate validate this event based on the context it is used
func (m *Event) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateFromLocation(ctx, formats); err != nil {
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

func (m *Event) contextValidateFromLocation(ctx context.Context, formats strfmt.Registry) error {

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

func (m *Event) contextValidateToLocation(ctx context.Context, formats strfmt.Registry) error {

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
func (m *Event) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Event) UnmarshalBinary(b []byte) error {
	var res Event
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
