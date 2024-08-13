// Connector for Tencent cloud

package connector

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"go.uber.org/ratelimit"
)

const (
	TENCENTCLOUD_SECRET_ID  = "TENCENTCLOUD_SECRET_ID"
	TENCENTCLOUD_SECRET_KEY = "TENCENTCLOUD_SECRET_KEY"
	TENCENTCLOUD_REGION     = "TENCENTCLOUD_REGION"
)

func createTencentCloudClient(p auth.IAuthProvider) (*common.Client, error) {
	if p == nil {
		return nil, errors.New("nil pointor of IAuthProvider")
	}

	v, err := p.GetProfile(def.TENCENT_CLOUD)
	if err != nil {
		return nil, err
	}
	if err := auth.IsAllSet(v, []string{TENCENTCLOUD_SECRET_ID, TENCENTCLOUD_SECRET_KEY, TENCENTCLOUD_REGION}); err != nil {
		return nil, err
	}

	credential := common.NewCredential(
		v.GetString(TENCENTCLOUD_SECRET_ID),
		v.GetString(TENCENTCLOUD_SECRET_KEY))
	cpf := profile.NewClientProfile()
	client := common.NewCommonClient(credential, v.GetString(TENCENTCLOUD_REGION), cpf)
	return client, nil
}

var (
	_mapTencentCloudClient internal.SyncMap[*common.Client]

	_rlTencentCloud = ratelimit.New(10, ratelimit.WithoutSlack)
)

func getTencentCloudClient(p auth.IAuthProvider) (*common.Client, error) {
	key := fmt.Sprintf("%p_default", p)
	return _mapTencentCloudClient.LoadOrCreate(key, func() (any, error) {
		return createTencentCloudClient(p)
	}, nil)
}

// CallTencentCloud: Send a request to Tencent cloud and parse response
// @param: authProvider: IAuthProvider to provide profile of auth
// @param: service: Parameter for Tencent cloud common request
// @param: version: Parameter for Tencent cloud common request
// @param: action: Parameter for Tencent cloud common request
// @param: extraParam: Extra Parameter provided to Tencent cloud
// @return: Response data from Tencent cloud
// @return: Error
func CallTencentCloud(authProvider auth.IAuthProvider, service string, version string, action string, extraParam map[string]any) (
	*json.RawMessage, error) {
	client, err := getTencentCloudClient(authProvider)
	if err != nil {
		return nil, err
	}

	request := tchttp.NewCommonRequest(service, version, action)
	if err := request.SetActionParameters(extraParam); err != nil {
		return nil, fmt.Errorf("failed to set extraParam: %w", err)
	}
	response := tchttp.NewCommonResponse()

	_rlTencentCloud.Take()
	if err := client.Send(request, response); err != nil {
		return nil, fmt.Errorf("failed to invoke api: %w", err)
	}

	responseMap := make(map[string]json.RawMessage, 1)
	if err := internal.JsonUnmarshal(response.GetBody(), &responseMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response as json: %w", err)
	}

	rawData, ok := responseMap["Response"]
	if !ok {
		return nil, errors.New("invalid response, missing key \"Response\"")
	}

	return &rawData, nil
}
