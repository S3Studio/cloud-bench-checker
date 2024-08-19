// Checker used to get prop from raw data and validate it

package framework

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"strings"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	"github.com/s3studio/cloud-bench-checker/pkg/connector"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	"github.com/xeipuuv/gojsonschema"
)

// Checker: Used to extract properties and validate them
//
//	Usage of Checker consists of 2 steps:
//
// 1. GetProp: Extract Id, Name (if required) and properties of the raw data
// from either IDataProvider or cloud connector
// 2. Validate: Validate properties and generate result according to the rule
// defined in JsonSchema
type Checker struct {
	// Definition of Checker
	conf *def.ConfChecker
	// Instance of validator of JsonSchema
	validator *gojsonschema.Schema
	// Flag to avoid repeating the same errors
	validatorLoadFailed bool
	// IAuthProvider to provide profile of auth
	authProvider auth.IAuthProvider
	// IDataProvider to provide raw data
	dataProvider IDataProvider
}

// NewChecker: Constructor of Checker
// @param: conf: Definition of Baseline
// @param: authProvider: IAuthProvider to provide profile of auth
// @param: dataProvider: IDataProvider to provide raw data
func NewChecker(conf *def.ConfChecker, authProvider auth.IAuthProvider, dataProvider IDataProvider) *Checker {
	return &Checker{conf: conf, authProvider: authProvider, dataProvider: dataProvider}
}

// SetAuthProvider: Set new authProvider
// @param: authProvider: New provider
func (c *Checker) SetAuthProvider(authProvider auth.IAuthProvider) {
	c.authProvider = authProvider
}

// SetDataProvider: Set new dataProvider
// @param: dataProvider: New provider
func (c *Checker) SetDataProvider(dataProvider IDataProvider) {
	c.dataProvider = dataProvider
}

// CheckerProp: Properties extracted from raw data that need to be validated
type CheckerProp struct {
	// Resource identifier used in cloud connector
	Id string
	// Human readable name of the resource
	Name string
	// Properties extracted
	Prop *json.RawMessage
}

type getPropOpt struct {
	// IAuthProvider used in call of GetProp instead of default value
	ap auth.IAuthProvider
	// IDataProvider used in call of GetProp instead of default value
	dp IDataProvider
}

// GetPropOption: Functional options used in GetProp in case more options are added
type GetPropOption func(opt *getPropOpt) error

// SetAuthProviderOpt: Set getPropOpt.ap
//
// IAuthProvider used in call of GetProp instead of default value
// @param: val: Value for IAuthProvider
func SetAuthProviderOpt(val auth.IAuthProvider) GetPropOption {
	return func(options *getPropOpt) error {
		options.ap = val
		return nil
	}
}

// SetDataProviderOpt: Set getPropOpt.dp
//
// IDataProvider used in call of GetProp instead of default value
// @param: val: Value for IDataProvider
func SetDataProviderOpt(val IDataProvider) GetPropOption {
	return func(options *getPropOpt) error {
		options.dp = val
		return nil
	}
}

// CheckerPropList: Type alias of list of CheckerProp
type CheckerPropList []*CheckerProp

// GetProp: Extract Id, Name (if required) and properties of the raw data
// @param: opts: Additional options
// @return: List of properties extracted from raw data
// @return: Error
func (c *Checker) GetProp(opts ...GetPropOption) (CheckerPropList, error) {
	var optAll getPropOpt
	for _, opt := range opts {
		err := opt(&optAll)
		if err != nil {
			return nil, err
		}
	}
	authProvider := optAll.ap
	if authProvider == nil {
		authProvider = c.authProvider
	}
	dataProvider := optAll.dp
	if dataProvider == nil {
		dataProvider = c.dataProvider
	}

	var checkerPropList CheckerPropList
	for _, listorId := range c.conf.Listor {
		var eachListorData []*json.RawMessage
		var err error
		if dataProvider == nil {
			return nil, errors.New("failed to get raw data, provider is nil")
		}
		if cloudType, err := dataProvider.GetCloudTypeByListorId(listorId); err != nil {
			return nil, fmt.Errorf("failed to get cloud type from provider: %w", err)
		} else if cloudType == "" {
			// No data of Listor in the cloud, bypass to the next Listor
			continue
		} else if cloudType != string(c.conf.CloudType) {
			return nil, fmt.Errorf("cloud type of data \"%s\" mismatch cloud type of Checker \"%s\"",
				cloudType, c.conf.CloudType)
		}
		if eachListorData, err = dataProvider.GetRawDataByListorId(listorId); err != nil {
			return nil, fmt.Errorf("failed to get raw data from provider: %w", err)
		}

		for _, rawData := range eachListorData {
			eachData := &CheckerProp{Prop: rawData}
			eachData, err = getPropWithCmd(authProvider, *eachData, &c.conf.ExtractCmd, c.conf.CloudType)
			if err != nil {
				return nil, err
			}

			checkerPropList = append(checkerPropList, eachData)
		}

	}

	return checkerPropList, nil
}

