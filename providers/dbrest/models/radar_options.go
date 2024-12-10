// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// RadarOptions radar options
//
// swagger:model RadarOptions
type RadarOptions struct {

	// compute frames for the next n seconds
	Duration *float64 `json:"duration,omitempty"`

	// parse & expose entrances of stops/stations?
	Entrances *bool `json:"entrances,omitempty"`

	// nr of frames to compute
	Frames *float64 `json:"frames,omitempty"`

	// return a shape for the trip?
	Polylines *bool `json:"polylines,omitempty"`

	// optionally an object of booleans
	Products Products `json:"products,omitempty"`

	// maximum number of vehicles
	Results *float64 `json:"results,omitempty"`

	// parse & expose sub-stops of stations?
	SubStops *bool `json:"subStops,omitempty"`

	// when
	// Format: date-time
	When strfmt.DateTime `json:"when,omitempty"`
}

// Validate validates this radar options
func (m *RadarOptions) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateProducts(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWhen(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RadarOptions) validateProducts(formats strfmt.Registry) error {
	if swag.IsZero(m.Products) { // not required
		return nil
	}

	if m.Products != nil {
		if err := m.Products.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("products")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("products")
			}
			return err
		}
	}

	return nil
}

func (m *RadarOptions) validateWhen(formats strfmt.Registry) error {
	if swag.IsZero(m.When) { // not required
		return nil
	}

	if err := validate.FormatOf("when", "body", "date-time", m.When.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this radar options based on the context it is used
func (m *RadarOptions) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateProducts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RadarOptions) contextValidateProducts(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Products) { // not required
		return nil
	}

	if err := m.Products.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("products")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("products")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *RadarOptions) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RadarOptions) UnmarshalBinary(b []byte) error {
	var res RadarOptions
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}