// Code generated by go-swagger; DO NOT EDIT.

package listor

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// NewGetListorGetIdsParams creates a new GetListorGetIdsParams object
//
// There are no default values defined in the spec.
func NewGetListorGetIdsParams() GetListorGetIdsParams {

	return GetListorGetIdsParams{}
}

// GetListorGetIdsParams contains all the bound params for the get listor get ids operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetListorGetIds
type GetListorGetIdsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Cloud type to filter
	  In: query
	*/
	CloudType *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetListorGetIdsParams() beforehand.
func (o *GetListorGetIdsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qCloudType, qhkCloudType, _ := qs.GetOK("cloud_type")
	if err := o.bindCloudType(qCloudType, qhkCloudType, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindCloudType binds and validates parameter CloudType from query.
func (o *GetListorGetIdsParams) bindCloudType(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.CloudType = &raw

	return nil
}