// getPropWithCmd: Extract Id, Name and properties from previous data acoording to ConfExtractCmd
//
// Priority:
// (TODO) 1. Extract from sub commands described in CmdChain (skip below)
// 2. Extract from jsonpath described in ExtractJsonPath (skip below)
// 3. Extract from cloud associated with cloudType
// @param: authProvider: IAuthProvider to provide profile of auth
// @param: previousData: Raw data or data extracted from previous command in the CmdChain
// @param: conf: Definition of the extraction commands
// @param: cloudType: Cloud type for additional data to be retrieve from
// @return: Extracted properties
// @return: Error
func getPropWithCmd(authProvider auth.IAuthProvider, previousData CheckerProp, conf *def.ConfExtractCmd, cloudType def.CloudType) (
	*CheckerProp, error) {
	// previousData is a copy of the original value, except for the pointer of Prop
	checkerProp := &previousData
	var err error

	if len(conf.IdConst) > 0 {
		checkerProp.Id = conf.IdConst
	} else if len(conf.IdJsonPath) > 0 {
		checkerProp.Id, err = internal.ParseJsonPathStr(checkerProp.Prop, conf.IdJsonPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get id: %w", err)
		}
	}
	if len(checkerProp.Id) == 0 {
		return nil, errors.New("invalid property, id is empty")
	}

	if len(conf.NameJsonPath) > 0 {
		checkerProp.Name, err = internal.ParseJsonPathStr(checkerProp.Prop, conf.NameJsonPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get name: %w", err)
		}
	}

	// Need to be checked before decommenting
	//
	// if len(conf.CmdChain) > 0 {
	// 	for _, subCmd := range conf.CmdChain {
	// 		checkerProp, err = getPropWithCmd(authProvider, *checkerProp, &subCmd, cloudType)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("failed to get prop with commands chain: %w", err)
	// 		}
	// 	}
	// } else
	if !reflect.DeepEqual(conf.ExtractJsonPath, def.ConfJsonPathCmd{}) {
		extractedProp, err := internal.ParseJsonPath(checkerProp.Prop, conf.ExtractJsonPath.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JsonPath: %w", err)
		}

		checkerProp.Prop = extractedProp
	} else {
		checkerProp.Prop, err = getPropWithCloud(authProvider, cloudType, checkerProp.Id, conf)
		if err != nil {
			return nil, err
		}
	}

	if conf.NormalizeId {
		splitRes := strings.Split(checkerProp.Id, "/")
		checkerProp.Id = splitRes[len(splitRes)-1]
	}

	return checkerProp, nil
}

