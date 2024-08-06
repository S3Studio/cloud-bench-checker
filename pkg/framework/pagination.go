// Merge data list from multiple pages
package framework

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"
)

// IPaginator: Interface to get single page of data
type IPaginator interface {
	// See function of GetEntireList for details of paginationParam
	// @param: paginationParam: Parameter of each page
	// @return: List of data on one page
	// @return: NextCondition, See function GetEntireList for detail
	// @return: Error
	GetOnePage(paginationParam map[string]any) ([]*json.RawMessage, NextCondition, error)
}

const DEFAULT_PAGE_SIZE = 10

// NextCondition: Indicate if data on the next page should be retrieved
//
// See function of GetEntireList for detail
type NextCondition struct {
	TotalCount int
	NextMarker string
}

// GetEntireList: Get list of all raw data according to definition of ConfPaginator
//
// There are several ways of pagination:
// Note: [i, j) means starts with index i (inclusive) and ends with index j (exclusive)
//
// - PaginationType == PAGE_OFFSET_LIMIT:
// List items of [offset, offset + limit), and offset starts with 0.
// We defines pageIndex as offset and pageSize as limit in paginationParam of GetOnePage.
// NextCondition from IPaginator.GetOnePage: Total count of items in the entire list returned by cloud
// (negative means not given)
//
// - PaginationType == PAGE_CURPAGE_SIZE:
// List items of [(curpage - 1) * pagesize, curpage * pagesize), and curpage starts with 1.
// We defines pageIndex as curpage and pageSize as pagesize in paginationParam of GetOnePage.
// NextCondition from IPaginator.GetOnePage: Total count of items in the entire list returned by cloud
// (negative means not given)
//
// - PaginationType == PAGE_MARKER:
// List items with marker of empty string on 1st page,
// and use NextMarkerName as marker for the next page if value of NextMarkerName is not empty.
// We defines marker and pagesize in paginationParam of GetOnePage.
// NextCondition from IPaginator.GetOnePage: Value of next marker
//
// @param: p: Implementation of interface IPaginator to get data of one page
// @param: conf: Definition of ConfPaginator
// @return: List of data merged from all pages
// @return: Error
func GetEntireList(p IPaginator, conf def.ConfPaginator) ([]*json.RawMessage, error) {
	if conf.PaginationType == def.PAGEINATION_DEFAULT {
		// Default value must be set to a valid value before this function is called
		return nil, errors.New("PaginationType not set")
	}

	paginationParam := make(map[string]any)

	if conf.PaginationType == def.PAGE_NOPAGEINATION {
		data, _, err := p.GetOnePage(paginationParam)
		if err != nil {
			pndError := auth.ProfileNotDefinedError{}
			if errors.As(err, &pndError) {
				// It's ok to bypass here
				return nil, nil
			}
		}

		return data, err
	}

	offset := 0
	limit := _opt.PageSize
	if limit <= 0 {
		limit = DEFAULT_PAGE_SIZE
	}
	marker := ""

	var fullList []*json.RawMessage
GetPageLoop:
	for {
		var pageIndex int
		switch conf.PaginationType {
		case def.PAGE_OFFSET_LIMIT:
			pageIndex = offset
		case def.PAGE_CURPAGE_SIZE:
			pageIndex = offset/limit + 1
		}
		pageSize := limit

		if len(conf.OffsetName) > 0 {
			if err := internal.AddParamInt(
				paginationParam,
				conf.OffsetName,
				pageIndex,
				conf.OffsetType); err != nil {
				return nil, err
			}
		}
		if len(conf.LimitName) > 0 {
			if err := internal.AddParamInt(
				paginationParam,
				conf.LimitName,
				pageSize,
				conf.LimitType); err != nil {
				return nil, err
			}
		}
		if len(conf.MarkerName) > 0 {
			if err := internal.AddParamString(
				paginationParam,
				conf.MarkerName,
				marker,
				def.PARAM_STRING); err != nil {
				return nil, err
			}
		}

		pageList, nextCondition, err := p.GetOnePage(paginationParam)
		if err != nil {
			pndError := auth.ProfileNotDefinedError{}
			if errors.As(err, &pndError) {
				// It's ok to bypass here
				return nil, nil
			}

			return nil, fmt.Errorf("failed to get data of page (offset %d/limit %d): %w", offset, limit, err)
		}

		if len(pageList) > 0 {
			fullList = append(fullList, pageList...)
		}
		switch conf.PaginationType {
		case def.PAGE_OFFSET_LIMIT, def.PAGE_CURPAGE_SIZE:
			if len(pageList) < limit ||
				(nextCondition.TotalCount >= 0 && len(fullList) >= nextCondition.TotalCount) {
				// All items are now on the list. It is possible that number of items
				// may have changed between the time function was called,
				// so it is ok if the number of items differs from the returned totalCount
				break GetPageLoop
			}
		case def.PAGE_MARKER:
			if len(nextCondition.NextMarker) == 0 {
				break GetPageLoop
			} else {
				marker = nextCondition.NextMarker
			}
		default:
			break GetPageLoop // Avoid infinite loop
		}

		offset += limit
	}

	return fullList, nil
}

