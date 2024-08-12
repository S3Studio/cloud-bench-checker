// Connector for Aliyun OSS

package connector

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/ratelimit"
)

func createAliyunOSSClient(p auth.IAuthProvider) (*oss.Client, error) {
	if p == nil {
		return nil, errors.New("nil pointor of IAuthProvider")
	}

	v, err := p.GetProfile(def.ALIYUN_OSS)
	if err != nil {
		return nil, err
	}
	if err := auth.IsAllSet(v, []string{ALIYUN_ACCESS_KEY_ID, ALIYUN_ACCESS_KEY_SECRET, ALIYUN_REGION}); err != nil {
		return nil, err
	}

	region := v.GetString(ALIYUN_REGION)
	if region[:4] != "oss-" {
		region = fmt.Sprintf("oss-%s", region)
	}
	return oss.New(
		fmt.Sprintf("https://%s.aliyuncs.com", region),
		v.GetString(ALIYUN_ACCESS_KEY_ID),
		v.GetString(ALIYUN_ACCESS_KEY_SECRET))
}

var (
	_mapAliyunOSSClient internal.SyncMap[*oss.Client]

	_rlAliyunOSS = ratelimit.New(10, ratelimit.WithoutSlack)
)

func getAliyunOSSClient(p auth.IAuthProvider) (*oss.Client, error) {
	key := fmt.Sprintf("%p_default", p)
	return _mapAliyunOSSClient.LoadOrCreate(key, func() (any, error) {
		return createAliyunOSSClient(p)
	}, nil)
}

const ALIYUN_OSS_MARKER_KEY = "marker"

// CallAliyunOSS: Send a request to Aliyun OSS and parse response
//
// TODO: Deal with more extra parameters for different reflect call
// @param: authProvider: IAuthProvider to provide profile of auth
// @param: bucketName: Name of the bucket. List all buckets if empty string given
// @param: action: Parameter for the reflection of Aliyun OSS API
// @param: extraParam: Currently only used to transfer marker when listing buckets
// @return: Response data from Aliyun OSS
// @return: Error
func CallAliyunOSS(authProvider auth.IAuthProvider, bucketName string, action string, extraParam map[string]any) (
	*json.RawMessage, error) {
	client, err := getAliyunOSSClient(authProvider)
	if err != nil {
		return nil, err
	}

	if len(bucketName) == 0 {
		// Only ListBuckets action available when listing buckets
		action = "ListBuckets"
	}

	clientType := reflect.TypeOf(client)
	clientValue := reflect.ValueOf(client)
	if method, found := clientType.MethodByName(action); !found {
		return nil, fmt.Errorf("action method not found on reflection of Aliyun OSS client: %s", action)
	} else {
		param := make([]reflect.Value, 2)
		param[0] = clientValue
		if action == "ListBuckets" {
			if marker, ok := extraParam[ALIYUN_OSS_MARKER_KEY].(string); ok {
				param[1] = reflect.ValueOf(oss.Marker(marker))
			}
		} else {
			param[1] = reflect.ValueOf(bucketName)
		}

		_rlAliyunOSS.Take()
		callResult := method.Func.Call(param)
		if len(callResult) != 2 {
			panic("internal error, invalid length of call results on reflection of Aliyun OSS client")
		}

		if callResult[1].Interface() != nil {
			if err, ok := callResult[1].Interface().(error); !ok {
				return nil, errors.New("failed to parse call result[1] of Aliyun OSS client as an error")
			} else {
				return nil, fmt.Errorf("failed to call \"%s\" on Aliyun OSS: %w", action, err)
			}
		}

		return internal.JsonMarshal(callResult[0].Interface())
	}
}
