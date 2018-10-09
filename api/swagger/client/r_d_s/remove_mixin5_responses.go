// Code generated by go-swagger; DO NOT EDIT.

package r_d_s

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/percona/pmm-managed/api/swagger/models"
)

// RemoveMixin5Reader is a Reader for the RemoveMixin5 structure.
type RemoveMixin5Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RemoveMixin5Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewRemoveMixin5OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewRemoveMixin5OK creates a RemoveMixin5OK with default headers values
func NewRemoveMixin5OK() *RemoveMixin5OK {
	return &RemoveMixin5OK{}
}

/*RemoveMixin5OK handles this case with default header values.

(empty)
*/
type RemoveMixin5OK struct {
	Payload models.APIRDSRemoveResponse
}

func (o *RemoveMixin5OK) Error() string {
	return fmt.Sprintf("[DELETE /v0/rds][%d] removeMixin5OK  %+v", 200, o.Payload)
}

func (o *RemoveMixin5OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
