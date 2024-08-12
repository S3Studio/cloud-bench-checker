// Connector for Tencent cloud

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
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

func setupEnvTencent() {
	envMap := map[string]string{
		"TENCENTCLOUD_SECRET_ID":  "mock_secretid",
		"TENCENTCLOUD_SECRET_KEY": "mock_secretkey",
		"TENCENTCLOUD_REGION":     "mock_region",
	}
	for k, v := range envMap {
		os.Setenv(k, v)
	}
}

func Test_createTencentCloudClient(t *testing.T) {
	setupEnvTencent()

	type args struct {
		p auth.IAuthProvider
	}
	tests := []struct {
		name string
		args args
		//want    *common.Client
		wantErr bool
	}{
		{
			"Valid result",
			args{auth.NewAuthFileProvider(test.Test_conf_env)},
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
		{
			"nil pointor of IAuthProvider",
			args{nil},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createTencentCloudClient(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("createTencentCloudClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("createTencentCloudClient() = %v, want a valid pointer", got)
			}
		})
	}
}

func Test_getTencentCloudClient(t *testing.T) {
	setupEnvTencent()

	type args struct {
		p auth.IAuthProvider
	}
	tests := []struct {
		name string
		args args
		//want    *common.Client
		wantErr bool
	}{
		{
			"Valid result",
			args{auth.NewAuthFileProvider(test.Test_conf_env)},
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
			got, err := getTencentCloudClient(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTencentCloudClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("getTencentCloudClient() = %v, want a valid pointer", got)
			}
		})
	}
}

func TestCallTencentCloud(t *testing.T) {
	setupEnvTencent()
	patchSend := gomonkey.ApplyMethodFunc(&common.Client{}, "Send",
		func(request tchttp.Request, response tchttp.Response) (err error) {
			return nil
		})
	defer patchSend.Reset()
	patchJsonUnmarshalCalledTime := 0
	patchJsonUnmarshal := gomonkey.ApplyFunc(internal.JsonUnmarshal,
		func(data []byte, v any) error {
			if patchJsonUnmarshalCalledTime == 0 {
				json.Unmarshal([]byte("{\"Response\":\"mock\"}"), v)
			} else {
				json.Unmarshal([]byte("{}"), v)
			}

			patchJsonUnmarshalCalledTime += 1
			return nil
		})
	defer patchJsonUnmarshal.Reset()
	rmValid, _ := internal.JsonMarshal("mock")

	type args struct {
		authProvider auth.IAuthProvider
		service      string
		version      string
		action       string
		extraParam   map[string]any
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
				auth.NewAuthFileProvider(test.Test_conf_env),
				"", "", "", make(map[string]any),
			},
			rmValid,
			false,
		},
		{
			"invalid response, missing key \"Response\"",
			args{
				auth.NewAuthFileProvider(test.Test_conf_env),
				"", "", "", make(map[string]any),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CallTencentCloud(tt.args.authProvider, tt.args.service, tt.args.version, tt.args.action, tt.args.extraParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("CallTencentCloud() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CallTencentCloud() = %v, want %v", got, tt.want)
			}
		})
	}
}
