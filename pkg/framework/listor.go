// Listor to get raw data list from cloud
package framework

import (
	"encoding/json"
	"fmt"

	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	"github.com/s3studio/cloud-bench-checker/pkg/connector"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"
)

// Listor: Used to retrieve a list of resources from the cloud
//
// Implements the interface of IPaginator
type Listor struct {
	conf         *def.ConfListor
	authProvider auth.IAuthProvider
}

const AZURE_NEXT_MARKER = "nextLink"

// _defaultPaginatorConf: Default paginator definition of different cloud connector
var _defaultPaginatorConf = map[def.CloudType]def.ConfPaginator{
	def.TENCENT_CLOUD: {
		PaginationType: def.PAGE_OFFSET_LIMIT,
		OffsetType:     def.PARAM_INT,
		OffsetName:     "Offset",
		LimitType:      def.PARAM_INT,
		LimitName:      "Limit",
		RespTotalName:  "TotalCount",
	},
	def.TENCENT_COS: {
		PaginationType: def.PAGE_NOPAGEINATION,
	},
	def.ALIYUN_OSS: {
		PaginationType: def.PAGE_MARKER,
		MarkerName:     connector.ALIYUN_OSS_MARKER_KEY,
		NextMarkerName: "NextMarker",
		TruncatedName:  "IsTruncated",
	},
	def.K8S: {
		PaginationType: def.PAGE_NOPAGEINATION,
	},
	def.AZURE: {
		PaginationType: def.PAGE_MARKER,
		MarkerName:     AZURE_NEXT_MARKER,
		NextMarkerName: AZURE_NEXT_MARKER,
	},
	// No default definition for def.ALIYUN_CLOUD as it varies from API to API
}

// NewListor:Constructor of Listor
// @param: conf: Definition of Listor
// @param: authProvider: IAuthProvider to provide profile of auth
func NewListor(conf *def.ConfListor, authProvider auth.IAuthProvider) *Listor {
	listor := Listor{conf: conf, authProvider: authProvider}
	if listor.conf.Paginator.PaginationType == def.PAGEINATION_DEFAULT {
		listor.conf.Paginator = _defaultPaginatorConf[listor.conf.CloudType]
	}
	return &listor
}

// SetAuthProvider: Set new authProvider
// @param: authProvider: New provider
func (l *Listor) SetAuthProvider(authProvider auth.IAuthProvider) {
	l.authProvider = authProvider
}

// ListData: Get list of all raw data according to Listor.conf
//
// Raw data from different pages are merged where necessary.
// Listor.GetOnePage is called to retrieve data as an implementation of IPaginator.
// @return: List of raw data
// @return: Error
func (l *Listor) ListData() ([]*json.RawMessage, error) {
	return GetEntireList(l, l.conf.Paginator)
}

