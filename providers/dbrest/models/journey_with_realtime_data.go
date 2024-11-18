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

// JourneyWithRealtimeData journey with realtime data
//
// swagger:model JourneyWithRealtimeData
type JourneyWithRealtimeData struct {

	// journey
	Journey *Journey `json:"journey,omitempty"`

	// realtime data updated at
	RealtimeDataUpdatedAt float64 `json:"realtimeDataUpdatedAt,omitempty"`
}

// Validate validates this journey with realtime data
func (m *JourneyWithRealtimeData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateJourney(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *JourneyWithRealtimeData) validateJourney(formats strfmt.Registry) error {
	if swag.IsZero(m.Journey) { // not required
		return nil
	}

	if m.Journey != nil {
		if err := m.Journey.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("journey")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("journey")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this journey with realtime data based on the context it is used
func (m *JourneyWithRealtimeData) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateJourney(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *JourneyWithRealtimeData) contextValidateJourney(ctx context.Context, formats strfmt.Registry) error {

	if m.Journey != nil {

		if swag.IsZero(m.Journey) { // not required
			return nil
		}

		if err := m.Journey.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("journey")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("journey")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *JourneyWithRealtimeData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *JourneyWithRealtimeData) UnmarshalBinary(b []byte) error {
	var res JourneyWithRealtimeData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
