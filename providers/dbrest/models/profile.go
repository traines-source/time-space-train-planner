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

// Profile A profile is a specific customisation for each endpoint.
// It parses data from the API differently, add additional information, or enable non-default methods.
//
// swagger:model Profile
type Profile struct {

	// endpoint
	Endpoint string `json:"endpoint,omitempty"`

	// journeys from trip
	JourneysFromTrip bool `json:"journeysFromTrip,omitempty"`

	// journeys walking speed
	JourneysWalkingSpeed bool `json:"journeysWalkingSpeed,omitempty"`

	// lines
	Lines bool `json:"lines,omitempty"`

	// locale
	Locale string `json:"locale,omitempty"`

	// products
	Products []*ProductType `json:"products"`

	// radar
	Radar bool `json:"radar,omitempty"`

	// reachable from
	ReachableFrom bool `json:"reachableFrom,omitempty"`

	// refresh journey
	RefreshJourney bool `json:"refreshJourney,omitempty"`

	// remarks
	Remarks bool `json:"remarks,omitempty"`

	// remarks get polyline
	RemarksGetPolyline bool `json:"remarksGetPolyline,omitempty"`

	// timezone
	Timezone string `json:"timezone,omitempty"`

	// trip
	Trip bool `json:"trip,omitempty"`

	// trips by name
	TripsByName bool `json:"tripsByName,omitempty"`
}

// Validate validates this profile
func (m *Profile) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateProducts(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Profile) validateProducts(formats strfmt.Registry) error {
	if swag.IsZero(m.Products) { // not required
		return nil
	}

	for i := 0; i < len(m.Products); i++ {
		if swag.IsZero(m.Products[i]) { // not required
			continue
		}

		if m.Products[i] != nil {
			if err := m.Products[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("products" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("products" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this profile based on the context it is used
func (m *Profile) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateProducts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Profile) contextValidateProducts(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Products); i++ {

		if m.Products[i] != nil {

			if swag.IsZero(m.Products[i]) { // not required
				return nil
			}

			if err := m.Products[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("products" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("products" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Profile) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Profile) UnmarshalBinary(b []byte) error {
	var res Profile
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}