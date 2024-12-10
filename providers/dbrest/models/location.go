// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Location A location object is used by other items to indicate their locations.
//
// swagger:model Location
type Location struct {

	// address
	Address string `json:"address,omitempty"`

	// altitude
	Altitude float64 `json:"altitude,omitempty"`

	// distance
	Distance float64 `json:"distance,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// latitude
	Latitude float64 `json:"latitude,omitempty"`

	// longitude
	Longitude float64 `json:"longitude,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// poi
	Poi bool `json:"poi,omitempty"`

	// type
	// Enum: ["location"]
	Type string `json:"type,omitempty"`
}

// Validate validates this location
func (m *Location) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var locationTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["location"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		locationTypeTypePropEnum = append(locationTypeTypePropEnum, v)
	}
}

const (

	// LocationTypeLocation captures enum value "location"
	LocationTypeLocation string = "location"
)

// prop value enum
func (m *Location) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, locationTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Location) validateType(formats strfmt.Registry) error {
	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this location based on context it is used
func (m *Location) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Location) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Location) UnmarshalBinary(b []byte) error {
	var res Location
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}