// Connector for Azure using Azure Resource Manager(arm)
package connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"go.uber.org/ratelimit"
)

const (
	AZURE_CLIENT_ID       = "AZURE_CLIENT_ID"
	AZURE_TENANT_ID       = "AZURE_TENANT_ID"
	AZURE_CLIENT_SECRET   = "AZURE_CLIENT_SECRET"
	AZURE_SUBSCRIPTION_ID = "AZURE_SUBSCRIPTION_ID"
)

func createAzureClient(p auth.IAuthProvider) (*arm.Client, error) {
	v, err := p.GetProfile(def.AZURE)
	if err != nil {
		return nil, err
	}
	if err := auth.IsAllSet(v, []string{AZURE_CLIENT_ID, AZURE_TENANT_ID, AZURE_CLIENT_SECRET, AZURE_SUBSCRIPTION_ID}); err != nil {
		return nil, err
	}

	credential, err := azidentity.NewClientSecretCredential(
		v.GetString(AZURE_TENANT_ID),
		v.GetString(AZURE_CLIENT_ID),
		v.GetString(AZURE_CLIENT_SECRET),
		nil)
	if err != nil {
		return nil, err
	}

	return arm.NewClient("cloud-bench-checker", "v0.0.1", credential, nil)
}

var (
	_mapAzureClient sync.Map

	_rlAzure = ratelimit.New(10, ratelimit.WithoutSlack)
)

func getAzureClient(p auth.IAuthProvider) (*arm.Client, error) {
	key := fmt.Sprintf("%p_default", p)
	client, ok := _mapAzureClient.Load(key)
	if !ok {
		newClient, err := createAzureClient(p)
		if err != nil {
			return nil, fmt.Errorf("failed to create Azure client: %w", err)
		}
		// May have already been created by other goroutions,
		// but it's ok to spend a little more time creating them
		client, _ = _mapAzureClient.LoadOrStore(key, newClient)
	}

	return client.(*arm.Client), nil
}

// Need to be checked before decommenting
// func CallAzure(authProvider auth.IAuthProvider, provider string, version string, rsType string, rsName string, action string) (
// 	*json.RawMessage, error) {
// 	v, err := authProvider.GetProfile(def.AZURE)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	endpoint := fmt.Sprintf("/subscriptions/%s/providers/Microsoft.%s/%s",
// 		v.GetString(AZURE_SUBSCRIPTION_ID),
// 		provider, rsType)
// 	if len(rsName) > 0 {
// 		endpoint = fmt.Sprintf("%s/%s", endpoint, rsName)
// 		if len(action) > 0 {
// 			endpoint = fmt.Sprintf("%s/%s", endpoint, action)
// 		}
// 	}
//
// 	return CallAzureWithEndpoint(authProvider, version, endpoint, "")
// }

// CallAzureList: Send a request to Azure to list resources
//
// TODO: Deal with extra parameters for Azure request
// @param: authProvider: IAuthProvider to provide profile of auth
// @param: provider: Parameter for Azure common request
// @param: version: Parameter for Azure common request
// @param: rsType: Parameter for Azure common request
// @param: nextLink: Returned from the previous call for pagination
// @return: Response data from Azure
// @return: Error
func CallAzureList(authProvider auth.IAuthProvider, provider string, version string, rsType string, nextLink string) (
	*json.RawMessage, error) {
	var endpoint string
	if len(nextLink) > 0 {
		// treat nextLink as endpoint
		endpoint = nextLink
	} else {
		v, err := authProvider.GetProfile(def.AZURE)
		if err != nil {
			return nil, err
		}

		endpoint = fmt.Sprintf("/subscriptions/%s/providers/Microsoft.%s/%s",
			v.GetString(AZURE_SUBSCRIPTION_ID),
			provider, rsType)
	}

	// ignore action when listing
	return CallAzureWithEndpoint(authProvider, version, endpoint, "")
}

// CallAzureList: Send a request to Azure with an endpoint provided
//
// The endpoint may be returned from the previous call as nextLink or resource id.
// Other functions like CallAzureList also concats endpoint string from their own parameters.
//
// TODO: Deal with extra parameters for Azure request
// @param: authProvider: IAuthProvider to provide profile of auth
// @param: version: Parameter for Azure common request
// @param: endpoint: Parameter for Azure common request
// @param: action: Parameter for Azure common request
// @return: Response data from Azure
// @return: Error
func CallAzureWithEndpoint(authProvider auth.IAuthProvider, version string, endpoint string, action string) (
	*json.RawMessage, error) {
	if len(endpoint) == 0 {
		return nil, errors.New("endpoint for Azure is empty")
	}

	if len(action) > 0 {
		endpoint = fmt.Sprintf("%s/%s", endpoint, action)
	}

	client, err := getAzureClient(authProvider)
	if err != nil {
		return nil, err
	}

	URL := endpoint
	if endpoint[:4] != "http" {
		URL = runtime.JoinPaths(client.Endpoint(), endpoint)
	}

	ctx := context.Background()
	req, err := runtime.NewRequest(ctx, http.MethodGet, URL)
	if err != nil {
		return nil, err
	}

	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", version)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}

	_rlAzure.Take()
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke api: %w", err)
	}

	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return nil, fmt.Errorf("response indicates failure: %w", runtime.NewResponseError(resp))
	}

	defer resp.Body.Close()
	byResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	responseMap := make(map[string]json.RawMessage, 1)
	if err := internal.JsonUnmarshal(byResp, &responseMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response as json: %w", err)
	}

	var rm json.RawMessage = byResp
	return &rm, nil
}
