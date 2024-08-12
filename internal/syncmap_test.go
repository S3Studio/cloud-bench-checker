// SyncMap util
package internal

import (
	"errors"
	"reflect"
	"testing"
)

func TestSyncMap_LoadOrCreate(t *testing.T) {
	mockKey := "mock"
	mockVal := 42

	type args struct {
		key      string
		fnCreate func() (any, error)
		nilVal   int
	}
	tests := []struct {
		name    string
		m       *SyncMap[int]
		args    args
		want    int
		wantErr bool
	}{
		{
			"Valid result with generics type of int",
			&SyncMap[int]{},
			args{
				mockKey,
				func() (any, error) { return mockVal, nil },
				0,
			},
			mockVal,
			false,
		},
		{
			"Valid result with generics type of int",
			&SyncMap[int]{},
			args{
				"invalid",
				func() (any, error) { return nil, errors.New("mock error") },
				0,
			},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.LoadOrCreate(tt.args.key, tt.args.fnCreate, tt.args.nilVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncMap.LoadOrCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SyncMap.LoadOrCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}
