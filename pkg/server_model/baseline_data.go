// Code generated by go-swagger; DO NOT EDIT.

package server_model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// BaselineData baseline data
//
// swagger:model baseline_data
type BaselineData struct {

	// baseline hash
	BaselineHash *ItemHash `json:"baseline_hash,omitempty"`

	// checker prop
	CheckerProp []string `json:"checker_prop"`

	// id
	ID int64 `json:"id,omitempty"`
}

// Validate validates this baseline data
func (m *BaselineData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBaselineHash(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BaselineData) validateBaselineHash(formats strfmt.Registry) error {
	if swag.IsZero(m.BaselineHash) { // not required
		return nil
	}

	if m.BaselineHash != nil {
		if err := m.BaselineHash.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("baseline_hash")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("baseline_hash")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this baseline data based on the context it is used
func (m *BaselineData) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBaselineHash(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BaselineData) contextValidateBaselineHash(ctx context.Context, formats strfmt.Registry) error {

	if m.BaselineHash != nil {

		if swag.IsZero(m.BaselineHash) { // not required
			return nil
		}

		if err := m.BaselineHash.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("baseline_hash")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("baseline_hash")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *BaselineData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BaselineData) UnmarshalBinary(b []byte) error {
	var res BaselineData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