// GetOnePage: Implementation of IPaginator.GetOnePage
//
// See function of GetEntireList in pagination for details of paginationParam
// @param: paginationParam: Parameter of each page
// @return: List of data on one page
// @return: NextCondition, See function of GetEntireList in pagination for detail
// @return: Error
func (l *Listor) GetOnePage(paginationParam map[string]any) ([]*json.RawMessage, NextCondition, error) {
	switch l.conf.CloudType {
	case def.TENCENT_CLOUD:
		mergeMaps(&paginationParam, l.conf.ListCmd.TencentCloud.ExtraParam)

		pageRes, err := connector.CallTencentCloud(
			l.authProvider,
			l.conf.ListCmd.TencentCloud.Service,
			l.conf.ListCmd.TencentCloud.Version,
			l.conf.ListCmd.TencentCloud.Action,
			paginationParam,
		)
		if err != nil {
			return nil, NextCondition{}, err
		}

		return ResultDataParse(pageRes, l.conf.Paginator, l.conf.ListCmd.DataListJsonPath,
			SetConvertObjectToList(l.conf.ListCmd.ConvertObjectToList),
		)
	case def.TENCENT_COS:
		// service and action are ignored, and bucketName is set to empty,
		// so that CallTencentCOS returns a list of all buckets
		res, err := connector.CallTencentCOS(
			l.authProvider,
			"", "", "",
		)
		if err != nil {
			return nil, NextCondition{}, err
		}

		dataListJsonPath := l.conf.ListCmd.DataListJsonPath
		if len(dataListJsonPath) == 0 {
			dataListJsonPath = "$.Buckets" // Default value
		}

		return ResultDataParse(res, l.conf.Paginator, dataListJsonPath,
			SetConvertObjectToList(l.conf.ListCmd.ConvertObjectToList),
		)
	case def.ALIYUN_CLOUD:
		mergeMaps(&paginationParam, l.conf.ListCmd.Aliyun.ExtraParam)

		pageRes, err := connector.CallAliyunCloud(
			l.authProvider,
			l.conf.ListCmd.Aliyun.Endpoint,
			l.conf.ListCmd.Aliyun.EndpointWithRegion,
			l.conf.ListCmd.Aliyun.Version,
			l.conf.ListCmd.Aliyun.Action,
			paginationParam,
		)
		if err != nil {
			return nil, NextCondition{}, err
		}

		return ResultDataParse(pageRes, l.conf.Paginator, l.conf.ListCmd.DataListJsonPath,
			SetConvertObjectToList(l.conf.ListCmd.ConvertObjectToList),
		)
	case def.ALIYUN_OSS:
		// action is ignored, and bucketName is set to empty,
		// so that CallAliyunOSS returns a list of all buckets
		pageRes, err := connector.CallAliyunOSS(
			l.authProvider,
			"", "",
			paginationParam,
		)
		if err != nil {
			return nil, NextCondition{}, err
		}

		dataListJsonPath := l.conf.ListCmd.DataListJsonPath
		if len(dataListJsonPath) == 0 {
			dataListJsonPath = "$.Buckets" // Default value
		}

		return ResultDataParse(pageRes, l.conf.Paginator, dataListJsonPath,
			SetConvertObjectToList(l.conf.ListCmd.ConvertObjectToList),
		)
	case def.K8S:
		mergeMaps(&paginationParam, l.conf.ListCmd.K8sList.ListOptions)
		res, err := connector.CallK8sList(
			l.authProvider,
			l.conf.ListCmd.K8sList.Namespace,
			l.conf.ListCmd.K8sList.Group,
			l.conf.ListCmd.K8sList.Version,
			l.conf.ListCmd.K8sList.Resource,
			paginationParam,
		)
		if err != nil {
			return nil, NextCondition{}, err
		}

		dataListJsonPath := l.conf.ListCmd.DataListJsonPath
		if len(dataListJsonPath) == 0 {
			dataListJsonPath = "$.items" // Default value
		}

		return ResultDataParse(res, l.conf.Paginator, dataListJsonPath,
			SetConvertObjectToList(l.conf.ListCmd.ConvertObjectToList),
		)
	case def.AZURE:
		// nextLink is empty on the first call of listing
		nextLink, _ := paginationParam[AZURE_NEXT_MARKER].(string)

		res, err := connector.CallAzureList(
			l.authProvider,
			l.conf.ListCmd.Azure.Provider,
			l.conf.ListCmd.Azure.Version,
			l.conf.ListCmd.Azure.RsType,
			nextLink,
		)
		if err != nil {
			return nil, NextCondition{}, err
		}

		dataListJsonPath := l.conf.ListCmd.DataListJsonPath
		if len(dataListJsonPath) == 0 {
			dataListJsonPath = "$.value" // Default value
		}

		return ResultDataParse(res, l.conf.Paginator, dataListJsonPath,
			SetConvertObjectToList(l.conf.ListCmd.ConvertObjectToList),
		)
	default:
		return nil, NextCondition{}, fmt.Errorf("invalid cloud type of %s", l.conf.CloudType)
	}
}

func mergeMaps(target *map[string]any, from ...map[string]any) {
	for _, m := range from {
		for k, v := range m {
			(*target)[k] = v
		}
	}
}
