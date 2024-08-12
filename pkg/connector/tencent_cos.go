// Connector for Tencent COS

package connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	cos "github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/ratelimit"
)

func createTencentCOSClient(p auth.IAuthProvider, bucketName string) (*cos.Client, error) {
	if p == nil {
		return nil, errors.New("nil pointor of IAuthProvider")
	}

	v, err := p.GetProfile(def.TENCENT_COS)
	if err != nil {
		return nil, err
	}
	if err := auth.IsAllSet(v, []string{TENCENTCLOUD_SECRET_ID, TENCENTCLOUD_SECRET_KEY, TENCENTCLOUD_REGION}); err != nil {
		return nil, err
	}

	var u *url.URL
	if len(bucketName) > 0 {
		u, err = url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com",
			bucketName,
			v.GetString(TENCENTCLOUD_REGION),
		))
	} else {
		u, err = url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com",
			v.GetString(TENCENTCLOUD_REGION),
		))
	}
	if err != nil {
		panic(fmt.Sprintf("internal error when parsing url: %v", err))
	}

	client := cos.NewClient(
		&cos.BaseURL{BucketURL: u},
		&http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  v.GetString(TENCENTCLOUD_SECRET_ID),
				SecretKey: v.GetString(TENCENTCLOUD_SECRET_KEY),
			},
		})

	return client, nil
}

var (
	_mapTencentCOSClient internal.SyncMap[*cos.Client]

	_rlTencentCOS = ratelimit.New(10, ratelimit.WithoutSlack)
)

func getTencentCOSClient(authProvider auth.IAuthProvider, bucketName string) (*cos.Client, error) {
	key := fmt.Sprintf("%p_%s", authProvider, bucketName)
	return _mapTencentCOSClient.LoadOrCreate(key, func() (any, error) {
		return createTencentCOSClient(authProvider, bucketName)
	}, nil)
}

// CallTencentCOS: Send a request to Tencent COS and parse response
//
// TODO: Deal with extra parameters for different reflect call
// @param: authProvider: IAuthProvider to provide profile of auth
// @param: bucketName: Name of bucket. List all buckets if empty string given
// @param: service: Parameter for the reflection of Tencent COS API
// @param: action: Parameter for the reflection of Tencent COS API
// @return: Response data from Tencent COS
// @return: Error
func CallTencentCOS(authProvider auth.IAuthProvider, bucketName string, service string, action string) (
	*json.RawMessage, error) {
	client, err := getTencentCOSClient(authProvider, bucketName)
	if err != nil {
		return nil, err
	}

	if len(bucketName) == 0 {
		// Only Service.Get is available when listing buckets
		service = "Service"
		action = "Get"
	}

	clientType := reflect.TypeOf(client)
	clientValue := reflect.ValueOf(client)
	if _, found := clientType.Elem().FieldByName(service); !found {
		return nil, fmt.Errorf("service field not found on reflection of Tencent COS client: %s", service)
	}

	serviceValue := clientValue.Elem().FieldByName(service)
	serviceType := serviceValue.Type()

	if method, found := serviceType.MethodByName(action); !found {
		return nil, fmt.Errorf("action method not found on reflection of \"%s\" of Tencent COS client: %s", service, action)
	} else {
		_rlTencentCOS.Take()
		callResult := method.Func.Call([]reflect.Value{serviceValue, reflect.ValueOf(context.Background())})
		if len(callResult) != 3 {
			panic("internal error, invalid length of call results on reflection of Tencent COS client")
		}

		if callResult[2].Interface() != nil {
			if err, ok := callResult[2].Interface().(error); !ok {
				return nil, errors.New("failed to parse call result[2] of Tencent COS client as an error")
			} else {
				return nil, fmt.Errorf("failed to call \"%s.%s\" on Tencent COS: %w", service, action, err)
			}
		}

		return internal.JsonMarshal(callResult[0].Interface())
	}
}
