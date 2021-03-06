// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"traines.eu/time-space-train-planner/providers/dbtimetables/models"
)

// GetTimetableFchgEvaNoReader is a Reader for the GetTimetableFchgEvaNo structure.
type GetTimetableFchgEvaNoReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetTimetableFchgEvaNoReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetTimetableFchgEvaNoOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetTimetableFchgEvaNoNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetTimetableFchgEvaNoOK creates a GetTimetableFchgEvaNoOK with default headers values
func NewGetTimetableFchgEvaNoOK() *GetTimetableFchgEvaNoOK {
	return &GetTimetableFchgEvaNoOK{}
}

/*GetTimetableFchgEvaNoOK handles this case with default header values.

successful operation
*/
type GetTimetableFchgEvaNoOK struct {
	Payload *models.Timetable
}

func (o *GetTimetableFchgEvaNoOK) Error() string {
	return fmt.Sprintf("[GET /timetable/fchg/{evaNo}][%d] getTimetableFchgEvaNoOK  %+v", 200, o.Payload)
}

func (o *GetTimetableFchgEvaNoOK) GetPayload() *models.Timetable {
	return o.Payload
}

func (o *GetTimetableFchgEvaNoOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Timetable)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetTimetableFchgEvaNoNotFound creates a GetTimetableFchgEvaNoNotFound with default headers values
func NewGetTimetableFchgEvaNoNotFound() *GetTimetableFchgEvaNoNotFound {
	return &GetTimetableFchgEvaNoNotFound{}
}

/*GetTimetableFchgEvaNoNotFound handles this case with default header values.

resource not found
*/
type GetTimetableFchgEvaNoNotFound struct {
}

func (o *GetTimetableFchgEvaNoNotFound) Error() string {
	return fmt.Sprintf("[GET /timetable/fchg/{evaNo}][%d] getTimetableFchgEvaNoNotFound ", 404)
}

func (o *GetTimetableFchgEvaNoNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
