// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// IcoCrd ico crd
//
// swagger:model IcoCrd
type IcoCrd struct {

	// type
	Type string `json:"type,omitempty"`

	// x
	X float64 `json:"x,omitempty"`

	// y
	Y float64 `json:"y,omitempty"`
}

// Validate validates this ico crd
func (m *IcoCrd) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this ico crd based on context it is used
func (m *IcoCrd) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *IcoCrd) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IcoCrd) UnmarshalBinary(b []byte) error {
	var res IcoCrd
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
