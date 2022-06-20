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
	"github.com/go-openapi/validate"
)

// DepartureArrival Dep or Arr
//
// swagger:model departureArrival
type DepartureArrival struct {

	// delay
	// Required: true
	Delay *int64 `json:"delay"`

	// direction
	// Required: true
	Direction *string `json:"direction"`

	// line
	// Required: true
	Line *DepartureArrivalLine `json:"line"`

	// load factor
	// Required: true
	LoadFactor *string `json:"loadFactor"`

	// planned platform
	// Required: true
	PlannedPlatform *string `json:"plannedPlatform"`

	// planned when
	// Required: true
	// Format: date-time
	PlannedWhen *strfmt.DateTime `json:"plannedWhen"`

	// platform
	// Required: true
	Platform *string `json:"platform"`

	// provenance
	// Required: true
	Provenance interface{} `json:"provenance"`

	// remarks
	// Required: true
	Remarks []*DepartureArrivalRemarksItems0 `json:"remarks"`

	// stop
	// Required: true
	Stop *DepartureArrivalStop `json:"stop"`

	// trip Id
	// Required: true
	TripID *string `json:"tripId"`

	// when
	// Required: true
	// Format: date-time
	When *strfmt.DateTime `json:"when"`
}

// Validate validates this departure arrival
func (m *DepartureArrival) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDelay(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDirection(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLine(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLoadFactor(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePlannedPlatform(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePlannedWhen(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePlatform(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProvenance(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRemarks(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStop(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTripID(formats); err != nil {
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

func (m *DepartureArrival) validateDelay(formats strfmt.Registry) error {

	if err := validate.Required("delay", "body", m.Delay); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrival) validateDirection(formats strfmt.Registry) error {

	if err := validate.Required("direction", "body", m.Direction); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrival) validateLine(formats strfmt.Registry) error {

	if err := validate.Required("line", "body", m.Line); err != nil {
		return err
	}

	if m.Line != nil {
		if err := m.Line.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("line")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("line")
			}
			return err
		}
	}

	return nil
}

func (m *DepartureArrival) validateLoadFactor(formats strfmt.Registry) error {

	if err := validate.Required("loadFactor", "body", m.LoadFactor); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrival) validatePlannedPlatform(formats strfmt.Registry) error {

	if err := validate.Required("plannedPlatform", "body", m.PlannedPlatform); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrival) validatePlannedWhen(formats strfmt.Registry) error {

	if err := validate.Required("plannedWhen", "body", m.PlannedWhen); err != nil {
		return err
	}

	if err := validate.FormatOf("plannedWhen", "body", "date-time", m.PlannedWhen.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrival) validatePlatform(formats strfmt.Registry) error {

	if err := validate.Required("platform", "body", m.Platform); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrival) validateProvenance(formats strfmt.Registry) error {

	return nil
}

func (m *DepartureArrival) validateRemarks(formats strfmt.Registry) error {

	if err := validate.Required("remarks", "body", m.Remarks); err != nil {
		return err
	}

	for i := 0; i < len(m.Remarks); i++ {
		if swag.IsZero(m.Remarks[i]) { // not required
			continue
		}

		if m.Remarks[i] != nil {
			if err := m.Remarks[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("remarks" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("remarks" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *DepartureArrival) validateStop(formats strfmt.Registry) error {

	if err := validate.Required("stop", "body", m.Stop); err != nil {
		return err
	}

	if m.Stop != nil {
		if err := m.Stop.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("stop")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("stop")
			}
			return err
		}
	}

	return nil
}

func (m *DepartureArrival) validateTripID(formats strfmt.Registry) error {

	if err := validate.Required("tripId", "body", m.TripID); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrival) validateWhen(formats strfmt.Registry) error {

	if err := validate.Required("when", "body", m.When); err != nil {
		return err
	}

	if err := validate.FormatOf("when", "body", "date-time", m.When.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this departure arrival based on the context it is used
func (m *DepartureArrival) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLine(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRemarks(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStop(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DepartureArrival) contextValidateLine(ctx context.Context, formats strfmt.Registry) error {

	if m.Line != nil {
		if err := m.Line.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("line")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("line")
			}
			return err
		}
	}

	return nil
}

func (m *DepartureArrival) contextValidateRemarks(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Remarks); i++ {

		if m.Remarks[i] != nil {
			if err := m.Remarks[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("remarks" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("remarks" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *DepartureArrival) contextValidateStop(ctx context.Context, formats strfmt.Registry) error {

	if m.Stop != nil {
		if err := m.Stop.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("stop")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("stop")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DepartureArrival) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DepartureArrival) UnmarshalBinary(b []byte) error {
	var res DepartureArrival
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DepartureArrivalLine departure arrival line
//
// swagger:model DepartureArrivalLine
type DepartureArrivalLine struct {

	// admin code
	// Required: true
	AdminCode *string `json:"adminCode"`

	// fahrt nr
	// Required: true
	FahrtNr *string `json:"fahrtNr"`

	// id
	// Required: true
	ID *string `json:"id"`

	// mode
	// Required: true
	Mode *string `json:"mode"`

	// name
	// Required: true
	Name *string `json:"name"`

	// operator
	// Required: true
	Operator *DepartureArrivalLineOperator `json:"operator"`

	// product
	// Required: true
	Product *string `json:"product"`

	// product name
	// Required: true
	ProductName *string `json:"productName"`

	// public
	// Required: true
	Public *bool `json:"public"`

	// type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this departure arrival line
func (m *DepartureArrivalLine) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAdminCode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFahrtNr(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOperator(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProduct(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProductName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePublic(formats); err != nil {
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

func (m *DepartureArrivalLine) validateAdminCode(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"adminCode", "body", m.AdminCode); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalLine) validateFahrtNr(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"fahrtNr", "body", m.FahrtNr); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalLine) validateID(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalLine) validateMode(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"mode", "body", m.Mode); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalLine) validateName(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalLine) validateOperator(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"operator", "body", m.Operator); err != nil {
		return err
	}

	if m.Operator != nil {
		if err := m.Operator.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("line" + "." + "operator")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("line" + "." + "operator")
			}
			return err
		}
	}

	return nil
}

func (m *DepartureArrivalLine) validateProduct(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"product", "body", m.Product); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalLine) validateProductName(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"productName", "body", m.ProductName); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalLine) validatePublic(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"public", "body", m.Public); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalLine) validateType(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this departure arrival line based on the context it is used
func (m *DepartureArrivalLine) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateOperator(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DepartureArrivalLine) contextValidateOperator(ctx context.Context, formats strfmt.Registry) error {

	if m.Operator != nil {
		if err := m.Operator.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("line" + "." + "operator")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("line" + "." + "operator")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DepartureArrivalLine) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DepartureArrivalLine) UnmarshalBinary(b []byte) error {
	var res DepartureArrivalLine
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DepartureArrivalLineOperator departure arrival line operator
//
// swagger:model DepartureArrivalLineOperator
type DepartureArrivalLineOperator struct {

	// id
	// Required: true
	ID *string `json:"id"`

	// name
	// Required: true
	Name *string `json:"name"`

	// type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this departure arrival line operator
func (m *DepartureArrivalLineOperator) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
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

func (m *DepartureArrivalLineOperator) validateID(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"operator"+"."+"id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalLineOperator) validateName(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"operator"+"."+"name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalLineOperator) validateType(formats strfmt.Registry) error {

	if err := validate.Required("line"+"."+"operator"+"."+"type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this departure arrival line operator based on context it is used
func (m *DepartureArrivalLineOperator) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DepartureArrivalLineOperator) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DepartureArrivalLineOperator) UnmarshalBinary(b []byte) error {
	var res DepartureArrivalLineOperator
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DepartureArrivalRemarksItems0 departure arrival remarks items0
//
// swagger:model DepartureArrivalRemarksItems0
type DepartureArrivalRemarksItems0 struct {

	// text
	Text string `json:"text,omitempty"`
}

// Validate validates this departure arrival remarks items0
func (m *DepartureArrivalRemarksItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this departure arrival remarks items0 based on context it is used
func (m *DepartureArrivalRemarksItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DepartureArrivalRemarksItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DepartureArrivalRemarksItems0) UnmarshalBinary(b []byte) error {
	var res DepartureArrivalRemarksItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DepartureArrivalStop departure arrival stop
//
// swagger:model DepartureArrivalStop
type DepartureArrivalStop struct {

	// id
	// Required: true
	ID *string `json:"id"`

	// location
	// Required: true
	Location *DepartureArrivalStopLocation `json:"location"`

	// name
	// Required: true
	Name *string `json:"name"`

	// products
	// Required: true
	Products *DepartureArrivalStopProducts `json:"products"`

	// type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this departure arrival stop
func (m *DepartureArrivalStop) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLocation(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProducts(formats); err != nil {
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

func (m *DepartureArrivalStop) validateID(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStop) validateLocation(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"location", "body", m.Location); err != nil {
		return err
	}

	if m.Location != nil {
		if err := m.Location.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("stop" + "." + "location")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("stop" + "." + "location")
			}
			return err
		}
	}

	return nil
}

func (m *DepartureArrivalStop) validateName(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStop) validateProducts(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"products", "body", m.Products); err != nil {
		return err
	}

	if m.Products != nil {
		if err := m.Products.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("stop" + "." + "products")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("stop" + "." + "products")
			}
			return err
		}
	}

	return nil
}

func (m *DepartureArrivalStop) validateType(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this departure arrival stop based on the context it is used
func (m *DepartureArrivalStop) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLocation(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateProducts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DepartureArrivalStop) contextValidateLocation(ctx context.Context, formats strfmt.Registry) error {

	if m.Location != nil {
		if err := m.Location.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("stop" + "." + "location")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("stop" + "." + "location")
			}
			return err
		}
	}

	return nil
}

func (m *DepartureArrivalStop) contextValidateProducts(ctx context.Context, formats strfmt.Registry) error {

	if m.Products != nil {
		if err := m.Products.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("stop" + "." + "products")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("stop" + "." + "products")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DepartureArrivalStop) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DepartureArrivalStop) UnmarshalBinary(b []byte) error {
	var res DepartureArrivalStop
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DepartureArrivalStopLocation departure arrival stop location
//
// swagger:model DepartureArrivalStopLocation
type DepartureArrivalStopLocation struct {

	// id
	// Required: true
	ID *string `json:"id"`

	// latitude
	// Required: true
	Latitude *float64 `json:"latitude"`

	// longitude
	// Required: true
	Longitude *float64 `json:"longitude"`

	// type
	// Required: true
	Type *string `json:"type"`
}

// Validate validates this departure arrival stop location
func (m *DepartureArrivalStopLocation) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLatitude(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLongitude(formats); err != nil {
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

func (m *DepartureArrivalStopLocation) validateID(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"location"+"."+"id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopLocation) validateLatitude(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"location"+"."+"latitude", "body", m.Latitude); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopLocation) validateLongitude(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"location"+"."+"longitude", "body", m.Longitude); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopLocation) validateType(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"location"+"."+"type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this departure arrival stop location based on context it is used
func (m *DepartureArrivalStopLocation) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DepartureArrivalStopLocation) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DepartureArrivalStopLocation) UnmarshalBinary(b []byte) error {
	var res DepartureArrivalStopLocation
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// DepartureArrivalStopProducts departure arrival stop products
//
// swagger:model DepartureArrivalStopProducts
type DepartureArrivalStopProducts struct {

	// bus
	// Required: true
	Bus *bool `json:"bus"`

	// ferry
	// Required: true
	Ferry *bool `json:"ferry"`

	// national
	// Required: true
	National *bool `json:"national"`

	// national express
	// Required: true
	NationalExpress *bool `json:"nationalExpress"`

	// regional
	// Required: true
	Regional *bool `json:"regional"`

	// regional exp
	// Required: true
	RegionalExp *bool `json:"regionalExp"`

	// suburban
	// Required: true
	Suburban *bool `json:"suburban"`

	// subway
	// Required: true
	Subway *bool `json:"subway"`

	// taxi
	// Required: true
	Taxi *bool `json:"taxi"`

	// tram
	// Required: true
	Tram *bool `json:"tram"`
}

// Validate validates this departure arrival stop products
func (m *DepartureArrivalStopProducts) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFerry(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNational(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNationalExpress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRegional(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRegionalExp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSuburban(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSubway(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTaxi(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTram(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DepartureArrivalStopProducts) validateBus(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"products"+"."+"bus", "body", m.Bus); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopProducts) validateFerry(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"products"+"."+"ferry", "body", m.Ferry); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopProducts) validateNational(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"products"+"."+"national", "body", m.National); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopProducts) validateNationalExpress(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"products"+"."+"nationalExpress", "body", m.NationalExpress); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopProducts) validateRegional(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"products"+"."+"regional", "body", m.Regional); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopProducts) validateRegionalExp(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"products"+"."+"regionalExp", "body", m.RegionalExp); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopProducts) validateSuburban(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"products"+"."+"suburban", "body", m.Suburban); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopProducts) validateSubway(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"products"+"."+"subway", "body", m.Subway); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopProducts) validateTaxi(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"products"+"."+"taxi", "body", m.Taxi); err != nil {
		return err
	}

	return nil
}

func (m *DepartureArrivalStopProducts) validateTram(formats strfmt.Registry) error {

	if err := validate.Required("stop"+"."+"products"+"."+"tram", "body", m.Tram); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this departure arrival stop products based on context it is used
func (m *DepartureArrivalStopProducts) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DepartureArrivalStopProducts) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DepartureArrivalStopProducts) UnmarshalBinary(b []byte) error {
	var res DepartureArrivalStopProducts
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
