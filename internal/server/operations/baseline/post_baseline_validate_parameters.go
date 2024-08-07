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

// NewPostBaselineValidateParams creates a new PostBaselineValidateParams object
//
// There are no default values defined in the spec.
func NewPostBaselineValidateParams() PostBaselineValidateParams {

	return PostBaselineValidateParams{}
}

// PostBaselineValidateParams contains all the bound params for the post baseline validate operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostBaselineValidate
type PostBaselineValidateParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*List of properties to be validated
	  Required: true
	  In: body
	*/
	Data *server_model.BaselineData
	/*Id of Baseline
	  Required: true
	  In: query
	*/
	ID int64
	/*Metadata of Baseline to be outputted
	  In: query
	  Collection Format: multi
	*/
	Metadata []string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostBaselineValidateParams() beforehand.
func (o *PostBaselineValidateParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body server_model.BaselineData
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("data", "body", ""))
			} else {
				res = append(res, errors.NewParseError("data", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			ctx := validate.WithOperationRequest(r.Context())
			if err := body.ContextValidate(ctx, route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Data = &body
			}
		}
	} else {
		res = append(res, errors.Required("data", "body", ""))
	}

	qID, qhkID, _ := qs.GetOK("id")
	if err := o.bindID(qID, qhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	qMetadata, qhkMetadata, _ := qs.GetOK("metadata")
	if err := o.bindMetadata(qMetadata, qhkMetadata, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindID binds and validates parameter ID from query.
func (o *PostBaselineValidateParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindMetadata binds and validates array parameter Metadata from query.
//
// Arrays are parsed according to CollectionFormat: "multi" (defaults to "csv" when empty).
func (o *PostBaselineValidateParams) bindMetadata(rawData []string, hasKey bool, formats strfmt.Registry) error {
	// CollectionFormat: multi
	metadataIC := rawData
	if len(metadataIC) == 0 {
		return nil
	}

	var metadataIR []string
	for _, metadataIV := range metadataIC {
		metadataI := metadataIV

		metadataIR = append(metadataIR, metadataI)
	}

	o.Metadata = metadataIR

	return nil
}
