// Code generated by go-swagger; DO NOT EDIT.

package alert

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

	"github.com/percona/pmm/api/alertmanager/ammodels"
)

// NewPostAlertsParams creates a new PostAlertsParams object
// with the default values initialized.
func NewPostAlertsParams() *PostAlertsParams {
	var ()
	return &PostAlertsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostAlertsParamsWithTimeout creates a new PostAlertsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostAlertsParamsWithTimeout(timeout time.Duration) *PostAlertsParams {
	var ()
	return &PostAlertsParams{

		timeout: timeout,
	}
}

// NewPostAlertsParamsWithContext creates a new PostAlertsParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostAlertsParamsWithContext(ctx context.Context) *PostAlertsParams {
	var ()
	return &PostAlertsParams{

		Context: ctx,
	}
}

// NewPostAlertsParamsWithHTTPClient creates a new PostAlertsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostAlertsParamsWithHTTPClient(client *http.Client) *PostAlertsParams {
	var ()
	return &PostAlertsParams{
		HTTPClient: client,
	}
}

/*PostAlertsParams contains all the parameters to send to the API endpoint
for the post alerts operation typically these are written to a http.Request
*/
type PostAlertsParams struct {

	/*Alerts
	  The alerts to create

	*/
	Alerts ammodels.PostableAlerts

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post alerts params
func (o *PostAlertsParams) WithTimeout(timeout time.Duration) *PostAlertsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post alerts params
func (o *PostAlertsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post alerts params
func (o *PostAlertsParams) WithContext(ctx context.Context) *PostAlertsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post alerts params
func (o *PostAlertsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post alerts params
func (o *PostAlertsParams) WithHTTPClient(client *http.Client) *PostAlertsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post alerts params
func (o *PostAlertsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAlerts adds the alerts to the post alerts params
func (o *PostAlertsParams) WithAlerts(alerts ammodels.PostableAlerts) *PostAlertsParams {
	o.SetAlerts(alerts)
	return o
}

// SetAlerts adds the alerts to the post alerts params
func (o *PostAlertsParams) SetAlerts(alerts ammodels.PostableAlerts) {
	o.Alerts = alerts
}

// WriteToRequest writes these params to a swagger request
func (o *PostAlertsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Alerts != nil {
		if err := r.SetBodyParam(o.Alerts); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
