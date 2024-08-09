// Listor to get raw data list from cloud

package framework

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	"github.com/s3studio/cloud-bench-checker/pkg/connector"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	"github.com/agiledragon/gomonkey/v2"
)

func TestNewListor(t *testing.T) {
	validConf := def.ConfListor{}

	type args struct {
		conf         *def.ConfListor
		authProvider auth.IAuthProvider
	}
	tests := []struct {
		name string
		args args
		want *Listor
	}{
		{
			"Valid result",
			args{&validConf, nil},
			&Listor{conf: &validConf},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewListor(tt.args.conf, tt.args.authProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewListor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListor_SetAuthProvider(t *testing.T) {
	listor := NewListor(&def.ConfListor{}, nil)

	type args struct {
		authProvider auth.IAuthProvider
	}
	tests := []struct {
		name string
		l    *Listor
		args args
	}{
		{
			"Valid result",
			listor,
			args{nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.SetAuthProvider(tt.args.authProvider)
		})
	}
}

func TestListor_GetOnePage(t *testing.T) {
	rm, _ := internal.JsonMarshal("mock")
	rmList := []*json.RawMessage{rm}
	patchCallTencentCloud := gomonkey.ApplyFunc(connector.CallTencentCloud,
		func(authProvider auth.IAuthProvider, service string, version string, action string, extraParam map[string]any) (*json.RawMessage, error) {
			return rm, nil
		})
	defer patchCallTencentCloud.Reset()
	patchCallTencentCOS := gomonkey.ApplyFunc(connector.CallTencentCOS,
		func(authProvider auth.IAuthProvider, bucketName string, service string, action string) (*json.RawMessage, error) {
			return rm, nil
		})
	defer patchCallTencentCOS.Reset()
	patchCallAliyunCloud := gomonkey.ApplyFunc(connector.CallAliyunCloud,
		func(authProvider auth.IAuthProvider, endpoint string, bEpWithRegion bool, version string, action string, extraParam map[string]any) (*json.RawMessage, error) {
			return rm, nil
		})
	defer patchCallAliyunCloud.Reset()
	patchCallAliyunOSS := gomonkey.ApplyFunc(connector.CallAliyunOSS,
		func(authProvider auth.IAuthProvider, bucketName string, action string, extraParam map[string]any) (*json.RawMessage, error) {
			return rm, nil
		})
	defer patchCallAliyunOSS.Reset()
	patchCallK8sList := gomonkey.ApplyFunc(connector.CallK8sList,
		func(authProvider auth.IAuthProvider, namespace string, group string, version string, resource string, extraParam map[string]any) (*json.RawMessage, error) {
			return rm, nil
		})
	defer patchCallK8sList.Reset()
	patchCallAzureList := gomonkey.ApplyFunc(connector.CallAzureList,
		func(authProvider auth.IAuthProvider, provider string, version string, rsType string, nextLink string) (*json.RawMessage, error) {
			return rm, nil
		})
	defer patchCallAzureList.Reset()
	patchRDP := gomonkey.ApplyFunc(ResultDataParse,
		func(resultData *json.RawMessage, conf def.ConfPaginator, dataListJsonPath string, opts ...RDPOption) (
			[]*json.RawMessage, NextCondition, error) {
			return rmList, NextCondition{}, nil
		})
	defer patchRDP.Reset()

	type args struct {
		paginationParam map[string]any
	}
	tests := []struct {
		name    string
		l       *Listor
		v       def.ConfListCmd
		args    args
		want    []*json.RawMessage
		want1   NextCondition
		wantErr bool
	}{
		{
			"Valid result of TencentCloud",
			NewListor(&def.ConfListor{CloudType: def.TENCENT_CLOUD}, nil),
			def.ConfListCmd{},
			args{nil},
			rmList,
			NextCondition{},
			false,
		},
		{
			"Valid result of TencentCOS",
			NewListor(&def.ConfListor{CloudType: def.TENCENT_COS}, nil),
			def.ConfListCmd{},
			args{nil},
			rmList,
			NextCondition{},
			false,
		},
		{
			"Valid result of AliyunCloud",
			NewListor(&def.ConfListor{CloudType: def.ALIYUN_CLOUD}, nil),
			def.ConfListCmd{},
			args{nil},
			rmList,
			NextCondition{},
			false,
		},
		{
			"Valid result of AliyunOSS",
			NewListor(&def.ConfListor{CloudType: def.ALIYUN_OSS}, nil),
			def.ConfListCmd{},
			args{nil},
			rmList,
			NextCondition{},
			false,
		},
		{
			"Valid result of K8s",
			NewListor(&def.ConfListor{CloudType: def.K8S}, nil),
			def.ConfListCmd{},
			args{nil},
			rmList,
			NextCondition{},
			false,
		},
		{
			"Valid result of Azure",
			NewListor(&def.ConfListor{CloudType: def.AZURE}, nil),
			def.ConfListCmd{},
			args{map[string]any{}},
			rmList,
			NextCondition{},
			false,
		},
		{
			"Valid result with mergeMaps",
			NewListor(&def.ConfListor{CloudType: def.TENCENT_CLOUD}, nil),
			def.ConfListCmd{
				TencentCloud: def.ConfTencentCloudCmd{
					ExtraParam: map[string]any{"mock_key": "mock_val"},
				},
			},
			args{make(map[string]any)},
			rmList,
			NextCondition{},
			false,
		},
		{
			"invalid cloud type",
			NewListor(&def.ConfListor{CloudType: "invalid"}, nil),
			def.ConfListCmd{},
			args{nil},
			nil,
			NextCondition{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.conf.ListCmd = tt.v
			got, got1, err := tt.l.GetOnePage(tt.args.paginationParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("Listor.GetOnePage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Listor.GetOnePage() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Listor.GetOnePage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
