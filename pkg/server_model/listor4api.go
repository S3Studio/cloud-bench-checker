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

// Listor4api listor4api
//
// swagger:model listor4api
type Listor4api struct {

	// cloud type
	CloudType Cloudtype4api `json:"cloud_type,omitempty"`

	// hash
	Hash *ItemHash `json:"hash,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// rs type
	RsType string `json:"rs_type,omitempty"`

	// yaml
	Yaml string `json:"yaml,omitempty"`

	// yaml hidden
	YamlHidden bool `json:"yaml_hidden"`
}

// Validate validates this listor4api
func (m *Listor4api) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCloudType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHash(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Listor4api) validateCloudType(formats strfmt.Registry) error {
	if swag.IsZero(m.CloudType) { // not required
		return nil
	}

	if err := m.CloudType.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("cloud_type")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("cloud_type")
		}
		return err
	}

	return nil
}

func (m *Listor4api) validateHash(formats strfmt.Registry) error {
	if swag.IsZero(m.Hash) { // not required
		return nil
	}

	if m.Hash != nil {
		if err := m.Hash.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("hash")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("hash")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this listor4api based on the context it is used
func (m *Listor4api) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCloudType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateHash(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Listor4api) contextValidateCloudType(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.CloudType) { // not required
		return nil
	}

	if err := m.CloudType.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("cloud_type")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("cloud_type")
		}
		return err
	}

	return nil
}

func (m *Listor4api) contextValidateHash(ctx context.Context, formats strfmt.Registry) error {

	if m.Hash != nil {

		if swag.IsZero(m.Hash) { // not required
			return nil
		}

		if err := m.Hash.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("hash")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("hash")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Listor4api) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Listor4api) UnmarshalBinary(b []byte) error {
	var res Listor4api
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
