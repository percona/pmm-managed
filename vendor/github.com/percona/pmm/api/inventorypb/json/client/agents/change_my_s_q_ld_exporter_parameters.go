// Code generated by go-swagger; DO NOT EDIT.

package agents

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewChangeMySQLdExporterParams creates a new ChangeMySQLdExporterParams object
// with the default values initialized.
func NewChangeMySQLdExporterParams() *ChangeMySQLdExporterParams {
	var ()
	return &ChangeMySQLdExporterParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewChangeMySQLdExporterParamsWithTimeout creates a new ChangeMySQLdExporterParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewChangeMySQLdExporterParamsWithTimeout(timeout time.Duration) *ChangeMySQLdExporterParams {
	var ()
	return &ChangeMySQLdExporterParams{

		timeout: timeout,
	}
}

// NewChangeMySQLdExporterParamsWithContext creates a new ChangeMySQLdExporterParams object
// with the default values initialized, and the ability to set a context for a request
func NewChangeMySQLdExporterParamsWithContext(ctx context.Context) *ChangeMySQLdExporterParams {
	var ()
	return &ChangeMySQLdExporterParams{

		Context: ctx,
	}
}

// NewChangeMySQLdExporterParamsWithHTTPClient creates a new ChangeMySQLdExporterParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewChangeMySQLdExporterParamsWithHTTPClient(client *http.Client) *ChangeMySQLdExporterParams {
	var ()
	return &ChangeMySQLdExporterParams{
		HTTPClient: client,
	}
}

/*ChangeMySQLdExporterParams contains all the parameters to send to the API endpoint
for the change my s q ld exporter operation typically these are written to a http.Request
*/
type ChangeMySQLdExporterParams struct {

	/*Body*/
	Body ChangeMySQLdExporterBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the change my s q ld exporter params
func (o *ChangeMySQLdExporterParams) WithTimeout(timeout time.Duration) *ChangeMySQLdExporterParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the change my s q ld exporter params
func (o *ChangeMySQLdExporterParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the change my s q ld exporter params
func (o *ChangeMySQLdExporterParams) WithContext(ctx context.Context) *ChangeMySQLdExporterParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the change my s q ld exporter params
func (o *ChangeMySQLdExporterParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the change my s q ld exporter params
func (o *ChangeMySQLdExporterParams) WithHTTPClient(client *http.Client) *ChangeMySQLdExporterParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the change my s q ld exporter params
func (o *ChangeMySQLdExporterParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the change my s q ld exporter params
func (o *ChangeMySQLdExporterParams) WithBody(body ChangeMySQLdExporterBody) *ChangeMySQLdExporterParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the change my s q ld exporter params
func (o *ChangeMySQLdExporterParams) SetBody(body ChangeMySQLdExporterBody) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *ChangeMySQLdExporterParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
