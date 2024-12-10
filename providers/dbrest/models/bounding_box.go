// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// BoundingBox bounding box
//
// swagger:model BoundingBox
type BoundingBox struct {

	// east
	East float64 `json:"east,omitempty"`

	// north
	North float64 `json:"north,omitempty"`

	// south
	South float64 `json:"south,omitempty"`

	// west
	West float64 `json:"west,omitempty"`
}

// Validate validates this bounding box
func (m *BoundingBox) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this bounding box based on context it is used
func (m *BoundingBox) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *BoundingBox) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BoundingBox) UnmarshalBinary(b []byte) error {
	var res BoundingBox
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}