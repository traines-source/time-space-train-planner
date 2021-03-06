// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HistoricDelay It's the history of all delay-messages for a stop. This element extends HistoricChange.
//
// swagger:model historicDelay
type HistoricDelay struct {

	// The arrival event. The time, in ten digit 'YYMMddHHmm' format, e.g. '1404011437' for 14:37 on April the 1st of 2014.
	Ar string `json:"ar,omitempty" xml:"ar,attr,omitempty"`

	// Detailed description of delay cause.
	Cod string `json:"cod,omitempty" xml:"cod,attr,omitempty"`

	// The departure event. The time, in ten digit 'YYMMddHHmm' format, e.g. '1404011437' for 14:37 on April the 1st of 2014.
	Dp string `json:"dp,omitempty" xml:"dp,attr,omitempty"`

	// Source of the message.
	Src DelaySource `json:"src,omitempty" xml:"src,attr,omitempty"`

	// Timestamp. The time, in ten digit 'YYMMddHHmm' format, e.g. '1404011437' for 14:37 on April the 1st of 2014.
	Ts string `json:"ts,omitempty" xml:"ts,attr,omitempty"`
}

// Validate validates this historic delay
func (m *HistoricDelay) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSrc(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HistoricDelay) validateSrc(formats strfmt.Registry) error {

	if swag.IsZero(m.Src) { // not required
		return nil
	}

	if err := m.Src.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("src")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HistoricDelay) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HistoricDelay) UnmarshalBinary(b []byte) error {
	var res HistoricDelay
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
