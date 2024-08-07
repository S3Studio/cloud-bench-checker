// Connector for Aliyun

package connector

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	"github.com/s3studio/cloud-bench-checker/test"

	"github.com/agiledragon/gomonkey/v2"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

func setupEnvAliyun() {
	envMap := map[string]string{
		"ALIBABA_CLOUD_ACCESS_KEY_ID":     "mock_secretid",
		"ALIBABA_CLOUD_ACCESS_KEY_SECRET": "mock_secretkey",
		"ALIBABA_CLOUD_REGION":            "mock_region",
	}
	for k, v := range envMap {
		os.Setenv(k, v)
	}
}

func Test_createAliyunCloudClient(t *testing.T) {
	setupEnvAliyun()

	type args struct {
		p             auth.IAuthProvider
		endpoint      string
		bEpWithRegion bool
	}
	tests := []struct {
		name string
		args args
		//want    *openapi.Client
		wantErr bool
	}{
		{
			"Valid result with bEpWithRegion==false",
			args{auth.NewAuthFileProvider(test.Test_conf_aliyun), "", false},
			false,
		},
		{
			"Valid result with bEpWithRegion==true",
			args{auth.NewAuthFileProvider(test.Test_conf_aliyun), "", true},
			false,
		},
		{
			"Profile not defined",
			args{auth.NewAuthFileProvider(test.Test_conf_invalid), "", false},
			true,
		},
		{
			"Profile not defined",
			args{&test.MockKeyNotSetAuthProvider{}, "", false},
			true,
		},
		{
			"nil pointor of IAuthProvider",
			args{nil, "", false},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createAliyunCloudClient(tt.args.p, tt.args.endpoint, tt.args.bEpWithRegion)
			if (err != nil) != tt.wantErr {
				t.Errorf("createAliyunCloudClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("createAliyunCloudClient() = %v, want a valid pointer", got)
			}
		})
	}
}

func Test_getAliyunCloudClient(t *testing.T) {
	setupEnvAliyun()

	type args struct {
		p             auth.IAuthProvider
		endpoint      string
		bEpWithRegion bool
	}
	tests := []struct {
		name string
		args args
		//want    *openapi.Client
		wantErr bool
	}{
		{
			"Valid result",
			args{auth.NewAuthFileProvider(test.Test_conf_aliyun), "", false},
			false,
		},
		{
			"Profile not defined",
			args{auth.NewAuthFileProvider(test.Test_conf_invalid), "", false},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAliyunCloudClient(tt.args.p, tt.args.endpoint, tt.args.bEpWithRegion)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAliyunCloudClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("getAliyunCloudClient() = %v, want a valid pointer", got)
			}
		})
	}
}

func TestCallAliyunCloud(t *testing.T) {
	setupEnvAliyun()
	actionValid := "Valid"
	valValid := "mock"
	rmValid, _ := internal.JsonMarshal(valValid)
	patches := gomonkey.ApplyMethodFunc(&openapi.Client{}, "CallApi",
		func(params *openapi.Params, request *openapi.OpenApiRequest, runtime *util.RuntimeOptions) (_result map[string]interface{}, _err error) {
			if reflect.DeepEqual(params.Action, tea.String(actionValid)) {
				return map[string]any{"body": valValid}, nil
			} else {
				return make(map[string]any), nil
			}
		})
	defer patches.Reset()

	type args struct {
		authProvider  auth.IAuthProvider
		endpoint      string
		bEpWithRegion bool
		version       string
		action        string
		extraParam    map[string]any
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
				auth.NewAuthFileProvider(test.Test_conf_aliyun),
				"", false, "", actionValid,
				make(map[string]any),
			},
			rmValid,
			false,
		},
		{
			"Valid extraParam",
			args{
				auth.NewAuthFileProvider(test.Test_conf_aliyun),
				"", false, "", actionValid,
				map[string]any{"mock1": 1, "mock2": "mock_val"},
			},
			rmValid,
			false,
		},
		{
			"invalid response, missing key \"body\"",
			args{
				auth.NewAuthFileProvider(test.Test_conf_aliyun),
				"", false, "", "invalid",
				make(map[string]any),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CallAliyunCloud(tt.args.authProvider, tt.args.endpoint, tt.args.bEpWithRegion, tt.args.version, tt.args.action, tt.args.extraParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("CallAliyunCloud() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CallAliyunCloud() = %v, want %v", got, tt.want)
			}
		})
	}
}
