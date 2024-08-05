// Jsonpath util

package internal

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
)

func rmHelper(v any) *json.RawMessage {
	r, _ := JsonMarshal(v)
	return r
}

func TestParseJsonPath(t *testing.T) {
	type args struct {
		input *json.RawMessage
		path  string
	}
	tests := []struct {
		name    string
		args    args
		want    *json.RawMessage
		wantErr bool
	}{
		{
			"Valid result",
			args{rmHelper([]any{}), "$"},
			rmHelper([]any{}),
			false,
		},
		{
			"Invalid jsonpath",
			args{rmHelper([]any{}), "invalid"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJsonPath(tt.args.input, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJsonPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJsonPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJsonPathStr(t *testing.T) {
	type args struct {
		input *json.RawMessage
		path  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Type string",
			args{rmHelper("mock"), "$"},
			"mock",
			false,
		},
		{
			"Type int",
			args{rmHelper(1), "$"},
			"1",
			false,
		},
		{
			"Type float",
			args{rmHelper(0.1), "$"},
			"0.1",
			false,
		},
		{
			"Type bool",
			args{rmHelper(true), "$"},
			strconv.FormatBool(true),
			false,
		},
		{
			"Type []string",
			args{rmHelper([]string{"mock1", "mock2"}), "$.[*]"},
			"[mock1,mock2]",
			false,
		},
		{
			"Type []any",
			args{rmHelper([]any{"mock", []int{1, 2}}), "$.[*]"},
			"[mock,[1,2]]",
			false,
		},
		// TODO: Deal with unordered map of go
		// {
		// 	"Type map[string]bool",
		// 	args{rmHelper(map[string]any{"mock_k1": true, "mock_k2": false}), "$"},
		// 	"{mock_k1:true,mock_k2:false}",
		// 	false,
		// },
		// TODO: Deal with unordered map of go
		// {
		// 	"Type map[string]string",
		// 	args{rmHelper(map[string]any{"mock_k1": "v1", "mock_k2": "v2"}), "$"},
		// 	"{mock_k1:v1,mock_k2:v2}",
		// 	false,
		// },
		{
			"Type nil",
			args{rmHelper(nil), "$"},
			"",
			false,
		},
		{
			"Type empty",
			args{rmHelper(map[string]string{"key": "value"}), "$[invalid]"},
			"",
			false,
		},
		{
			"failed to parse to string from data",
			args{rmHelper(map[bool]string{}), "$"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJsonPathStr(tt.args.input, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJsonPathStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseJsonPathStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJsonPathList(t *testing.T) {
	type args struct {
		input      *json.RawMessage
		path       string
		bObjAsList bool
	}
	tests := []struct {
		name    string
		args    args
		want    []*json.RawMessage
		wantErr bool
	}{
		{
			"Valid result",
			args{rmHelper([]string{"mock"}), "$", false},
			[]*json.RawMessage{rmHelper("mock")},
			false,
		},
		{
			"Valid result with bObjAsList==true",
			args{rmHelper(map[string]string{}), "$", true},
			[]*json.RawMessage{rmHelper(map[string]string{})},
			false,
		},
		{
			"Invalid type",
			args{rmHelper(map[string]string{}), "$", false},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJsonPathList(tt.args.input, tt.args.path, tt.args.bObjAsList)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJsonPathList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJsonPathList() = %v, want %v", got, tt.want)
			}
		})
	}
}
