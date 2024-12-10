// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// TripsWithRealtimeData trips with realtime data
//
// swagger:model TripsWithRealtimeData
type TripsWithRealtimeData struct {

	// realtime data updated at
	RealtimeDataUpdatedAt float64 `json:"realtimeDataUpdatedAt,omitempty"`

	// trips
	Trips []*Trip `json:"trips"`
}

// Validate validates this trips with realtime data
func (m *TripsWithRealtimeData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTrips(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TripsWithRealtimeData) validateTrips(formats strfmt.Registry) error {
	if swag.IsZero(m.Trips) { // not required
		return nil
	}

	for i := 0; i < len(m.Trips); i++ {
		if swag.IsZero(m.Trips[i]) { // not required
			continue
		}

		if m.Trips[i] != nil {
			if err := m.Trips[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("trips" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("trips" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this trips with realtime data based on the context it is used
func (m *TripsWithRealtimeData) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateTrips(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TripsWithRealtimeData) contextValidateTrips(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Trips); i++ {

		if m.Trips[i] != nil {

			if swag.IsZero(m.Trips[i]) { // not required
				return nil
			}

			if err := m.Trips[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("trips" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("trips" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *TripsWithRealtimeData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TripsWithRealtimeData) UnmarshalBinary(b []byte) error {
	var res TripsWithRealtimeData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}