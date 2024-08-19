// Data provider for apiserver
package server

import (
	"encoding/json"
	"fmt"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/server_model"
)

type listDataProvider struct {
	param []*server_model.ListorData
}

// GetRawDataByListorId: Implementation of IDataProvider.GetRawDataByListorId
// @param: listorId: Id of listor
// @return: Raw data of listor
// @return: Error
func (p *listDataProvider) GetRawDataByListorId(listorId int) ([]*json.RawMessage, error) {
	var value *server_model.ListorData

	for _, data := range p.param {
		if data != nil && data.ListorID == int64(listorId) {
			value = data
			break
		}
	}
	if value == nil {
		// No data of Listor in the cloud
		return nil, nil
	}

	var provideValue []*json.RawMessage
	if err := internal.JsonUnmarshal([]byte(value.Data), &provideValue); err != nil {
		return nil, err
	}

	return provideValue, nil
}

// GetCloudTypeByListorId: Implementation of IDataProvider.GetCloudTypeByListorId
// @param: listorId: Id of listor
// @return: Cloud type of listor
// @return: Error
func (p *listDataProvider) GetCloudTypeByListorId(listorId int) (string, error) {
	var value *server_model.ListorData

	for _, data := range p.param {
		if data != nil && data.ListorID == int64(listorId) {
			value = data
			break
		}
	}
	if value == nil {
		// No data of Listor in the cloud
		return "", nil
	}

	return string(value.CloudType), nil
}

// GetListorHashByListorId: Get hash of listor from request data
// @param: listorId: Id of listor
// @return: Hash of listor
// @return: Error
func (p *listDataProvider) GetListorHashByListorId(listorId int) (*server_model.ItemHash, error) {
	var value *server_model.ListorData

	for _, data := range p.param {
		if data != nil && data.ListorID == int64(listorId) {
			value = data
			break
		}
	}
	if value == nil {
		// No data of Listor in the cloud
		return nil, nil
	} else if value.ListorHash == nil {
		return nil, fmt.Errorf("nil hash pointer in data with id of: %d", listorId)
	}

	return value.ListorHash, nil
}
