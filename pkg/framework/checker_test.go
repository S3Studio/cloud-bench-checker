// Checker used to get prop from raw data and validate it

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

func TestNewChecker(t *testing.T) {
	validConf := def.ConfChecker{}

	type args struct {
		conf         *def.ConfChecker
		authProvider auth.IAuthProvider
		dataProvider IDataProvider
	}
	tests := []struct {
		name string
		args args
		want *Checker
	}{
		{
			"Valid result",
			args{&validConf, nil, nil},
			&Checker{conf: &validConf},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChecker(tt.args.conf, tt.args.authProvider, tt.args.dataProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChecker() = %v, want %v", got, tt.want)
			}
		})
	}
}

const VALID_CT = "cloud"

func setupCheckerData() *Checker {
	rm, _ := internal.JsonMarshal("mock")
	p := SyncMapDataProvider{}
	p.DataMap.Store(1, []*json.RawMessage{rm})
	p.CtMap.Store(1, VALID_CT)

	c := NewChecker(&def.ConfChecker{
		CloudType: VALID_CT,
		Listor:    []int{1},
		ExtractCmd: def.ConfExtractCmd{
			IdJsonPath:      "$",
			ExtractJsonPath: def.ConfJsonPathCmd{Path: "$"},
		},
	}, nil, &p)

	return c
}

func TestChecker_SetAuthProvider(t *testing.T) {
	checker := setupCheckerData()

	type args struct {
		authProvider auth.IAuthProvider
	}
	tests := []struct {
		name string
		c    *Checker
		args args
	}{
		{
			"Valid result",
			checker,
			args{nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetAuthProvider(tt.args.authProvider)
		})
	}
}

func TestChecker_SetDataProvider(t *testing.T) {
	checker := setupCheckerData()

	type args struct {
		dataProvider IDataProvider
	}
	tests := []struct {
		name string
		c    *Checker
		args args
	}{
		{
			"Valid result",
			checker,
			args{nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetDataProvider(tt.args.dataProvider)
		})
	}
}

