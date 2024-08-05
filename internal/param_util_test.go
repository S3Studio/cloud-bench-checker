// Param util

package internal

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/s3studio/cloud-bench-checker/pkg/definition"
)

func TestAddParamString(t *testing.T) {
	testMap := make(map[string]any)
	mockKey := "mock_key"
	mockVal := "mock_value"

	type args struct {
		m          map[string]any
		key        string
		value      string
		targetType definition.ParamType
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		actualEqual func(v any) bool
	}{
		{
			"Add as int",
			args{testMap, mockKey, "1", definition.PARAM_INT},
			false,
			func(v any) bool { return reflect.DeepEqual(v.(int), 1) },
		},
		{
			"Add as string",
			args{testMap, mockKey, mockVal, definition.PARAM_STRING},
			false,
			func(v any) bool { return reflect.DeepEqual(v.(string), mockVal) },
		},
		{
			"Add as []string",
			args{testMap, mockKey, mockVal, definition.PARAM_STRING_LIST},
			false,
			func(v any) bool { return reflect.DeepEqual(v.([]string), []string{mockVal}) },
		},
		{
			"failed to convert string to int",
			args{testMap, mockKey, "invalid", definition.PARAM_INT},
			true,
			nil,
		},
		{
			"invalid param type",
			args{testMap, "", "", "invalid"},
			true,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddParamString(tt.args.m, tt.args.key, tt.args.value, tt.args.targetType); (err != nil) != tt.wantErr {
				t.Errorf("AddParamString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !tt.actualEqual(testMap[mockKey]) {
				t.Errorf("Value after add = %v", testMap[mockKey])
			}
		})
	}
}

func TestAddParamInt(t *testing.T) {
	testMap := make(map[string]any)
	mockKey := "mock_key"
	mockVal := 1

	type args struct {
		m          map[string]any
		key        string
		value      int
		targetType definition.ParamType
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		actualEqual func(v any) bool
	}{
		{
			"Add as int",
			args{testMap, mockKey, mockVal, definition.PARAM_INT},
			false,
			func(v any) bool { return reflect.DeepEqual(v.(int), mockVal) },
		},
		{
			"Add as string",
			args{testMap, mockKey, mockVal, definition.PARAM_STRING},
			false,
			func(v any) bool { return reflect.DeepEqual(v.(string), fmt.Sprint(mockVal)) },
		},
		{
			"Add as []string",
			args{testMap, mockKey, mockVal, definition.PARAM_STRING_LIST},
			false,
			func(v any) bool { return reflect.DeepEqual(v.([]string), []string{fmt.Sprint(mockVal)}) },
		},
		{
			"invalid param type",
			args{testMap, mockKey, mockVal, "invalid"},
			true,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddParamInt(tt.args.m, tt.args.key, tt.args.value, tt.args.targetType); (err != nil) != tt.wantErr {
				t.Errorf("AddParamInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !tt.actualEqual(testMap[mockKey]) {
				t.Errorf("Value after add = %v", testMap[mockKey])
			}
		})
	}
}
