// Code generated by go-swagger; DO NOT EDIT.

package baseline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostBaselineValidateOKCode is the HTTP code returned for type PostBaselineValidateOK
const PostBaselineValidateOKCode int = 200

/*
PostBaselineValidateOK List of validation results

swagger:response postBaselineValidateOK
*/
type PostBaselineValidateOK struct {

	/*
	  In: Body
	*/
	Payload *PostBaselineValidateOKBody `json:"body,omitempty"`
}

// NewPostBaselineValidateOK creates PostBaselineValidateOK with default headers values
func NewPostBaselineValidateOK() *PostBaselineValidateOK {

	return &PostBaselineValidateOK{}
}

// WithPayload adds the payload to the post baseline validate o k response
func (o *PostBaselineValidateOK) WithPayload(payload *PostBaselineValidateOKBody) *PostBaselineValidateOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post baseline validate o k response
func (o *PostBaselineValidateOK) SetPayload(payload *PostBaselineValidateOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostBaselineValidateOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostBaselineValidateBadRequestCode is the HTTP code returned for type PostBaselineValidateBadRequest
const PostBaselineValidateBadRequestCode int = 400

/*
PostBaselineValidateBadRequest Error occurs

swagger:response postBaselineValidateBadRequest
*/
type PostBaselineValidateBadRequest struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewPostBaselineValidateBadRequest creates PostBaselineValidateBadRequest with default headers values
func NewPostBaselineValidateBadRequest() *PostBaselineValidateBadRequest {

	return &PostBaselineValidateBadRequest{}
}

// WithPayload adds the payload to the post baseline validate bad request response
func (o *PostBaselineValidateBadRequest) WithPayload(payload interface{}) *PostBaselineValidateBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post baseline validate bad request response
func (o *PostBaselineValidateBadRequest) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostBaselineValidateBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// PostBaselineValidateNotFoundCode is the HTTP code returned for type PostBaselineValidateNotFound
const PostBaselineValidateNotFoundCode int = 404

/*
PostBaselineValidateNotFound Id not found

swagger:response postBaselineValidateNotFound
*/
type PostBaselineValidateNotFound struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewPostBaselineValidateNotFound creates PostBaselineValidateNotFound with default headers values
func NewPostBaselineValidateNotFound() *PostBaselineValidateNotFound {

	return &PostBaselineValidateNotFound{}
}

// WithPayload adds the payload to the post baseline validate not found response
func (o *PostBaselineValidateNotFound) WithPayload(payload interface{}) *PostBaselineValidateNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post baseline validate not found response
func (o *PostBaselineValidateNotFound) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostBaselineValidateNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
