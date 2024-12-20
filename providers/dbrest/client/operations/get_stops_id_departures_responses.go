// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"traines.eu/time-space-train-planner/providers/dbrest/models"
)

// GetStopsIDDeparturesReader is a Reader for the GetStopsIDDepartures structure.
type GetStopsIDDeparturesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetStopsIDDeparturesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetStopsIDDeparturesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /stops/{id}/departures] GetStopsIDDepartures", response, response.Code())
	}
}

// NewGetStopsIDDeparturesOK creates a GetStopsIDDeparturesOK with default headers values
func NewGetStopsIDDeparturesOK() *GetStopsIDDeparturesOK {
	return &GetStopsIDDeparturesOK{}
}

/*
GetStopsIDDeparturesOK describes a response with status code 200, with default header values.

An object with an array of departures, in the [`hafas-client` format](https://github.com/public-transport/hafas-client/blob/6/docs/departures.md).
*/
type GetStopsIDDeparturesOK struct {
	Payload *GetStopsIDDeparturesOKBody
}

// IsSuccess returns true when this get stops Id departures o k response has a 2xx status code
func (o *GetStopsIDDeparturesOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get stops Id departures o k response has a 3xx status code
func (o *GetStopsIDDeparturesOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get stops Id departures o k response has a 4xx status code
func (o *GetStopsIDDeparturesOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get stops Id departures o k response has a 5xx status code
func (o *GetStopsIDDeparturesOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get stops Id departures o k response a status code equal to that given
func (o *GetStopsIDDeparturesOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get stops Id departures o k response
func (o *GetStopsIDDeparturesOK) Code() int {
	return 200
}

func (o *GetStopsIDDeparturesOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /stops/{id}/departures][%d] getStopsIdDeparturesOK %s", 200, payload)
}

func (o *GetStopsIDDeparturesOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /stops/{id}/departures][%d] getStopsIdDeparturesOK %s", 200, payload)
}

func (o *GetStopsIDDeparturesOK) GetPayload() *GetStopsIDDeparturesOKBody {
	return o.Payload
}

func (o *GetStopsIDDeparturesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetStopsIDDeparturesOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
GetStopsIDDeparturesOKBody get stops ID departures o k body
swagger:model GetStopsIDDeparturesOKBody
*/
type GetStopsIDDeparturesOKBody struct {

	// departures
	// Required: true
	Departures []*models.Alternative `json:"departures"`

	// realtime data updated at
	RealtimeDataUpdatedAt int64 `json:"realtimeDataUpdatedAt,omitempty"`
}

// Validate validates this get stops ID departures o k body
func (o *GetStopsIDDeparturesOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDepartures(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetStopsIDDeparturesOKBody) validateDepartures(formats strfmt.Registry) error {

	if err := validate.Required("getStopsIdDeparturesOK"+"."+"departures", "body", o.Departures); err != nil {
		return err
	}

	for i := 0; i < len(o.Departures); i++ {
		if swag.IsZero(o.Departures[i]) { // not required
			continue
		}

		if o.Departures[i] != nil {
			if err := o.Departures[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getStopsIdDeparturesOK" + "." + "departures" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getStopsIdDeparturesOK" + "." + "departures" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get stops ID departures o k body based on the context it is used
func (o *GetStopsIDDeparturesOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateDepartures(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetStopsIDDeparturesOKBody) contextValidateDepartures(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Departures); i++ {

		if o.Departures[i] != nil {

			if swag.IsZero(o.Departures[i]) { // not required
				return nil
			}

			if err := o.Departures[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getStopsIdDeparturesOK" + "." + "departures" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getStopsIdDeparturesOK" + "." + "departures" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetStopsIDDeparturesOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetStopsIDDeparturesOKBody) UnmarshalBinary(b []byte) error {
	var res GetStopsIDDeparturesOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
