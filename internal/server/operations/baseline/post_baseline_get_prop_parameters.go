// Code generated by go-swagger; DO NOT EDIT.

package baseline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/s3studio/cloud-bench-checker/pkg/server_model"
)

// NewPostBaselineGetPropParams creates a new PostBaselineGetPropParams object
//
// There are no default values defined in the spec.
func NewPostBaselineGetPropParams() PostBaselineGetPropParams {

	return PostBaselineGetPropParams{}
}

// PostBaselineGetPropParams contains all the bound params for the post baseline get prop operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostBaselineGetProp
type PostBaselineGetPropParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Id of Baseline
	  Required: true
	  In: query
	*/
	ID int64
	/*List of raw data from Listor
	  Required: true
	  In: body
	*/
	ListorData []*server_model.ListorData
	/*Name of authentication profile
	  Required: true
	  In: header
	*/
	Profile string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostBaselineGetPropParams() beforehand.
func (o *PostBaselineGetPropParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qID, qhkID, _ := qs.GetOK("id")
	if err := o.bindID(qID, qhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body []*server_model.ListorData
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("listorData", "body", ""))
			} else {
				res = append(res, errors.NewParseError("listorData", "body", "", err))
			}
		} else {

			// validate array of body objects
			for i := range body {
				if body[i] == nil {
					continue
				}
				if err := body[i].Validate(route.Formats); err != nil {
					res = append(res, err)
					break
				}
			}

			if len(res) == 0 {
				o.ListorData = body
			}
		}
	} else {
		res = append(res, errors.Required("listorData", "body", ""))
	}

	if err := o.bindProfile(r.Header[http.CanonicalHeaderKey("profile")], true, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindID binds and validates parameter ID from query.
func (o *PostBaselineGetPropParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("id", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("id", "query", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("id", "query", "int64", raw)
	}
	o.ID = value

	return nil
}

// bindProfile binds and validates parameter Profile from header.
func (o *PostBaselineGetPropParams) bindProfile(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("profile", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("profile", "header", raw); err != nil {
		return err
	}
	o.Profile = raw

	return nil
}
