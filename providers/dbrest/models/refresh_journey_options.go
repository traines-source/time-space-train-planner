// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// RefreshJourneyOptions refresh journey options
//
// swagger:model RefreshJourneyOptions
type RefreshJourneyOptions struct {

	// parse & expose entrances of stops/stations?
	Entrances *bool `json:"entrances,omitempty"`

	// language
	Language *string `json:"language,omitempty"`

	// return a shape for each leg?
	Polylines *bool `json:"polylines,omitempty"`

	// parse & expose hints & warnings?
	Remarks *bool `json:"remarks,omitempty"`

	// parse & expose dates the journey is valid on?
	ScheduledDays *bool `json:"scheduledDays,omitempty"`

	// return stations on the way?
	Stopovers *bool `json:"stopovers,omitempty"`

	// parse & expose sub-stops of stations?
	SubStops *bool `json:"subStops,omitempty"`

	// return tickets? only available with some profiles
	Tickets *bool `json:"tickets,omitempty"`
}

// Validate validates this refresh journey options
func (m *RefreshJourneyOptions) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this refresh journey options based on context it is used
func (m *RefreshJourneyOptions) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RefreshJourneyOptions) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RefreshJourneyOptions) UnmarshalBinary(b []byte) error {
	var res RefreshJourneyOptions
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}