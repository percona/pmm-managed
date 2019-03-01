// Code generated by go-swagger; DO NOT EDIT.

package metrics

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

// NewGetMetricsByDigestParams creates a new GetMetricsByDigestParams object
// with the default values initialized.
func NewGetMetricsByDigestParams() *GetMetricsByDigestParams {
	var ()
	return &GetMetricsByDigestParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetMetricsByDigestParamsWithTimeout creates a new GetMetricsByDigestParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetMetricsByDigestParamsWithTimeout(timeout time.Duration) *GetMetricsByDigestParams {
	var ()
	return &GetMetricsByDigestParams{

		timeout: timeout,
	}
}

// NewGetMetricsByDigestParamsWithContext creates a new GetMetricsByDigestParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetMetricsByDigestParamsWithContext(ctx context.Context) *GetMetricsByDigestParams {
	var ()
	return &GetMetricsByDigestParams{

		Context: ctx,
	}
}

// NewGetMetricsByDigestParamsWithHTTPClient creates a new GetMetricsByDigestParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetMetricsByDigestParamsWithHTTPClient(client *http.Client) *GetMetricsByDigestParams {
	var ()
	return &GetMetricsByDigestParams{
		HTTPClient: client,
	}
}

/*GetMetricsByDigestParams contains all the parameters to send to the API endpoint
for the get metrics by digest operation typically these are written to a http.Request
*/
type GetMetricsByDigestParams struct {

	/*Body*/
	Body GetMetricsByDigestBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get metrics by digest params
func (o *GetMetricsByDigestParams) WithTimeout(timeout time.Duration) *GetMetricsByDigestParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get metrics by digest params
func (o *GetMetricsByDigestParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get metrics by digest params
func (o *GetMetricsByDigestParams) WithContext(ctx context.Context) *GetMetricsByDigestParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get metrics by digest params
func (o *GetMetricsByDigestParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get metrics by digest params
func (o *GetMetricsByDigestParams) WithHTTPClient(client *http.Client) *GetMetricsByDigestParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get metrics by digest params
func (o *GetMetricsByDigestParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the get metrics by digest params
func (o *GetMetricsByDigestParams) WithBody(body GetMetricsByDigestBody) *GetMetricsByDigestParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the get metrics by digest params
func (o *GetMetricsByDigestParams) SetBody(body GetMetricsByDigestBody) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *GetMetricsByDigestParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
