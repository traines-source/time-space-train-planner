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

// NewGetLocationsParams creates a new GetLocationsParams object
// with the default values initialized.
func NewGetLocationsParams() *GetLocationsParams {
	var (
		addressesDefault    = bool(true)
		fuzzyDefault        = bool(true)
		languageDefault     = string("en")
		linesOfStopsDefault = bool(false)
		poiDefault          = bool(true)
		resultsDefault      = int64(10)
		stopsDefault        = bool(true)
	)
	return &GetLocationsParams{
		Addresses:    &addressesDefault,
		Fuzzy:        &fuzzyDefault,
		Language:     &languageDefault,
		LinesOfStops: &linesOfStopsDefault,
		Poi:          &poiDefault,
		Results:      &resultsDefault,
		Stops:        &stopsDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewGetLocationsParamsWithTimeout creates a new GetLocationsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetLocationsParamsWithTimeout(timeout time.Duration) *GetLocationsParams {
	var (
		addressesDefault    = bool(true)
		fuzzyDefault        = bool(true)
		languageDefault     = string("en")
		linesOfStopsDefault = bool(false)
		poiDefault          = bool(true)
		resultsDefault      = int64(10)
		stopsDefault        = bool(true)
	)
	return &GetLocationsParams{
		Addresses:    &addressesDefault,
		Fuzzy:        &fuzzyDefault,
		Language:     &languageDefault,
		LinesOfStops: &linesOfStopsDefault,
		Poi:          &poiDefault,
		Results:      &resultsDefault,
		Stops:        &stopsDefault,

		timeout: timeout,
	}
}

// NewGetLocationsParamsWithContext creates a new GetLocationsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetLocationsParamsWithContext(ctx context.Context) *GetLocationsParams {
	var (
		addressesDefault    = bool(true)
		fuzzyDefault        = bool(true)
		languageDefault     = string("en")
		linesOfStopsDefault = bool(false)
		poiDefault          = bool(true)
		resultsDefault      = int64(10)
		stopsDefault        = bool(true)
	)
	return &GetLocationsParams{
		Addresses:    &addressesDefault,
		Fuzzy:        &fuzzyDefault,
		Language:     &languageDefault,
		LinesOfStops: &linesOfStopsDefault,
		Poi:          &poiDefault,
		Results:      &resultsDefault,
		Stops:        &stopsDefault,

		Context: ctx,
	}
}

// NewGetLocationsParamsWithHTTPClient creates a new GetLocationsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetLocationsParamsWithHTTPClient(client *http.Client) *GetLocationsParams {
	var (
		addressesDefault    = bool(true)
		fuzzyDefault        = bool(true)
		languageDefault     = string("en")
		linesOfStopsDefault = bool(false)
		poiDefault          = bool(true)
		resultsDefault      = int64(10)
		stopsDefault        = bool(true)
	)
	return &GetLocationsParams{
		Addresses:    &addressesDefault,
		Fuzzy:        &fuzzyDefault,
		Language:     &languageDefault,
		LinesOfStops: &linesOfStopsDefault,
		Poi:          &poiDefault,
		Results:      &resultsDefault,
		Stops:        &stopsDefault,
		HTTPClient:   client,
	}
}

