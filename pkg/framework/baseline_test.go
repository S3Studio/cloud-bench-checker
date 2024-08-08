// Baseline of management for process of checking

package framework

import (
	"crypto"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"
)

var (
	mockValidCheckProp = CheckerProp{}
	mockValidResult    = ValidateResult{}
)

func TestNewBaseline(t *testing.T) {
	validConf := def.ConfBaseline{Checker: make([]def.ConfChecker, 1)}

	type args struct {
		conf         *def.ConfBaseline
		authProvider auth.IAuthProvider
		dataProvider IDataProvider
	}
	tests := []struct {
		name string
		args args
		want *Baseline
	}{
		{
			"Valid result",
			args{&validConf, nil, nil},
			&Baseline{
				conf: &validConf,
				checker: []*Checker{
					{conf: &validConf.Checker[0]},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBaseline(tt.args.conf, tt.args.authProvider, tt.args.dataProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBaseline() = %v, want %v", got, tt.want)
			}
		})
	}
}

var mockMetadata = map[string]string{
	"mock_key": "mock_value",
}

var mockBaseline = NewBaseline(
	&def.ConfBaseline{
		Metadata: mockMetadata,
		Checker: []def.ConfChecker{
			{CloudType: "mock", Listor: []int{1}},
			{CloudType: "invalid", Listor: []int{1}},
		},
	},
	nil, nil,
)

func TestBaseline_SetAuthProvider(t *testing.T) {
	type args struct {
		authProvider auth.IAuthProvider
	}
	tests := []struct {
		name string
		b    *Baseline
		args args
	}{
		{
			"Valid result",
			mockBaseline,
			args{nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetAuthProvider(tt.args.authProvider)
		})
	}
}

func TestBaseline_SetDataProvider(t *testing.T) {
	type args struct {
		dataProvider IDataProvider
	}
	tests := []struct {
		name string
		b    *Baseline
		args args
	}{
		{
			"Valid result",
			mockBaseline,
			args{nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetDataProvider(tt.args.dataProvider)
		})
	}
}

func TestBaseline_GetListorId(t *testing.T) {
	tests := []struct {
		name string
		b    *Baseline
		want []int
	}{
		{
			"Valid result",
			mockBaseline,
			[]int{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GetListorId(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Baseline.GetListorId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setupChecker() func() {
	patchGetProp := gomonkey.ApplyFunc((*Checker).GetProp,
		func(c *Checker) (CheckerPropList, error) {
			if c.conf.CloudType == "invalid" {
				return nil, errors.New("mock invalid Checker.GetProp")
			}

			return CheckerPropList{&mockValidCheckProp}, nil
		})
	patchValidate := gomonkey.ApplyMethodFunc(&Checker{}, "Validate",
		func(data CheckerPropList) ([]*ValidateResult, error) {
			if data == nil {
				return nil, errors.New("mock invalid Checker.Validate")
			}

			return []*ValidateResult{&mockValidResult}, nil
		})

	return func() {
		patchValidate.Reset()
		patchGetProp.Reset()
	}
}

func TestBaseline_GetProp(t *testing.T) {
	deferFn := setupChecker()
	defer deferFn()

	tests := []struct {
		name string
		b    *Baseline
		want BaselinePropList
	}{
		{
			"Valid result",
			mockBaseline,
			BaselinePropList{
				{&mockValidCheckProp},
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GetProp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Baseline.GetProp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseline_Validate(t *testing.T) {
	deferFn := setupChecker()
	defer deferFn()
	prop := mockBaseline.GetProp()

	type args struct {
		data BaselinePropList
	}
	tests := []struct {
		name    string
		b       *Baseline
		args    args
		want    []*ValidateResult
		wantErr bool
	}{
		{
			"Valid result",
			mockBaseline,
			args{prop},
			[]*ValidateResult{&mockValidResult},
			false,
		},
		{
			"size mismatch between props and checkers",
			mockBaseline,
			args{prop[:1]},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Validate(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Baseline.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Baseline.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseline_GetMetadata(t *testing.T) {
	tests := []struct {
		name string
		b    *Baseline
		want *map[string]string
	}{
		{
			"Valid result",
			mockBaseline,
			&mockMetadata,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.GetMetadata(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Baseline.GetMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseline_GetHash(t *testing.T) {
	mockHash := []byte("1")

	type args struct {
		hashType       crypto.Hash
		listorHashList [][]*[]byte
	}
	tests := []struct {
		name    string
		b       *Baseline
		args    args
		want    string
		wantErr bool
	}{
		{
			"Valid result",
			NewBaseline(&def.ConfBaseline{
				Checker: []def.ConfChecker{{Listor: []int{1}}},
			}, nil, nil),
			args{
				crypto.SHA256,
				[][]*[]byte{
					{&mockHash},
				},
			},
			"d00701fd7e5e81a594329de7bf063a60e8dc803f95a30b3afc089b1b14338589", // hardcode value
			false,
		},
		{
			"size mismatch between Checker and given hash list",
			NewBaseline(&def.ConfBaseline{
				Checker: []def.ConfChecker{{Listor: []int{1}}},
			}, nil, nil),
			args{
				crypto.SHA256,
				[][]*[]byte{},
			},
			"",
			true,
		},
		{
			"size mismatch between Checker and given hash list",
			NewBaseline(&def.ConfBaseline{
				Checker: []def.ConfChecker{{Listor: []int{1}}},
			}, nil, nil),
			args{
				crypto.SHA256,
				[][]*[]byte{
					{&mockHash, &mockHash},
				},
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.GetHash(tt.args.hashType, tt.args.listorHashList)
			if (err != nil) != tt.wantErr {
				t.Errorf("Baseline.GetHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && fmt.Sprintf("%x", got) != tt.want {
				t.Errorf("Baseline.GetHash() = %v, want %v", fmt.Sprintf("%x", got), tt.want)
			}
		})
	}
}
