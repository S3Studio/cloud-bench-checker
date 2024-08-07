// Code generated by go-swagger; DO NOT EDIT.

package listor

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"

	"github.com/go-openapi/swag"
)

// GetListorGetDefinitionURL generates an URL for the get listor get definition operation
type GetListorGetDefinitionURL struct {
	ID       int64
	WithHash *bool
	WithYaml *bool

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetListorGetDefinitionURL) WithBasePath(bp string) *GetListorGetDefinitionURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetListorGetDefinitionURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetListorGetDefinitionURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/listor/getDefinition"

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/api"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	idQ := swag.FormatInt64(o.ID)
	if idQ != "" {
		qs.Set("id", idQ)
	}

	var withHashQ string
	if o.WithHash != nil {
		withHashQ = swag.FormatBool(*o.WithHash)
	}
	if withHashQ != "" {
		qs.Set("with_hash", withHashQ)
	}

	var withYamlQ string
	if o.WithYaml != nil {
		withYamlQ = swag.FormatBool(*o.WithYaml)
	}
	if withYamlQ != "" {
		qs.Set("with_yaml", withYamlQ)
	}

	_result.RawQuery = qs.Encode()

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetListorGetDefinitionURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetListorGetDefinitionURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetListorGetDefinitionURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetListorGetDefinitionURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetListorGetDefinitionURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *GetListorGetDefinitionURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}