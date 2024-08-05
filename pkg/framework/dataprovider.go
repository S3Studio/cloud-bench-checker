// interface IDataProvider
package framework

import (
	"encoding/json"
	"errors"
	"sync"
)

// IDataProvider: Interface that provides different management of Listor
type IDataProvider interface {
	// GetRawDataByListorId: Get raw data of given id of Listor.
	//
	// Returns (nil, nil) if there is no data of Listor in the cloud
	//
	// IMPORTANT: The function must be goroutine safe
	// @param: listorId: Id of Listor
	// @return: Raw data of Listor
	// @return: Error
	GetRawDataByListorId(listorId int) ([]*json.RawMessage, error)
	// GetCloudTypeByListorId: Get raw data of given id of Listor.
	//
	// Returns ("", nil) if there is no data of Listor in the cloud
	//
	// IMPORTANT: The function must be goroutine safe
	// @param: listorId: Id of Listor
	// @return: Cloud type of Listor
	// @return: Error
	GetCloudTypeByListorId(listorId int) (string, error)
}

// SyncMapDataProvider: Simple implementation of IDataProvider using sync.Map
type SyncMapDataProvider struct {
	// sync.Map of data
	DataMap sync.Map
	// sync.Map of cloud_type
	CtMap sync.Map
}

// GetRawDataByListorId: Implementation of IDataProvider.GetRawDataByListorId
// @param: listorId: Id of listor
// @return: Raw data of listor
// @return: Error
func (p *SyncMapDataProvider) GetRawDataByListorId(listorId int) ([]*json.RawMessage, error) {
	value, ok := p.DataMap.Load(listorId)
	if !ok {
		// No data of Listor in the cloud
		return nil, nil
	}

	provideValue, ok := value.([]*json.RawMessage)
	if !ok {
		return nil, errors.New("data in the sync.Map is not a type of \"[]*json.RawMessage\"")
	}

	return provideValue, nil
}

// GetCloudTypeByListorId: Implementation of IDataProvider.GetCloudTypeByListorId
// @param: listorId: Id of listor
// @return: Cloud type of listor
// @return: Error
func (p *SyncMapDataProvider) GetCloudTypeByListorId(listorId int) (string, error) {
	value, ok := p.CtMap.Load(listorId)
	if !ok {
		// No data of Listor in the cloud
		return "", nil
	}

	cloudType, ok := value.(string)
	if !ok {
		return "", errors.New("data in the sync.Map is not a type of string")
	}

	return cloudType, nil
}
