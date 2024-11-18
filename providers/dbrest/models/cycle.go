// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Cycle cycle
//
// swagger:model Cycle
type Cycle struct {

	// max
	Max float64 `json:"max,omitempty"`

	// min
	Min float64 `json:"min,omitempty"`

	// nr
	Nr float64 `json:"nr,omitempty"`
}

// Validate validates this cycle
func (m *Cycle) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this cycle based on context it is used
func (m *Cycle) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Cycle) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Cycle) UnmarshalBinary(b []byte) error {
	var res Cycle
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
