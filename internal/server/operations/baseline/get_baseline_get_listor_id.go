// Code generated by go-swagger; DO NOT EDIT.

package baseline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetBaselineGetListorIDHandlerFunc turns a function with the right signature into a get baseline get listor ID handler
type GetBaselineGetListorIDHandlerFunc func(GetBaselineGetListorIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetBaselineGetListorIDHandlerFunc) Handle(params GetBaselineGetListorIDParams) middleware.Responder {
	return fn(params)
}

// GetBaselineGetListorIDHandler interface for that can handle valid get baseline get listor ID params
type GetBaselineGetListorIDHandler interface {
	Handle(GetBaselineGetListorIDParams) middleware.Responder
}

// NewGetBaselineGetListorID creates a new http.Handler for the get baseline get listor ID operation
func NewGetBaselineGetListorID(ctx *middleware.Context, handler GetBaselineGetListorIDHandler) *GetBaselineGetListorID {
	return &GetBaselineGetListorID{Context: ctx, Handler: handler}
}

/*
	GetBaselineGetListorID swagger:route GET /baseline/getListorId baseline getBaselineGetListorId

Get the ids of the Listors used in all the Checkers of the Baseline
*/
type GetBaselineGetListorID struct {
	Context *middleware.Context
	Handler GetBaselineGetListorIDHandler
}

func (o *GetBaselineGetListorID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetBaselineGetListorIDParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetBaselineGetListorIDOKBody get baseline get listor ID o k body
//
// swagger:model GetBaselineGetListorIDOKBody
type GetBaselineGetListorIDOKBody struct {

	// code
	Code int64 `json:"code,omitempty"`

	// data
	Data []int64 `json:"data"`

	// msg
	Msg string `json:"msg,omitempty"`
}

// Validate validates this get baseline get listor ID o k body
func (o *GetBaselineGetListorIDOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get baseline get listor ID o k body based on context it is used
func (o *GetBaselineGetListorIDOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetBaselineGetListorIDOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetBaselineGetListorIDOKBody) UnmarshalBinary(b []byte) error {
	var res GetBaselineGetListorIDOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