func getPropWithCloud(authProvider auth.IAuthProvider, cloudType def.CloudType, id string, conf *def.ConfExtractCmd) (*json.RawMessage, error) {
	if authProvider == nil {
		return nil, errors.New("nil pointor of IAuthProvider of Checker")
	}

	switch cloudType {
	case def.TENCENT_CLOUD:
		if len(conf.IdParamName) == 0 {
			return nil, errors.New("missing IdParamName for getting prop from Tencent cloud")
		}

		extraParam := maps.Clone(conf.TencentCloud.ExtraParam)
		if extraParam == nil {
			extraParam = make(map[string]any)
		}
		if err := internal.AddParamString(extraParam, conf.IdParamName, id, conf.IdParamType); err != nil {
			return nil, err
		}

		return connector.CallTencentCloud(
			authProvider,
			conf.TencentCloud.Service,
			conf.TencentCloud.Version,
			conf.TencentCloud.Action,
			extraParam,
		)
	case def.TENCENT_COS:
		return connector.CallTencentCOS(
			authProvider,
			id, // Treat id as bucketName
			conf.TencentCOS.Service,
			conf.TencentCOS.Action,
		)
	case def.ALIYUN_CLOUD:
		if len(conf.IdParamName) == 0 {
			return nil, errors.New("missing IdParamName for getting prop from Aliyun")
		}

		extraParam := maps.Clone(conf.Aliyun.ExtraParam)
		if extraParam == nil {
			extraParam = make(map[string]any)
		}
		if err := internal.AddParamString(extraParam, conf.IdParamName, id, conf.IdParamType); err != nil {
			return nil, err
		}

		return connector.CallAliyunCloud(
			authProvider,
			conf.Aliyun.Endpoint,
			conf.Aliyun.EndpointWithRegion,
			conf.Aliyun.Version,
			conf.Aliyun.Action,
			extraParam,
		)
	case def.ALIYUN_OSS:
		return connector.CallAliyunOSS(
			authProvider,
			id, // Treat id as bucketName
			conf.AliyunOSS.Action,
			nil, // Ignored if not listing buckets
		)
	case def.AZURE:
		return connector.CallAzureWithEndpoint(
			authProvider,
			conf.Azure.Version,
			id,
			conf.Azure.Action,
		)
	default:
		return nil, fmt.Errorf("invalid cloud type: %s", cloudType)
	}
}

// ValidateResult: Result of validation
type ValidateResult struct {
	// Name of cloud
	CloudType def.CloudType
	// Resource identifier on the cloud
	Id string
	// Human readable name of the resource
	Name string
	// Indicate if the property has failed the benchmark check
	InRisk bool
	// Actual value of the property to be displayed
	Value string
}

// Validate: Validate properties and generate result
// @param: data: Properties extracted from the step of GetProp
// @return: Result of validation
// @return: Error
func (c *Checker) Validate(data CheckerPropList) ([]*ValidateResult, error) {
	if err := c.createValidator(); err != nil {
		return nil, err
	}
	if c.validator == nil {
		// createValidator has failed for more than once,
		// just return to avoid repeating the same errors
		return nil, nil
	}

	var validateResultList = make([]*ValidateResult, 0, len(data))
	for _, eachProp := range data {
		eachResult := ValidateResult{
			CloudType: c.conf.CloudType,
			Id:        eachProp.Id,
			Name:      eachProp.Name,
		}

		jsResult, err := c.validator.Validate(gojsonschema.NewBytesLoader(*eachProp.Prop))
		if err != nil {
			// Print error and skip the current property
			glog().Printf("Failed to validate prop: %v\n", err)
			continue
		}
		eachResult.InRisk = jsResult.Valid()

		if len(c.conf.Validator.ValueJsonPath) > 0 {
			eachResult.Value, err = internal.ParseJsonPathStr(eachProp.Prop, c.conf.Validator.ValueJsonPath)
			if err != nil {
				// Print error and skip the current property
				glog().Printf("Failed to get actual value with jsonPath: %v\n", err)
				continue
			}
		}

		validateResultList = append(validateResultList, &eachResult)
	}

	return validateResultList, nil
}

// createValidator: Try to create validator according to definition of Checker.conf.Validator
//
// If failed, the same error will not be repeated on the next call
// @return: Error
func (c *Checker) createValidator() error {
	if c.validator == nil && !c.validatorLoadFailed {
		schema := c.conf.Validator.ValidateSchema
		for k, v := range c.conf.Validator.DynValidateValue {
			schema = strings.ReplaceAll(schema, fmt.Sprintf("%%%s%%", k), v) // replace %k% to v
		}

		var err error
		c.validator, err = gojsonschema.NewSchema(gojsonschema.NewStringLoader(schema))
		if err != nil {
			c.validator = nil
			c.validatorLoadFailed = true
			return fmt.Errorf("failed to create jsonschema: %w", err)
		}
	}

	return nil
}
