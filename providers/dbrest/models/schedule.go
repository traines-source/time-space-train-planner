// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Schedule There are many ways to format schedules of public transport routes.
// This one tries to balance the amount of data and consumability.
// It is specifically geared towards urban public transport, with frequent trains and homogenous travels.
//
// swagger:model Schedule
type Schedule struct {

	// id
	ID string `json:"id,omitempty"`

	// mode
	// Enum: ["aircraft","bicycle","bus","car","gondola","taxi","train","walking","watercraft"]
	Mode string `json:"mode,omitempty"`

	// route
	Route string `json:"route,omitempty"`

	// sequence
	Sequence []*ArrivalDeparture `json:"sequence"`

	// array of Unix timestamps
	Starts []string `json:"starts"`

	// type
	// Enum: ["schedule"]
	Type string `json:"type,omitempty"`
}

// Validate validates this schedule
func (m *Schedule) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSequence(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var scheduleTypeModePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["aircraft","bicycle","bus","car","gondola","taxi","train","walking","watercraft"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		scheduleTypeModePropEnum = append(scheduleTypeModePropEnum, v)
	}
}

const (

	// ScheduleModeAircraft captures enum value "aircraft"
	ScheduleModeAircraft string = "aircraft"

	// ScheduleModeBicycle captures enum value "bicycle"
	ScheduleModeBicycle string = "bicycle"

	// ScheduleModeBus captures enum value "bus"
	ScheduleModeBus string = "bus"

	// ScheduleModeCar captures enum value "car"
	ScheduleModeCar string = "car"

	// ScheduleModeGondola captures enum value "gondola"
	ScheduleModeGondola string = "gondola"

	// ScheduleModeTaxi captures enum value "taxi"
	ScheduleModeTaxi string = "taxi"

	// ScheduleModeTrain captures enum value "train"
	ScheduleModeTrain string = "train"

	// ScheduleModeWalking captures enum value "walking"
	ScheduleModeWalking string = "walking"

	// ScheduleModeWatercraft captures enum value "watercraft"
	ScheduleModeWatercraft string = "watercraft"
)

// prop value enum
func (m *Schedule) validateModeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, scheduleTypeModePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Schedule) validateMode(formats strfmt.Registry) error {
	if swag.IsZero(m.Mode) { // not required
		return nil
	}

	// value enum
	if err := m.validateModeEnum("mode", "body", m.Mode); err != nil {
		return err
	}

	return nil
}

func (m *Schedule) validateSequence(formats strfmt.Registry) error {
	if swag.IsZero(m.Sequence) { // not required
		return nil
	}

	for i := 0; i < len(m.Sequence); i++ {
		if swag.IsZero(m.Sequence[i]) { // not required
			continue
		}

		if m.Sequence[i] != nil {
			if err := m.Sequence[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("sequence" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("sequence" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

var scheduleTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["schedule"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		scheduleTypeTypePropEnum = append(scheduleTypeTypePropEnum, v)
	}
}

const (

	// ScheduleTypeSchedule captures enum value "schedule"
	ScheduleTypeSchedule string = "schedule"
)

// prop value enum
func (m *Schedule) validateTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, scheduleTypeTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Schedule) validateType(formats strfmt.Registry) error {
	if swag.IsZero(m.Type) { // not required
		return nil
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this schedule based on the context it is used
func (m *Schedule) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSequence(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Schedule) contextValidateSequence(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Sequence); i++ {

		if m.Sequence[i] != nil {

			if swag.IsZero(m.Sequence[i]) { // not required
				return nil
			}

			if err := m.Sequence[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("sequence" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("sequence" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Schedule) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Schedule) UnmarshalBinary(b []byte) error {
	var res Schedule
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
