// Code generated by go-swagger; DO NOT EDIT.

package baseline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/s3studio/cloud-bench-checker/pkg/server_model"
)

// PostBaselineGetPropHandlerFunc turns a function with the right signature into a post baseline get prop handler
type PostBaselineGetPropHandlerFunc func(PostBaselineGetPropParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostBaselineGetPropHandlerFunc) Handle(params PostBaselineGetPropParams) middleware.Responder {
	return fn(params)
}

// PostBaselineGetPropHandler interface for that can handle valid post baseline get prop params
type PostBaselineGetPropHandler interface {
	Handle(PostBaselineGetPropParams) middleware.Responder
}

// NewPostBaselineGetProp creates a new http.Handler for the post baseline get prop operation
func NewPostBaselineGetProp(ctx *middleware.Context, handler PostBaselineGetPropHandler) *PostBaselineGetProp {
	return &PostBaselineGetProp{Context: ctx, Handler: handler}
}

/*
	PostBaselineGetProp swagger:route POST /baseline/getProp baseline postBaselineGetProp

Extract properties from the raw data
*/
type PostBaselineGetProp struct {
	Context *middleware.Context
	Handler PostBaselineGetPropHandler
}

func (o *PostBaselineGetProp) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostBaselineGetPropParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PostBaselineGetPropOKBody post baseline get prop o k body
//
// swagger:model PostBaselineGetPropOKBody
type PostBaselineGetPropOKBody struct {

	// code
	Code int64 `json:"code,omitempty"`

	// data
	Data *server_model.BaselineData `json:"data,omitempty"`

	// msg
	Msg string `json:"msg,omitempty"`
}

// Validate validates this post baseline get prop o k body
func (o *PostBaselineGetPropOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostBaselineGetPropOKBody) validateData(formats strfmt.Registry) error {
	if swag.IsZero(o.Data) { // not required
		return nil
	}

	if o.Data != nil {
		if err := o.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("postBaselineGetPropOK" + "." + "data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("postBaselineGetPropOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this post baseline get prop o k body based on the context it is used
func (o *PostBaselineGetPropOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostBaselineGetPropOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if o.Data != nil {

		if swag.IsZero(o.Data) { // not required
			return nil
		}

		if err := o.Data.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("postBaselineGetPropOK" + "." + "data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("postBaselineGetPropOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostBaselineGetPropOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostBaselineGetPropOKBody) UnmarshalBinary(b []byte) error {
	var res PostBaselineGetPropOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
