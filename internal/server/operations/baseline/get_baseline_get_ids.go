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

// GetBaselineGetIdsHandlerFunc turns a function with the right signature into a get baseline get ids handler
type GetBaselineGetIdsHandlerFunc func(GetBaselineGetIdsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetBaselineGetIdsHandlerFunc) Handle(params GetBaselineGetIdsParams) middleware.Responder {
	return fn(params)
}

// GetBaselineGetIdsHandler interface for that can handle valid get baseline get ids params
type GetBaselineGetIdsHandler interface {
	Handle(GetBaselineGetIdsParams) middleware.Responder
}

// NewGetBaselineGetIds creates a new http.Handler for the get baseline get ids operation
func NewGetBaselineGetIds(ctx *middleware.Context, handler GetBaselineGetIdsHandler) *GetBaselineGetIds {
	return &GetBaselineGetIds{Context: ctx, Handler: handler}
}

/*
	GetBaselineGetIds swagger:route GET /baseline/getIds baseline getBaselineGetIds

Get ids of Baseline
*/
type GetBaselineGetIds struct {
	Context *middleware.Context
	Handler GetBaselineGetIdsHandler
}

func (o *GetBaselineGetIds) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetBaselineGetIdsParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetBaselineGetIdsOKBody get baseline get ids o k body
//
// swagger:model GetBaselineGetIdsOKBody
type GetBaselineGetIdsOKBody struct {

	// code
	Code int64 `json:"code,omitempty"`

	// data
	Data []int64 `json:"data"`

	// msg
	Msg string `json:"msg,omitempty"`
}

// Validate validates this get baseline get ids o k body
func (o *GetBaselineGetIdsOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this get baseline get ids o k body based on context it is used
func (o *GetBaselineGetIdsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *GetBaselineGetIdsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetBaselineGetIdsOKBody) UnmarshalBinary(b []byte) error {
	var res GetBaselineGetIdsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
