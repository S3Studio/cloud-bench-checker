// Connector for Azure

package connector

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"unsafe"

	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	"github.com/s3studio/cloud-bench-checker/test"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/agiledragon/gomonkey/v2"
)

func setupEnvAzure() {
	envMap := map[string]string{
		"AZURE_CLIENT_ID":       "mock_clientid",
		"AZURE_TENANT_ID":       "mock_tenantid",
		"AZURE_CLIENT_SECRET":   "mock_secret",
		"AZURE_SUBSCRIPTION_ID": "mock_subid",
	}
	for k, v := range envMap {
		os.Setenv(k, v)
	}
}

func setupAzureCredential() func() {
	patchNewCredential := gomonkey.ApplyFunc(azidentity.NewClientSecretCredential,
		func(tenantID string, clientID string, clientSecret string, options *azidentity.ClientSecretCredentialOptions) (*azidentity.ClientSecretCredential, error) {
			return &azidentity.ClientSecretCredential{}, nil
		})
	return patchNewCredential.Reset
}

func Test_createAzureClient(t *testing.T) {
	setupEnvAzure()
	deferFn := setupAzureCredential()
	defer deferFn()

	type args struct {
		p auth.IAuthProvider
	}
	tests := []struct {
		name string
		args args
		//want    *arm.Client
		wantErr bool
	}{
		{
			"Valid result",
			args{auth.NewAuthFileProvider(test.Test_conf_azure)},
			false,
		},
		{
			"Profile not defined",
			args{auth.NewAuthFileProvider(test.Test_conf_invalid)},
			true,
		},
		{
			"Key not set",
			args{&test.MockKeyNotSetAuthProvider{}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createAzureClient(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("createAzureClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("createAzureClient() = %v, want a valid pointer", got)
			}
		})
	}
}

func Test_getAzureClient(t *testing.T) {
	setupEnvAzure()
	deferFn := setupAzureCredential()
	defer deferFn()

	type args struct {
		p auth.IAuthProvider
	}
	tests := []struct {
		name string
		args args
		//want    *arm.Client
		wantErr bool
	}{
		{
			"Valid result",
			args{auth.NewAuthFileProvider(test.Test_conf_azure)},
			false,
		},
		{
			"Profile not defined",
			args{auth.NewAuthFileProvider(test.Test_conf_invalid)},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAzureClient(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAzureClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("getAzureClient() = %v, want a valid pointer", got)
			}
		})
	}
}

type mockAzureClient struct {
	Ep string
	Pl runtime.Pipeline
	Tr tracing.Tracer
}

func setupGetAzureClient(valValid []byte) func() {
	patchGetAzureClient := gomonkey.ApplyFunc(getAzureClient,
		func(p auth.IAuthProvider) (*arm.Client, error) {
			c := mockAzureClient{
				Ep: "https://mock.domain",
				Pl: runtime.Pipeline{},
			}
			return (*arm.Client)(unsafe.Pointer(&c)), nil
		})

	patchDo := gomonkey.ApplyMethodFunc(runtime.Pipeline{}, "Do",
		func(req *policy.Request) (*http.Response, error) {
			rr := httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer(valValid),
			}
			return rr.Result(), nil
		})

	return func() {
		patchDo.Reset()
		patchGetAzureClient.Reset()
	}
}

func TestCallAzureList(t *testing.T) {
	setupEnvAzure()
	deferFn := setupAzureCredential()
	defer deferFn()
	valValid := []byte("{}")
	var rmValid json.RawMessage = valValid
	deferFn2 := setupGetAzureClient(valValid)
	defer deferFn2()

	type args struct {
		authProvider auth.IAuthProvider
		provider     string
		version      string
		rsType       string
		nextLink     string
	}
	tests := []struct {
		name    string
		args    args
		want    *json.RawMessage
		wantErr bool
	}{
		{
			"Valid result",
			args{
				auth.NewAuthFileProvider(test.Test_conf_azure),
				"", "", "", "",
			},
			&rmValid,
			false,
		},
		{
			"Valid result with endpoint",
			args{
				auth.NewAuthFileProvider(test.Test_conf_azure),
				"", "", "", "mock",
			},
			&rmValid,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CallAzureList(tt.args.authProvider, tt.args.provider, tt.args.version, tt.args.rsType, tt.args.nextLink)
			if (err != nil) != tt.wantErr {
				t.Errorf("CallAzureList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CallAzureList() = %v, want %v", string(*got), string(*tt.want))
			}
		})
	}
}

func TestCallAzureWithEndpoint(t *testing.T) {
	setupEnvAzure()
	deferFn := setupAzureCredential()
	defer deferFn()
	valValid := []byte("{}")
	var rmValid json.RawMessage = valValid
	deferFn2 := setupGetAzureClient(valValid)
	defer deferFn2()

	type args struct {
		authProvider auth.IAuthProvider
		version      string
		endpoint     string
		action       string
	}
	tests := []struct {
		name    string
		args    args
		want    *json.RawMessage
		wantErr bool
	}{
		{
			"Valid result",
			args{
				auth.NewAuthFileProvider(test.Test_conf_azure),
				"", "mock", "",
			},
			&rmValid,
			false,
		},
		{
			"Valid result with prefix of http",
			args{
				auth.NewAuthFileProvider(test.Test_conf_azure),
				"", "https://mock", "",
			},
			&rmValid,
			false,
		},
		{
			"Valid result with prefix of /",
			args{
				auth.NewAuthFileProvider(test.Test_conf_azure),
				"", "/mock", "",
			},
			&rmValid,
			false,
		},
		{
			"Valid result with action",
			args{
				auth.NewAuthFileProvider(test.Test_conf_azure),
				"", "mock", "mock_action",
			},
			&rmValid,
			false,
		},
		{
			"Empty endpoint",
			args{
				auth.NewAuthFileProvider(test.Test_conf_azure),
				"", "", "",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CallAzureWithEndpoint(tt.args.authProvider, tt.args.version, tt.args.endpoint, tt.args.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("CallAzureWithEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CallAzureWithEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