/*GetLocationsParams contains all the parameters to send to the API endpoint
for the get locations operation typically these are written to a http.Request
*/
type GetLocationsParams struct {

	/*Addresses
	  Show points of interest?

	*/
	Addresses *bool
	/*Fuzzy
	  Find more than exact matches?

	*/
	Fuzzy *bool
	/*Language
	  Language of the results.

	*/
	Language *string
	/*LinesOfStops
	  Parse & return lines of each stop/station?

	*/
	LinesOfStops *bool
	/*Poi
	  Show addresses?

	*/
	Poi *bool
	/*Query
	  The term to search for.

	*/
	Query string
	/*Results
	  How many stations shall be shown?

	*/
	Results *int64
	/*Stops
	  Show stops/stations?

	*/
	Stops *bool

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get locations params
func (o *GetLocationsParams) WithTimeout(timeout time.Duration) *GetLocationsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get locations params
func (o *GetLocationsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get locations params
func (o *GetLocationsParams) WithContext(ctx context.Context) *GetLocationsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get locations params
func (o *GetLocationsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get locations params
func (o *GetLocationsParams) WithHTTPClient(client *http.Client) *GetLocationsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get locations params
func (o *GetLocationsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAddresses adds the addresses to the get locations params
func (o *GetLocationsParams) WithAddresses(addresses *bool) *GetLocationsParams {
	o.SetAddresses(addresses)
	return o
}

// SetAddresses adds the addresses to the get locations params
func (o *GetLocationsParams) SetAddresses(addresses *bool) {
	o.Addresses = addresses
}

// WithFuzzy adds the fuzzy to the get locations params
func (o *GetLocationsParams) WithFuzzy(fuzzy *bool) *GetLocationsParams {
	o.SetFuzzy(fuzzy)
	return o
}

// SetFuzzy adds the fuzzy to the get locations params
func (o *GetLocationsParams) SetFuzzy(fuzzy *bool) {
	o.Fuzzy = fuzzy
}

// WithLanguage adds the language to the get locations params
func (o *GetLocationsParams) WithLanguage(language *string) *GetLocationsParams {
	o.SetLanguage(language)
	return o
}

// SetLanguage adds the language to the get locations params
func (o *GetLocationsParams) SetLanguage(language *string) {
	o.Language = language
}

// WithLinesOfStops adds the linesOfStops to the get locations params
func (o *GetLocationsParams) WithLinesOfStops(linesOfStops *bool) *GetLocationsParams {
	o.SetLinesOfStops(linesOfStops)
	return o
}

// SetLinesOfStops adds the linesOfStops to the get locations params
func (o *GetLocationsParams) SetLinesOfStops(linesOfStops *bool) {
	o.LinesOfStops = linesOfStops
}

// WithPoi adds the poi to the get locations params
func (o *GetLocationsParams) WithPoi(poi *bool) *GetLocationsParams {
	o.SetPoi(poi)
	return o
}

// SetPoi adds the poi to the get locations params
func (o *GetLocationsParams) SetPoi(poi *bool) {
	o.Poi = poi
}

// WithQuery adds the query to the get locations params
func (o *GetLocationsParams) WithQuery(query string) *GetLocationsParams {
	o.SetQuery(query)
	return o
}

// SetQuery adds the query to the get locations params
func (o *GetLocationsParams) SetQuery(query string) {
	o.Query = query
}

// WithResults adds the results to the get locations params
func (o *GetLocationsParams) WithResults(results *int64) *GetLocationsParams {
	o.SetResults(results)
	return o
}

// SetResults adds the results to the get locations params
func (o *GetLocationsParams) SetResults(results *int64) {
	o.Results = results
}

// WithStops adds the stops to the get locations params
func (o *GetLocationsParams) WithStops(stops *bool) *GetLocationsParams {
	o.SetStops(stops)
	return o
}

// SetStops adds the stops to the get locations params
func (o *GetLocationsParams) SetStops(stops *bool) {
	o.Stops = stops
}

// WriteToRequest writes these params to a swagger request
func (o *GetLocationsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Addresses != nil {

		// query param addresses
		var qrAddresses bool
		if o.Addresses != nil {
			qrAddresses = *o.Addresses
		}
		qAddresses := swag.FormatBool(qrAddresses)
		if qAddresses != "" {
			if err := r.SetQueryParam("addresses", qAddresses); err != nil {
				return err
			}
		}

	}

	if o.Fuzzy != nil {

		// query param fuzzy
		var qrFuzzy bool
		if o.Fuzzy != nil {
			qrFuzzy = *o.Fuzzy
		}
		qFuzzy := swag.FormatBool(qrFuzzy)
		if qFuzzy != "" {
			if err := r.SetQueryParam("fuzzy", qFuzzy); err != nil {
				return err
			}
		}

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

	if o.Poi != nil {

		// query param poi
		var qrPoi bool
		if o.Poi != nil {
			qrPoi = *o.Poi
		}
		qPoi := swag.FormatBool(qrPoi)
		if qPoi != "" {
			if err := r.SetQueryParam("poi", qPoi); err != nil {
				return err
			}
		}

	}

	// query param query
	qrQuery := o.Query
	qQuery := qrQuery
	if qQuery != "" {
		if err := r.SetQueryParam("query", qQuery); err != nil {
			return err
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

	if o.Stops != nil {

		// query param stops
		var qrStops bool
		if o.Stops != nil {
			qrStops = *o.Stops
		}
		qStops := swag.FormatBool(qrStops)
		if qStops != "" {
			if err := r.SetQueryParam("stops", qStops); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}