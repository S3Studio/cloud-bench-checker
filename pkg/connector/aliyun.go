// Connector for Aliyun
package connector

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/ratelimit"
)

const (
	ALIYUN_ACCESS_KEY_ID     = "ALIBABA_CLOUD_ACCESS_KEY_ID"
	ALIYUN_ACCESS_KEY_SECRET = "ALIBABA_CLOUD_ACCESS_KEY_SECRET"
	ALIYUN_REGION            = "ALIBABA_CLOUD_REGION"
)

func createAliyunCloudClient(p auth.IAuthProvider, endpoint string, bEpWithRegion bool) (*openapi.Client, error) {
	if p == nil {
		return nil, errors.New("nil pointor of IAuthProvider")
	}

	v, err := p.GetProfile(def.ALIYUN_CLOUD)
	if err != nil {
		return nil, err
	}
	if err := auth.IsAllSet(v, []string{ALIYUN_ACCESS_KEY_ID, ALIYUN_ACCESS_KEY_SECRET, ALIYUN_REGION}); err != nil {
		return nil, err
	}

	config := &openapi.Config{
		AccessKeyId:     tea.String(v.GetString(ALIYUN_ACCESS_KEY_ID)),
		AccessKeySecret: tea.String(v.GetString(ALIYUN_ACCESS_KEY_SECRET)),
	}
	if bEpWithRegion {
		config.Endpoint = tea.String(fmt.Sprintf("%s.%s.aliyuncs.com", endpoint, v.GetString(ALIYUN_REGION)))
	} else {
		config.Endpoint = tea.String(fmt.Sprintf("%s.aliyuncs.com", endpoint))
	}

	return openapi.NewClient(config)
}

var (
	_mapAliyunCloudClient internal.SyncMap[*openapi.Client]

	_rlAliyunCloud = ratelimit.New(10, ratelimit.WithoutSlack)
)

func getAliyunCloudClient(p auth.IAuthProvider, endpoint string, bEpWithRegion bool) (*openapi.Client, error) {
	key := fmt.Sprintf("%p_%s", p, endpoint)
	return _mapAliyunCloudClient.LoadOrCreate(key, func() (any, error) {
		return createAliyunCloudClient(p, endpoint, bEpWithRegion)
	}, nil)
}

// CallAliyunCloud: Send a request to Aliyun and parse response
//
// TODO: Deal with more types of extraParam for tea
// @param: authProvider: IAuthProvider to provide profile of auth
// @param: endpoint: Parameter for Aliyun common request
// @param: bEpWithRegion: Indicate if region should be added to endpoint
// @param: version: Parameter for Aliyun common request
// @param: action: Parameter for Aliyun common request
// @param: extraParam: Extra parameters provided to Aliyun
// @return: Response data from Aliyun
// @return: Error
func CallAliyunCloud(authProvider auth.IAuthProvider, endpoint string, bEpWithRegion bool, version string, action string, extraParam map[string]any) (
	*json.RawMessage, error) {
	client, err := getAliyunCloudClient(authProvider, endpoint, bEpWithRegion)
	if err != nil {
		return nil, err
	}

	v, err := authProvider.GetProfile(def.ALIYUN_CLOUD)
	if err != nil {
		return nil, err
	}

	params := &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Pathname:    tea.String("/"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}

	queries := make(map[string]any)
	queries["RegionId"] = tea.String(v.GetString(ALIYUN_REGION))
	for k, v := range extraParam {
		switch p := v.(type) {
		case string:
			queries[k] = tea.String(p)
		case int:
			queries[k] = tea.Int(p)
		}
	}

	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}

	_rlAliyunCloud.Take()
	response, err := client.CallApi(params, request, runtime)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke api: %w", err)
	}

	responseBody, ok := response["body"]
	if !ok {
		return nil, errors.New("invalid response, missing key \"body\"")
	}

	return internal.JsonMarshal(responseBody)
}
