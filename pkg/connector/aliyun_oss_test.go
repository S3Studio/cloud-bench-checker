// Connector for Aliyun OSS

package connector

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	"github.com/s3studio/cloud-bench-checker/test"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func Test_createAliyunOSSClient(t *testing.T) {
	setupEnvAliyun()

	type args struct {
		p auth.IAuthProvider
	}
	tests := []struct {
		name string
		args args
		//want    *oss.Client
		wantErr bool
	}{
		{
			"Valid result",
			args{auth.NewAuthFileProvider(test.Test_conf_aliyun)},
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
			got, err := createAliyunOSSClient(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("createAliyunOSSClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("createAliyunOSSClient() = %v, want a valid pointer", got)
			}
		})
	}
}

func Test_getAliyunOSSClient(t *testing.T) {
	setupEnvAliyun()

	type args struct {
		p auth.IAuthProvider
	}
	tests := []struct {
		name string
		args args
		//want    *oss.Client
		wantErr bool
	}{
		{
			"Valid result",
			args{auth.NewAuthFileProvider(test.Test_conf_aliyun)},
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
			got, err := getAliyunOSSClient(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAliyunOSSClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("getAliyunOSSClient() = %v, want a valid pointer", got)
			}
		})
	}
}

func TestCallAliyunOSS(t *testing.T) {
	setupEnvAliyun()
	resListBucket := oss.ListBucketsResult{}
	rmListBucket, _ := internal.JsonMarshal(resListBucket)
	patchListBucket := gomonkey.ApplyMethodFunc(oss.Client{}, "ListBuckets",
		func(options ...oss.Option) (oss.ListBucketsResult, error) {
			return resListBucket, nil
		})
	defer patchListBucket.Reset()
	rmFalse, _ := internal.JsonMarshal(false)
	patchIsBucketExist := gomonkey.ApplyMethodFunc(oss.Client{}, "IsBucketExist",
		func(bucketName string) (bool, error) {
			return false, nil
		})
	defer patchIsBucketExist.Reset()

	type args struct {
		authProvider auth.IAuthProvider
		bucketName   string
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
			"Valid result of Client.ListBuckets",
			args{
				auth.NewAuthFileProvider(test.Test_conf_aliyun),
				"", "", map[string]any{"marker": ""},
			},
			rmListBucket,
			false,
		},
		{
			"Valid result of other method",
			args{
				auth.NewAuthFileProvider(test.Test_conf_aliyun),
				"mock-bucket", "IsBucketExist", make(map[string]any),
			},
			rmFalse,
			false,
		},
		{
			"action method not found on reflect of Aliyun OSS client",
			args{
				auth.NewAuthFileProvider(test.Test_conf_aliyun),
				"mock-bucket", "mock-method", make(map[string]any),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CallAliyunOSS(tt.args.authProvider, tt.args.bucketName, tt.args.action, tt.args.extraParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("CallAliyunOSS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				sgot, _ := got.MarshalJSON()
				swant, _ := tt.want.MarshalJSON()
				t.Errorf("CallAliyunOSS() = %v, want %v", string(sgot), string(swant))
			}
		})
	}
}
