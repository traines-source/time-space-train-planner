// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ArrivalDeparture arrival departure
//
// swagger:model ArrivalDeparture
type ArrivalDeparture struct {

	// arrival
	Arrival float64 `json:"arrival,omitempty"`

	// departure
	Departure float64 `json:"departure,omitempty"`
}

// Validate validates this arrival departure
func (m *ArrivalDeparture) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this arrival departure based on context it is used
func (m *ArrivalDeparture) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ArrivalDeparture) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ArrivalDeparture) UnmarshalBinary(b []byte) error {
	var res ArrivalDeparture
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
