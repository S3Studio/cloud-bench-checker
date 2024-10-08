// Code generated by go-swagger; DO NOT EDIT.

package baseline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// NewGetBaselineGetIdsParams creates a new GetBaselineGetIdsParams object
//
// There are no default values defined in the spec.
func NewGetBaselineGetIdsParams() GetBaselineGetIdsParams {

	return GetBaselineGetIdsParams{}
}

// GetBaselineGetIdsParams contains all the bound params for the get baseline get ids operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetBaselineGetIds
type GetBaselineGetIdsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Tag to filter
	  In: query
	  Collection Format: multi
	*/
	Tag []string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetBaselineGetIdsParams() beforehand.
func (o *GetBaselineGetIdsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qTag, qhkTag, _ := qs.GetOK("tag")
	if err := o.bindTag(qTag, qhkTag, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindTag binds and validates array parameter Tag from query.
//
// Arrays are parsed according to CollectionFormat: "multi" (defaults to "csv" when empty).
func (o *GetBaselineGetIdsParams) bindTag(rawData []string, hasKey bool, formats strfmt.Registry) error {
	// CollectionFormat: multi
	tagIC := rawData
	if len(tagIC) == 0 {
		return nil
	}

	var tagIR []string
	for _, tagIV := range tagIC {
		tagI := tagIV

		tagIR = append(tagIR, tagI)
	}

	o.Tag = tagIR

	return nil
}
