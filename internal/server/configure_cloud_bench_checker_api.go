// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/s3studio/cloud-bench-checker/internal/server/operations"
	"github.com/s3studio/cloud-bench-checker/internal/server/operations/baseline"
	"github.com/s3studio/cloud-bench-checker/internal/server/operations/listor"
)

//go:generate swagger generate server --target ../../../cloud-bench-checker --name CloudBenchCheckerAPI --spec ../../doc/api_swagger.yml --model-package ./pkg/server-model --server-package ./internal/server --principal interface{}

func configureFlags(api *operations.CloudBenchCheckerAPIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.CloudBenchCheckerAPIAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	setupHandler(api)

	if api.BaselineGetBaselineGetDefinitionHandler == nil {
		api.BaselineGetBaselineGetDefinitionHandler = baseline.GetBaselineGetDefinitionHandlerFunc(func(params baseline.GetBaselineGetDefinitionParams) middleware.Responder {
			return middleware.NotImplemented("operation baseline.GetBaselineGetDefinition has not yet been implemented")
		})
	}
	if api.BaselineGetBaselineGetIdsHandler == nil {
		api.BaselineGetBaselineGetIdsHandler = baseline.GetBaselineGetIdsHandlerFunc(func(params baseline.GetBaselineGetIdsParams) middleware.Responder {
			return middleware.NotImplemented("operation baseline.GetBaselineGetIds has not yet been implemented")
		})
	}
	if api.BaselineGetBaselineGetListorIDHandler == nil {
		api.BaselineGetBaselineGetListorIDHandler = baseline.GetBaselineGetListorIDHandlerFunc(func(params baseline.GetBaselineGetListorIDParams) middleware.Responder {
			return middleware.NotImplemented("operation baseline.GetBaselineGetListorID has not yet been implemented")
		})
	}
	if api.ListorGetListorGetDefinitionHandler == nil {
		api.ListorGetListorGetDefinitionHandler = listor.GetListorGetDefinitionHandlerFunc(func(params listor.GetListorGetDefinitionParams) middleware.Responder {
			return middleware.NotImplemented("operation listor.GetListorGetDefinition has not yet been implemented")
		})
	}
	if api.ListorGetListorGetIdsHandler == nil {
		api.ListorGetListorGetIdsHandler = listor.GetListorGetIdsHandlerFunc(func(params listor.GetListorGetIdsParams) middleware.Responder {
			return middleware.NotImplemented("operation listor.GetListorGetIds has not yet been implemented")
		})
	}
	if api.ListorGetListorListDataHandler == nil {
		api.ListorGetListorListDataHandler = listor.GetListorListDataHandlerFunc(func(params listor.GetListorListDataParams) middleware.Responder {
			return middleware.NotImplemented("operation listor.GetListorListData has not yet been implemented")
		})
	}
	if api.BaselinePostBaselineGetPropHandler == nil {
		api.BaselinePostBaselineGetPropHandler = baseline.PostBaselineGetPropHandlerFunc(func(params baseline.PostBaselineGetPropParams) middleware.Responder {
			return middleware.NotImplemented("operation baseline.PostBaselineGetProp has not yet been implemented")
		})
	}
	if api.BaselinePostBaselineValidateHandler == nil {
		api.BaselinePostBaselineValidateHandler = baseline.PostBaselineValidateHandlerFunc(func(params baseline.PostBaselineValidateParams) middleware.Responder {
			return middleware.NotImplemented("operation baseline.PostBaselineValidate has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
