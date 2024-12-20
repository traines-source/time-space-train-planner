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

// JourneysOptions journeys options
//
// swagger:model JourneysOptions
type JourneysOptions struct {
	JourneysOptionsCommon

	JourneysOptionsDbProfile
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *JourneysOptions) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 JourneysOptionsCommon
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.JourneysOptionsCommon = aO0

	// AO1
	var aO1 JourneysOptionsDbProfile
	if err := swag.ReadJSON(raw, &aO1); err != nil {
		return err
	}
	m.JourneysOptionsDbProfile = aO1

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m JourneysOptions) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.JourneysOptionsCommon)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	aO1, err := swag.WriteJSON(m.JourneysOptionsDbProfile)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this journeys options
func (m *JourneysOptions) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with JourneysOptionsCommon
	if err := m.JourneysOptionsCommon.Validate(formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with JourneysOptionsDbProfile
	if err := m.JourneysOptionsDbProfile.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validate this journeys options based on the context it is used
func (m *JourneysOptions) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with JourneysOptionsCommon
	if err := m.JourneysOptionsCommon.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with JourneysOptionsDbProfile
	if err := m.JourneysOptionsDbProfile.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *JourneysOptions) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *JourneysOptions) UnmarshalBinary(b []byte) error {
	var res JourneysOptions
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
