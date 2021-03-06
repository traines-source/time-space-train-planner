// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ReferenceTripRelation A reference trip relation holds how a reference trip is related to a stop, for instance the reference trip starts after the stop. Stop contains a collection of that type, only if reference trips are available.
//
// swagger:model referenceTripRelation
type ReferenceTripRelation struct {

	// Reference trip element.
	// Required: true
	Rt *ReferenceTrip `json:"rt" xml:"rt"`

	// Relation to stop element.
	// Required: true
	Rts ReferenceTripRelationToStop `json:"rts" xml:"rts"`
}

// Validate validates this reference trip relation
func (m *ReferenceTripRelation) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRts(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ReferenceTripRelation) validateRt(formats strfmt.Registry) error {

	if err := validate.Required("rt", "body", m.Rt); err != nil {
		return err
	}

	if m.Rt != nil {
		if err := m.Rt.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("rt")
			}
			return err
		}
	}

	return nil
}

func (m *ReferenceTripRelation) validateRts(formats strfmt.Registry) error {

	if err := m.Rts.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("rts")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ReferenceTripRelation) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ReferenceTripRelation) UnmarshalBinary(b []byte) error {
	var res ReferenceTripRelation
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
