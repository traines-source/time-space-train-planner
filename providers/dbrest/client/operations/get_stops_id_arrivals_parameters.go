// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewGetStopsIDArrivalsParams creates a new GetStopsIDArrivalsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetStopsIDArrivalsParams() *GetStopsIDArrivalsParams {
	return &GetStopsIDArrivalsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetStopsIDArrivalsParamsWithTimeout creates a new GetStopsIDArrivalsParams object
// with the ability to set a timeout on a request.
func NewGetStopsIDArrivalsParamsWithTimeout(timeout time.Duration) *GetStopsIDArrivalsParams {
	return &GetStopsIDArrivalsParams{
		timeout: timeout,
	}
}

// NewGetStopsIDArrivalsParamsWithContext creates a new GetStopsIDArrivalsParams object
// with the ability to set a context for a request.
func NewGetStopsIDArrivalsParamsWithContext(ctx context.Context) *GetStopsIDArrivalsParams {
	return &GetStopsIDArrivalsParams{
		Context: ctx,
	}
}

// NewGetStopsIDArrivalsParamsWithHTTPClient creates a new GetStopsIDArrivalsParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetStopsIDArrivalsParamsWithHTTPClient(client *http.Client) *GetStopsIDArrivalsParams {
	return &GetStopsIDArrivalsParams{
		HTTPClient: client,
	}
}

/* GetStopsIDArrivalsParams contains all the parameters to send to the API endpoint
   for the get stops ID arrivals operation.

   Typically these are written to a http.Request.
*/
type GetStopsIDArrivalsParams struct {

	/* Direction.

	   Filter departures by direction.
	*/
	Direction *string

	/* Duration.

	   Show departures for how many minutes?

	   Default: 10
	*/
	Duration *int64

	/* ID.

	   stop/station ID to show arrivals for
	*/
	ID string

	/* Language.

	   Language of the results.

	   Default: "en"
	*/
	Language *string

	/* LinesOfStops.

	   Parse & return lines of each stop/station?
	*/
	LinesOfStops *bool

	/* Remarks.

	   Parse & return hints & warnings?

	   Default: true
	*/
	Remarks *bool

	/* Results.

	   Max. number of departures. – Default: *whatever HAFAS wants*
	*/
	Results *int64

	/* When.

	   Date & time to get departures for. – Default: *now*

	   Format: date-time
	*/
	When *strfmt.DateTime

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get stops ID arrivals params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetStopsIDArrivalsParams) WithDefaults() *GetStopsIDArrivalsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get stops ID arrivals params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetStopsIDArrivalsParams) SetDefaults() {
	var (
		durationDefault = int64(10)

		languageDefault = string("en")

		linesOfStopsDefault = bool(false)

		remarksDefault = bool(true)
	)

	val := GetStopsIDArrivalsParams{
		Duration:     &durationDefault,
		Language:     &languageDefault,
		LinesOfStops: &linesOfStopsDefault,
		Remarks:      &remarksDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) WithTimeout(timeout time.Duration) *GetStopsIDArrivalsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) WithContext(ctx context.Context) *GetStopsIDArrivalsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) WithHTTPClient(client *http.Client) *GetStopsIDArrivalsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDirection adds the direction to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) WithDirection(direction *string) *GetStopsIDArrivalsParams {
	o.SetDirection(direction)
	return o
}

// SetDirection adds the direction to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) SetDirection(direction *string) {
	o.Direction = direction
}

// WithDuration adds the duration to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) WithDuration(duration *int64) *GetStopsIDArrivalsParams {
	o.SetDuration(duration)
	return o
}

// SetDuration adds the duration to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) SetDuration(duration *int64) {
	o.Duration = duration
}

// WithID adds the id to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) WithID(id string) *GetStopsIDArrivalsParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) SetID(id string) {
	o.ID = id
}

// WithLanguage adds the language to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) WithLanguage(language *string) *GetStopsIDArrivalsParams {
	o.SetLanguage(language)
	return o
}

// SetLanguage adds the language to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) SetLanguage(language *string) {
	o.Language = language
}

// WithLinesOfStops adds the linesOfStops to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) WithLinesOfStops(linesOfStops *bool) *GetStopsIDArrivalsParams {
	o.SetLinesOfStops(linesOfStops)
	return o
}

// SetLinesOfStops adds the linesOfStops to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) SetLinesOfStops(linesOfStops *bool) {
	o.LinesOfStops = linesOfStops
}

// WithRemarks adds the remarks to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) WithRemarks(remarks *bool) *GetStopsIDArrivalsParams {
	o.SetRemarks(remarks)
	return o
}

// SetRemarks adds the remarks to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) SetRemarks(remarks *bool) {
	o.Remarks = remarks
}

// WithResults adds the results to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) WithResults(results *int64) *GetStopsIDArrivalsParams {
	o.SetResults(results)
	return o
}

// SetResults adds the results to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) SetResults(results *int64) {
	o.Results = results
}

// WithWhen adds the when to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) WithWhen(when *strfmt.DateTime) *GetStopsIDArrivalsParams {
	o.SetWhen(when)
	return o
}

// SetWhen adds the when to the get stops ID arrivals params
func (o *GetStopsIDArrivalsParams) SetWhen(when *strfmt.DateTime) {
	o.When = when
}

// WriteToRequest writes these params to a swagger request
func (o *GetStopsIDArrivalsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Direction != nil {

		// query param direction
		var qrDirection string

		if o.Direction != nil {
			qrDirection = *o.Direction
		}
		qDirection := qrDirection
		if qDirection != "" {

			if err := r.SetQueryParam("direction", qDirection); err != nil {
				return err
			}
		}
	}

	if o.Duration != nil {

		// query param duration
		var qrDuration int64

		if o.Duration != nil {
			qrDuration = *o.Duration
		}
		qDuration := swag.FormatInt64(qrDuration)
		if qDuration != "" {

			if err := r.SetQueryParam("duration", qDuration); err != nil {
				return err
			}
		}
	}

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if o.Language != nil {

		// query param language
		var qrLanguage string

		if o.Language != nil {
			qrLanguage = *o.Language
		}
		qLanguage := qrLanguage
		if qLanguage != "" {

			if err := r.SetQueryParam("language", qLanguage); err != nil {
				return err
			}
		}
	}

	if o.LinesOfStops != nil {

		// query param linesOfStops
		var qrLinesOfStops bool

		if o.LinesOfStops != nil {
			qrLinesOfStops = *o.LinesOfStops
		}
		qLinesOfStops := swag.FormatBool(qrLinesOfStops)
		if qLinesOfStops != "" {

			if err := r.SetQueryParam("linesOfStops", qLinesOfStops); err != nil {
				return err
			}
		}
	}

	if o.Remarks != nil {

		// query param remarks
		var qrRemarks bool

		if o.Remarks != nil {
			qrRemarks = *o.Remarks
		}
		qRemarks := swag.FormatBool(qrRemarks)
		if qRemarks != "" {

			if err := r.SetQueryParam("remarks", qRemarks); err != nil {
				return err
			}
		}
	}

	if o.Results != nil {

		// query param results
		var qrResults int64

		if o.Results != nil {
			qrResults = *o.Results
		}
		qResults := swag.FormatInt64(qrResults)
		if qResults != "" {

			if err := r.SetQueryParam("results", qResults); err != nil {
				return err
			}
		}
	}

	if o.When != nil {

		// query param when
		var qrWhen strfmt.DateTime

		if o.When != nil {
			qrWhen = *o.When
		}
		qWhen := qrWhen.String()
		if qWhen != "" {

			if err := r.SetQueryParam("when", qWhen); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