func TestChecker_GetProp(t *testing.T) {
	checker := setupCheckerData()
	checkerEmpty := setupCheckerData()
	checkerEmpty.SetDataProvider(&SyncMapDataProvider{})
	checkerNil := setupCheckerData()
	checkerNil.SetDataProvider(nil)
	dpInvalidCt := SyncMapDataProvider{}
	dpInvalidCt.CtMap.Store(1, false)
	checkerInvalidCt := setupCheckerData()
	checkerInvalidCt.SetDataProvider(&dpInvalidCt)
	dpCtMismatch := SyncMapDataProvider{}
	dpCtMismatch.CtMap.Store(1, "invalid")
	checkerCtMismatch := setupCheckerData()
	checkerCtMismatch.SetDataProvider(&dpCtMismatch)
	dpInvalid := SyncMapDataProvider{}
	dpInvalid.CtMap.Store(1, VALID_CT)
	dpInvalid.DataMap.Store(1, "invalid")
	checkerInvalid := setupCheckerData()
	checkerInvalid.SetDataProvider(&dpInvalid)
	rm, _ := internal.JsonMarshal("mock")
	mockDp := SyncMapDataProvider{}
	mockDp.DataMap.Store(1, []*json.RawMessage{rm})
	mockDp.CtMap.Store(1, VALID_CT)

	type args struct {
		opts []GetPropOption
	}
	tests := []struct {
		name    string
		c       *Checker
		args    args
		want    CheckerPropList
		wantErr bool
	}{
		{
			"Valid result",
			checker,
			args{nil},
			CheckerPropList{
				{Id: "mock", Prop: rm},
			},
			false,
		},
		{
			"Valid result with empty IDataProvider",
			checkerEmpty,
			args{nil},
			nil,
			false,
		},
		{
			"Valid result with IAuthProvider in opts",
			checker,
			args{
				[]GetPropOption{SetAuthProviderOpt(&auth.AuthFileProvider{})},
			},
			CheckerPropList{
				{Id: "mock", Prop: rm},
			},
			false,
		},
		{
			"Valid result with IDataProvider in opts",
			checker,
			args{
				[]GetPropOption{SetDataProviderOpt(&mockDp)},
			},
			CheckerPropList{
				{Id: "mock", Prop: rm},
			},
			false,
		},
		{
			"failed to get raw data, provider is nil",
			checkerNil,
			args{nil},
			nil,
			true,
		},
		{
			"failed to get cloud type from provider",
			checkerInvalidCt,
			args{nil},
			nil,
			true,
		},
		{
			"cloud type of data \"%s\" mismatch cloud type of Checker",
			checkerCtMismatch,
			args{nil},
			nil,
			true,
		},
		{
			"failed to get raw data from provider",
			checkerInvalid,
			args{nil},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetProp(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Checker.GetProp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Checker.GetProp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPropWithCmd(t *testing.T) {
	rm, _ := internal.JsonMarshal("mock")

	type args struct {
		authProvider auth.IAuthProvider
		previousData CheckerProp
		conf         *def.ConfExtractCmd
		cloudType    def.CloudType
	}
	tests := []struct {
		name    string
		args    args
		want    *CheckerProp
		wantErr bool
	}{
		{
			"Valid result",
			args{
				nil,
				CheckerProp{Id: "mock_id", Prop: rm},
				&def.ConfExtractCmd{
					ExtractJsonPath: def.ConfJsonPathCmd{Path: "$"},
				},
				"",
			},
			&CheckerProp{Id: "mock_id", Prop: rm},
			false,
		},
		{
			"Valid result with IdConst",
			args{
				nil,
				CheckerProp{Id: "mock_id", Prop: rm},
				&def.ConfExtractCmd{
					IdConst:         "mock_const",
					ExtractJsonPath: def.ConfJsonPathCmd{Path: "$"},
				},
				"",
			},
			&CheckerProp{Id: "mock_const", Prop: rm},
			false,
		},
		{
			"Valid result with IdJsonPath",
			args{
				nil,
				CheckerProp{Id: "mock_id", Prop: rm},
				&def.ConfExtractCmd{
					IdJsonPath:      "$",
					ExtractJsonPath: def.ConfJsonPathCmd{Path: "$"},
				},
				"",
			},
			&CheckerProp{Id: "mock", Prop: rm},
			false,
		},
		{
			"Valid result with NameJsonPath",
			args{
				nil,
				CheckerProp{Id: "mock_id", Prop: rm},
				&def.ConfExtractCmd{
					NameJsonPath:    "$",
					ExtractJsonPath: def.ConfJsonPathCmd{Path: "$"},
				},
				"",
			},
			&CheckerProp{Id: "mock_id", Name: "mock", Prop: rm},
			false,
		},
		{
			"Valid result with NormalizeId",
			args{
				nil,
				CheckerProp{Id: "mock_group/mock_id", Prop: rm},
				&def.ConfExtractCmd{
					NormalizeId:     true,
					ExtractJsonPath: def.ConfJsonPathCmd{Path: "$"},
				},
				"",
			},
			&CheckerProp{Id: "mock_id", Prop: rm},
			false,
		},
		{
			"failed to get id",
			args{
				nil,
				CheckerProp{Id: "mock_id", Prop: rm},
				&def.ConfExtractCmd{
					IdJsonPath:      "invalid",
					ExtractJsonPath: def.ConfJsonPathCmd{Path: "$"},
				},
				"",
			},
			nil,
			true,
		},
		{
			"failed to get prop, id is empty",
			args{
				nil,
				CheckerProp{Prop: rm},
				&def.ConfExtractCmd{
					ExtractJsonPath: def.ConfJsonPathCmd{Path: "$"},
				},
				"",
			},
			nil,
			true,
		},
		{
			"failed to get name",
			args{
				nil,
				CheckerProp{Id: "mock_id", Prop: rm},
				&def.ConfExtractCmd{
					NameJsonPath:    "invalid",
					ExtractJsonPath: def.ConfJsonPathCmd{Path: "$"},
				},
				"",
			},
			nil,
			true,
		},
		{
			"failed to parse JsonPath",
			args{
				nil,
				CheckerProp{Id: "mock_id", Prop: rm},
				&def.ConfExtractCmd{
					ExtractJsonPath: def.ConfJsonPathCmd{Path: "invalid"},
				},
				"",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPropWithCmd(tt.args.authProvider, tt.args.previousData, tt.args.conf, tt.args.cloudType)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPropWithCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPropWithCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPropWithCloud(t *testing.T) {
	rm, _ := internal.JsonMarshal("mock")
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
	patchCallAzure := gomonkey.ApplyFunc(connector.CallAzureWithEndpoint,
		func(authProvider auth.IAuthProvider, version string, endpoint string, action string) (*json.RawMessage, error) {
			return rm, nil
		})
	defer patchCallAzure.Reset()
	mockAuthProvider := auth.NewAuthFileProvider(def.ConfProfile{})

	type args struct {
		authProvider auth.IAuthProvider
		cloudType    def.CloudType
		id           string
		conf         *def.ConfExtractCmd
	}
	tests := []struct {
		name    string
		args    args
		want    *json.RawMessage
		wantErr bool
	}{
		{
			"Valid result of TencentCloud",
			args{
				mockAuthProvider, def.TENCENT_CLOUD, "mock",
				&def.ConfExtractCmd{IdParamName: "mock_name", IdParamType: def.PARAM_STRING},
			},
			rm,
			false,
		},
		{
			"Valid result of TencentCOS",
			args{
				mockAuthProvider, def.TENCENT_COS, "mock", &def.ConfExtractCmd{},
			},
			rm,
			false,
		},
		{
			"Valid result of AliyunCloud",
			args{
				mockAuthProvider, def.ALIYUN_CLOUD, "mock",
				&def.ConfExtractCmd{IdParamName: "mock_name", IdParamType: def.PARAM_STRING},
			},
			rm,
			false,
		},
		{
			"Valid result of AliyunOSS",
			args{
				mockAuthProvider, def.ALIYUN_OSS, "mock", &def.ConfExtractCmd{},
			},
			rm,
			false,
		},
		{
			"Valid result of Azure",
			args{
				mockAuthProvider, def.AZURE, "mock", &def.ConfExtractCmd{},
			},
			rm,
			false,
		},
		{
			"missing IdParamName for getting prop from tencent cloud",
			args{
				mockAuthProvider, def.TENCENT_CLOUD, "mock", &def.ConfExtractCmd{},
			},
			nil,
			true,
		},
		{
			"missing IdParamName for getting prop from aliyun",
			args{
				mockAuthProvider, def.ALIYUN_CLOUD, "mock", &def.ConfExtractCmd{},
			},
			nil,
			true,
		},
		{
			"invalid cloud type",
			args{
				mockAuthProvider, "invalid", "mock", &def.ConfExtractCmd{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPropWithCloud(tt.args.authProvider, tt.args.cloudType, tt.args.id, tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPropWithCloud() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPropWithCloud() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChecker_Validate(t *testing.T) {
	rm, _ := internal.JsonMarshal("mock")

	type args struct {
		data CheckerPropList
	}
	tests := []struct {
		name    string
		c       *Checker
		v       def.ConfValidator
		args    args
		want    []*ValidateResult
		wantErr bool
	}{
		{
			"Valid result",
			NewChecker(&def.ConfChecker{}, nil, nil),
			def.ConfValidator{
				ValidateSchema: `{"type": "string"}`,
			},
			args{CheckerPropList{
				{Id: "mock_id", Prop: rm}},
			},
			[]*ValidateResult{
				{Id: "mock_id", InRisk: true},
			},
			false,
		},
		{
			"Valid result with Value",
			NewChecker(&def.ConfChecker{}, nil, nil),
			def.ConfValidator{
				ValidateSchema: `{"type": "string"}`,
				ValueJsonPath:  "$",
			},
			args{CheckerPropList{
				{Id: "mock_id", Prop: rm}},
			},
			[]*ValidateResult{
				{Id: "mock_id", InRisk: true, Value: "mock"},
			},
			false,
		},
		{
			"Failed to validate prop",
			NewChecker(&def.ConfChecker{}, nil, nil),
			def.ConfValidator{
				ValidateSchema: `{"type": "string"}`,
			},
			args{CheckerPropList{
				{Id: "mock_id", Prop: &json.RawMessage{}}},
			},
			[]*ValidateResult{},
			false,
		},
		{
			"Failed to get actual value with jsonPath",
			NewChecker(&def.ConfChecker{}, nil, nil),
			def.ConfValidator{
				ValidateSchema: `{"type": "string"}`,
				ValueJsonPath:  "invalid",
			},
			args{CheckerPropList{
				{Id: "mock_id", Prop: rm}},
			},
			[]*ValidateResult{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.conf.Validator = tt.v
			got, err := tt.c.Validate(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Checker.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Checker.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChecker_createValidator(t *testing.T) {
	tests := []struct {
		name    string
		c       *Checker
		v       def.ConfValidator
		wantErr bool
	}{
		{
			"Valid result",
			NewChecker(&def.ConfChecker{}, nil, nil),
			def.ConfValidator{
				ValidateSchema:   `{"enum": ["%%v%%"]}`,
				DynValidateValue: map[string]string{"v": "value"},
			},
			false,
		},
		{
			"failed to create jsonschema",
			NewChecker(&def.ConfChecker{}, nil, nil),
			def.ConfValidator{
				ValidateSchema: "",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.conf.Validator = tt.v
			if err := tt.c.createValidator(); (err != nil) != tt.wantErr {
				t.Errorf("Checker.createValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
