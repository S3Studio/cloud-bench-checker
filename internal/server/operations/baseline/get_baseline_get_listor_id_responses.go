// Code generated by go-swagger; DO NOT EDIT.

package baseline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetBaselineGetListorIDOKCode is the HTTP code returned for type GetBaselineGetListorIDOK
const GetBaselineGetListorIDOKCode int = 200

/*
GetBaselineGetListorIDOK Ids of Listors

swagger:response getBaselineGetListorIdOK
*/
type GetBaselineGetListorIDOK struct {

	/*
	  In: Body
	*/
	Payload *GetBaselineGetListorIDOKBody `json:"body,omitempty"`
}

// NewGetBaselineGetListorIDOK creates GetBaselineGetListorIDOK with default headers values
func NewGetBaselineGetListorIDOK() *GetBaselineGetListorIDOK {

	return &GetBaselineGetListorIDOK{}
}

// WithPayload adds the payload to the get baseline get listor Id o k response
func (o *GetBaselineGetListorIDOK) WithPayload(payload *GetBaselineGetListorIDOKBody) *GetBaselineGetListorIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get baseline get listor Id o k response
func (o *GetBaselineGetListorIDOK) SetPayload(payload *GetBaselineGetListorIDOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBaselineGetListorIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetBaselineGetListorIDNotFoundCode is the HTTP code returned for type GetBaselineGetListorIDNotFound
const GetBaselineGetListorIDNotFoundCode int = 404

/*
GetBaselineGetListorIDNotFound Id not found

swagger:response getBaselineGetListorIdNotFound
*/
type GetBaselineGetListorIDNotFound struct {

	/*
	  In: Body
	*/
	Payload interface{} `json:"body,omitempty"`
}

// NewGetBaselineGetListorIDNotFound creates GetBaselineGetListorIDNotFound with default headers values
func NewGetBaselineGetListorIDNotFound() *GetBaselineGetListorIDNotFound {

	return &GetBaselineGetListorIDNotFound{}
}

// WithPayload adds the payload to the get baseline get listor Id not found response
func (o *GetBaselineGetListorIDNotFound) WithPayload(payload interface{}) *GetBaselineGetListorIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get baseline get listor Id not found response
func (o *GetBaselineGetListorIDNotFound) SetPayload(payload interface{}) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBaselineGetListorIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
