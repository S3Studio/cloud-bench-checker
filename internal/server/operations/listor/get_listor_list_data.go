// Code generated by go-swagger; DO NOT EDIT.

package listor

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

// GetListorListDataHandlerFunc turns a function with the right signature into a get listor list data handler
type GetListorListDataHandlerFunc func(GetListorListDataParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetListorListDataHandlerFunc) Handle(params GetListorListDataParams) middleware.Responder {
	return fn(params)
}

// GetListorListDataHandler interface for that can handle valid get listor list data params
type GetListorListDataHandler interface {
	Handle(GetListorListDataParams) middleware.Responder
}

// NewGetListorListData creates a new http.Handler for the get listor list data operation
func NewGetListorListData(ctx *middleware.Context, handler GetListorListDataHandler) *GetListorListData {
	return &GetListorListData{Context: ctx, Handler: handler}
}

/*
	GetListorListData swagger:route GET /listor/listData listor getListorListData

Get list of all raw data according to the definition
*/
type GetListorListData struct {
	Context *middleware.Context
	Handler GetListorListDataHandler
}

func (o *GetListorListData) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetListorListDataParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetListorListDataOKBody get listor list data o k body
//
// swagger:model GetListorListDataOKBody
type GetListorListDataOKBody struct {

	// code
	Code int64 `json:"code,omitempty"`

	// data
	Data *server_model.ListorData `json:"data,omitempty"`

	// msg
	Msg string `json:"msg,omitempty"`
}

// Validate validates this get listor list data o k body
func (o *GetListorListDataOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateData(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetListorListDataOKBody) validateData(formats strfmt.Registry) error {
	if swag.IsZero(o.Data) { // not required
		return nil
	}

	if o.Data != nil {
		if err := o.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getListorListDataOK" + "." + "data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getListorListDataOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this get listor list data o k body based on the context it is used
func (o *GetListorListDataOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateData(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetListorListDataOKBody) contextValidateData(ctx context.Context, formats strfmt.Registry) error {

	if o.Data != nil {

		if swag.IsZero(o.Data) { // not required
			return nil
		}

		if err := o.Data.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("getListorListDataOK" + "." + "data")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("getListorListDataOK" + "." + "data")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetListorListDataOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetListorListDataOKBody) UnmarshalBinary(b []byte) error {
	var res GetListorListDataOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
