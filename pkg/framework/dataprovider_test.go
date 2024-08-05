// IDataProvider is interface which provide different management of listor

package framework

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/s3studio/cloud-bench-checker/internal"
)

func TestSyncMapDataProvider_GetRawDataByListorId(t *testing.T) {
	rm, _ := internal.JsonMarshal("mock")
	p := SyncMapDataProvider{}
	p.DataMap.Store(1, []*json.RawMessage{rm})
	p.DataMap.Store(2, "invalid")

	type args struct {
		listorId int
	}
	tests := []struct {
		name    string
		p       *SyncMapDataProvider
		args    args
		want    []*json.RawMessage
		wantErr bool
	}{
		{
			"Valid result",
			&p,
			args{1},
			[]*json.RawMessage{rm},
			false,
		},
		{
			"No data in the cloud",
			&p,
			args{0},
			nil,
			false,
		},
		{
			"data of map is not type of \"[]*json.RawMessage\"",
			&p,
			args{2},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetRawDataByListorId(tt.args.listorId)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncMapDataProvider.GetRawDataByListorId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SyncMapDataProvider.GetRawDataByListorId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncMapDataProvider_GetCloudTypeByListorId(t *testing.T) {
	cloudType := "mock_ct"
	p := SyncMapDataProvider{}
	p.CtMap.Store(1, cloudType)
	p.CtMap.Store(2, false)

	type args struct {
		listorId int
	}
	tests := []struct {
		name    string
		p       *SyncMapDataProvider
		args    args
		want    string
		wantErr bool
	}{
		{
			"Valid result",
			&p,
			args{1},
			cloudType,
			false,
		},
		{
			"No data in the cloud",
			&p,
			args{0},
			"",
			false,
		},
		{
			"data of map is not type of string",
			&p,
			args{2},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetCloudTypeByListorId(tt.args.listorId)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncMapDataProvider.GetCloudTypeByListorId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SyncMapDataProvider.GetCloudTypeByListorId() = %v, want %v", got, tt.want)
			}
		})
	}
}
