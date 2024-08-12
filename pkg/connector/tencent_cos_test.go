// Connector for Tencent COS

package connector

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	"github.com/s3studio/cloud-bench-checker/test"

	"github.com/agiledragon/gomonkey/v2"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

func Test_createTencentCOSClient(t *testing.T) {
	setupEnvTencent()

	type args struct {
		p          auth.IAuthProvider
		bucketName string
	}
	tests := []struct {
		name string
		args args
		//want    *cos.Client
		wantErr bool
	}{
		{
			"Valid result",
			args{auth.NewAuthFileProvider(test.Test_conf_env), "mock"},
			false,
		},
		{
			"Valid result with empty bucket name",
			args{auth.NewAuthFileProvider(test.Test_conf_env), ""},
			false,
		},
		{
			"Profile not defined",
			args{auth.NewAuthFileProvider(test.Test_conf_invalid), ""},
			true,
		},
		{
			"Key not set",
			args{&test.MockKeyNotSetAuthProvider{}, ""},
			true,
		},
		{
			"nil pointor of IAuthProvider",
			args{nil, ""},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createTencentCOSClient(tt.args.p, tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("createTencentCOSClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("createTencentCOSClient() = %v, want a valid pointer", got)
			}
		})
	}
}

func Test_getTencentCOSClient(t *testing.T) {
	setupEnvTencent()

	type args struct {
		authProvider auth.IAuthProvider
		bucketName   string
	}
	tests := []struct {
		name string
		args args
		//want    *cos.Client
		wantErr bool
	}{
		{
			"Valid result",
			args{auth.NewAuthFileProvider(test.Test_conf_env), ""},
			false,
		},
		{
			"Profile not defined",
			args{auth.NewAuthFileProvider(test.Test_conf_invalid), ""},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTencentCOSClient(tt.args.authProvider, tt.args.bucketName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTencentCOSClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("getTencentCOSClient() = %v, want a valid pointer", got)
			}
		})
	}
}

func TestCallTencentCOS(t *testing.T) {
	setupEnvTencent()
	resServiceGetResult := cos.ServiceGetResult{}
	rmServiceGetResult, _ := internal.JsonMarshal(resServiceGetResult)
	patchServiceGet := gomonkey.ApplyMethodFunc(&cos.ServiceService{}, "Get",
		func(ctx context.Context, opt ...*cos.ServiceGetOptions) (*cos.ServiceGetResult, *cos.Response, error) {
			return &resServiceGetResult, nil, nil
		})
	defer patchServiceGet.Reset()

	type args struct {
		authProvider auth.IAuthProvider
		bucketName   string
		service      string
		action       string
	}
	tests := []struct {
		name    string
		args    args
		want    *json.RawMessage
		wantErr bool
	}{
		{
			"Valid result of Service.Get",
			args{
				auth.NewAuthFileProvider(test.Test_conf_env), "", "", "",
			},
			rmServiceGetResult,
			false,
		},
		{
			"service field not found on reflect of Tencent COS client",
			args{
				auth.NewAuthFileProvider(test.Test_conf_env), "mock", "mock_service", "",
			},
			nil,
			true,
		},
		{
			"action method not found on reflect of Tencent COS client",
			args{
				auth.NewAuthFileProvider(test.Test_conf_env), "mock", "Service", "mock_action",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CallTencentCOS(tt.args.authProvider, tt.args.bucketName, tt.args.service, tt.args.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("CallTencentCOS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CallTencentCOS() = %v, want %v", got, tt.want)
			}
		})
	}
}