type rdpOpt struct {
	// Indicate whether to put an object got by dataListJsonPath into a list and return it
	convertObjectToList *bool
}

// RDPOption: Functional options used in ResultDataParse in case more options are added
type RDPOption func(opt *rdpOpt) error

// SetConvertObjectToList: Set rdpOpt.convertObjectToList
//
// Indicate whether to put an object got by dataListJsonPath into a list and return it
// @param: flag: Value for convertObjectToList
func SetConvertObjectToList(flag bool) RDPOption {
	return func(options *rdpOpt) error {
		options.convertObjectToList = &flag
		return nil
	}
}

// ResultDataParse: Helper function to parse the pagination data from the result of cloud according to ConfPaginator
// @param: resultData: Result of the cloud
// @param: conf: Definition of ConfPaginator
// @param: dataListJsonPath: JsonPath of how to get the list from resultData
// @param: opts: Additional options
// @return: List of data on one page
// @return: NextCondition, See function GetEntireList for detail
// @return: Error
func ResultDataParse(resultData *json.RawMessage, conf def.ConfPaginator, dataListJsonPath string, opts ...RDPOption) (
	[]*json.RawMessage, NextCondition, error) {

	var optAll rdpOpt
	for _, opt := range opts {
		err := opt(&optAll)
		if err != nil {
			return nil, NextCondition{}, err
		}
	}
	if optAll.convertObjectToList == nil {
		defaultValue := false
		optAll.convertObjectToList = &defaultValue
	}

	switch conf.PaginationType {
	case def.PAGE_OFFSET_LIMIT, def.PAGE_CURPAGE_SIZE:
		resultMap := make(map[string]json.RawMessage)
		if err := internal.JsonUnmarshal(*resultData, &resultMap); err != nil {
			return nil, NextCondition{}, fmt.Errorf("failed to unmarshal as map: %w", err)
		}

		rmTotalCount, ok := resultMap[conf.RespTotalName]
		if !ok {
			return nil, NextCondition{}, fmt.Errorf("invalid response, missing key \"%s\"", conf.RespTotalName)
		}
		var nTotalCount json.Number
		if err := internal.JsonUnmarshal(rmTotalCount, &nTotalCount); err != nil {
			return nil, NextCondition{}, fmt.Errorf("failed to convert to number: %w", err)
		}
		totalCount, err := nTotalCount.Int64()
		if err != nil {
			return nil, NextCondition{}, fmt.Errorf("failed to convert to int: %w", err)
		}

		dataList, err := internal.ParseJsonPathList(resultData, dataListJsonPath, *optAll.convertObjectToList)
		if err != nil {
			return nil, NextCondition{}, fmt.Errorf("failed to convert to list: %w", err)
		}

		// totalCount will be truncated from 64-bit to 32-bit, hope there are not too many resources
		return dataList, NextCondition{TotalCount: int(totalCount)}, nil
	case def.PAGE_NOPAGEINATION:
		dataList, err := internal.ParseJsonPathList(resultData, dataListJsonPath, *optAll.convertObjectToList)
		if err != nil {
			return nil, NextCondition{}, fmt.Errorf("failed to convert to list: %w", err)
		}

		return dataList, NextCondition{TotalCount: len(dataList)}, nil
	case def.PAGE_MARKER:
		resultMap := make(map[string]json.RawMessage)
		if err := internal.JsonUnmarshal(*resultData, &resultMap); err != nil {
			return nil, NextCondition{}, fmt.Errorf("failed to unmarshal as map: %w", err)
		}

		isTruncated := false
		if len(conf.TruncatedName) == 0 {
			isTruncated = true
		} else {
			rmTruncated, ok := resultMap[conf.TruncatedName]
			if !ok {
				return nil, NextCondition{}, fmt.Errorf("invalid response, missing key \"%s\"", conf.TruncatedName)
			}

			if err := internal.JsonUnmarshal(rmTruncated, &isTruncated); err != nil {
				return nil, NextCondition{}, fmt.Errorf("failed to convert to bool: %w", err)
			}
		}

		// If isTruncated == true or not set, try to get nextMarker
		nextMarker := ""
		if isTruncated {
			if len(conf.NextMarkerName) == 0 {
				return nil, NextCondition{}, errors.New("config of NextMarkerName is empty")
			}

			rmNextMarker, ok := resultMap[conf.NextMarkerName]
			if ok {
				if err := internal.JsonUnmarshal(rmNextMarker, &nextMarker); err != nil {
					return nil, NextCondition{}, fmt.Errorf("failed to convert to string: %w", err)
				}
			} else {
				nextMarker = ""
			}
		}

		dataList, err := internal.ParseJsonPathList(resultData, dataListJsonPath, *optAll.convertObjectToList)
		if err != nil {
			return nil, NextCondition{}, fmt.Errorf("failed to convert to list: %w", err)
		}

		return dataList, NextCondition{NextMarker: nextMarker}, nil
	default:
		return nil, NextCondition{}, fmt.Errorf("failed to deal with PaginationType: %v", conf.PaginationType)
	}
}
