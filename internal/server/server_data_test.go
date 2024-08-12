// Data provider for apiserver

package server

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/server_model"
)

func TestListDataProvider_GetRawDataByListorId(t *testing.T) {
	strRm := "mock"
	rm, _ := internal.JsonMarshal(strRm)
	p := listDataProvider{
		[]*server_model.ListorData{
			{ListorID: 1, Data: fmt.Sprintf("[\"%s\"]", strRm)},
		},
	}

	type args struct {
		listorId int
	}
	tests := []struct {
		name    string
		p       *listDataProvider
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetRawDataByListorId(tt.args.listorId)
			if (err != nil) != tt.wantErr {
				t.Errorf("listDataProvider.GetRawDataByListorId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listDataProvider.GetRawDataByListorId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListDataProvider_GetCloudTypeByListorId(t *testing.T) {
	cloudType := "mock_ct"
	p := listDataProvider{
		[]*server_model.ListorData{
			{ListorID: 1, CloudType: server_model.Cloudtype4api(cloudType)},
		},
	}

	type args struct {
		listorId int
	}
	tests := []struct {
		name    string
		p       *listDataProvider
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetCloudTypeByListorId(tt.args.listorId)
			if (err != nil) != tt.wantErr {
				t.Errorf("listDataProvider.GetCloudTypeByListorId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("listDataProvider.GetCloudTypeByListorId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listDataProvider_GetListorHashByListorId(t *testing.T) {
	hash := server_model.ItemHash{Sha256: "mock"}
	p := listDataProvider{
		[]*server_model.ListorData{
			{ListorID: 1, ListorHash: &hash},
			{ListorID: 2},
		},
	}

	type args struct {
		listorId int
	}
	tests := []struct {
		name    string
		p       *listDataProvider
		args    args
		want    *server_model.ItemHash
		wantErr bool
	}{
		{
			"Valid result",
			&p,
			args{1},
			&hash,
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
			"nil hash pointer in data",
			&p,
			args{2},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetListorHashByListorId(tt.args.listorId)
			if (err != nil) != tt.wantErr {
				t.Errorf("listDataProvider.GetListorHashByListorId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listDataProvider.GetListorHashByListorId() = %v, want %v", got, tt.want)
			}
		})
	}
}
