// Code generated by go-swagger; DO NOT EDIT.

package p_cloud_networks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/IBM-Cloud/power-go-client/power/models"
)

// PcloudNetworksDeleteReader is a Reader for the PcloudNetworksDelete structure.
type PcloudNetworksDeleteReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PcloudNetworksDeleteReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewPcloudNetworksDeleteOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewPcloudNetworksDeleteBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 410:
		result := NewPcloudNetworksDeleteGone()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 500:
		result := NewPcloudNetworksDeleteInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPcloudNetworksDeleteOK creates a PcloudNetworksDeleteOK with default headers values
func NewPcloudNetworksDeleteOK() *PcloudNetworksDeleteOK {
	return &PcloudNetworksDeleteOK{}
}

/*PcloudNetworksDeleteOK handles this case with default header values.

OK
*/
type PcloudNetworksDeleteOK struct {
	Payload models.Object
}

func (o *PcloudNetworksDeleteOK) Error() string {
	return fmt.Sprintf("[DELETE /pcloud/v1/cloud-instances/{cloud_instance_id}/networks/{network_id}][%d] pcloudNetworksDeleteOK  %+v", 200, o.Payload)
}

func (o *PcloudNetworksDeleteOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPcloudNetworksDeleteBadRequest creates a PcloudNetworksDeleteBadRequest with default headers values
func NewPcloudNetworksDeleteBadRequest() *PcloudNetworksDeleteBadRequest {
	return &PcloudNetworksDeleteBadRequest{}
}

/*PcloudNetworksDeleteBadRequest handles this case with default header values.

Bad Request
*/
type PcloudNetworksDeleteBadRequest struct {
	Payload *models.Error
}

func (o *PcloudNetworksDeleteBadRequest) Error() string {
	return fmt.Sprintf("[DELETE /pcloud/v1/cloud-instances/{cloud_instance_id}/networks/{network_id}][%d] pcloudNetworksDeleteBadRequest  %+v", 400, o.Payload)
}

func (o *PcloudNetworksDeleteBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPcloudNetworksDeleteGone creates a PcloudNetworksDeleteGone with default headers values
func NewPcloudNetworksDeleteGone() *PcloudNetworksDeleteGone {
	return &PcloudNetworksDeleteGone{}
}

/*PcloudNetworksDeleteGone handles this case with default header values.

Gone
*/
type PcloudNetworksDeleteGone struct {
	Payload *models.Error
}

func (o *PcloudNetworksDeleteGone) Error() string {
	return fmt.Sprintf("[DELETE /pcloud/v1/cloud-instances/{cloud_instance_id}/networks/{network_id}][%d] pcloudNetworksDeleteGone  %+v", 410, o.Payload)
}

func (o *PcloudNetworksDeleteGone) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPcloudNetworksDeleteInternalServerError creates a PcloudNetworksDeleteInternalServerError with default headers values
func NewPcloudNetworksDeleteInternalServerError() *PcloudNetworksDeleteInternalServerError {
	return &PcloudNetworksDeleteInternalServerError{}
}

/*PcloudNetworksDeleteInternalServerError handles this case with default header values.

Internal Server Error
*/
type PcloudNetworksDeleteInternalServerError struct {
	Payload *models.Error
}

func (o *PcloudNetworksDeleteInternalServerError) Error() string {
	return fmt.Sprintf("[DELETE /pcloud/v1/cloud-instances/{cloud_instance_id}/networks/{network_id}][%d] pcloudNetworksDeleteInternalServerError  %+v", 500, o.Payload)
}

func (o *PcloudNetworksDeleteInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}