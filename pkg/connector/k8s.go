// Connector for k8s

package connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	"go.uber.org/ratelimit"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

// Bind dynamic.DynamicClient with RESTMapper
type k8sClient struct {
	c *dynamic.DynamicClient
	m meta.RESTMapper
	// Version of server
	v string
}

func createK8sClient(p auth.IAuthProvider) (*k8sClient, error) {
	if p == nil {
		return nil, errors.New("nil pointor of IAuthProvider")
	}

	kubeconfigPathname, err := p.GetProfilePathname(def.K8S)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPathname)
	if err != nil {
		// Do not use value of err to avoid leaking the file path
		return nil, errors.New("unable to read kube config file")
	}

	var client k8sClient
	client.c, err = dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	// GetAPIGroupResources() is called only once and the result is cached
	// TODO: Manage how to refresh restmapper
	groupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return nil, err
	}

	client.m = restmapper.NewDiscoveryRESTMapper(groupResources)

	// ServerVersion() is called only once and the result is cached
	serverVersion, err := discoveryClient.ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %w", err)
	}

	client.v = fmt.Sprintf("%s.%s", serverVersion.Major, serverVersion.Minor)

	return &client, nil
}

var (
	_mapK8sClient internal.SyncMap[*k8sClient]

	_rlK8sCloud = ratelimit.New(10, ratelimit.WithoutSlack)
)

func getK8sClient(p auth.IAuthProvider) (*k8sClient, error) {
	key := fmt.Sprintf("%p_default", p)
	return _mapK8sClient.LoadOrCreate(key, func() (any, error) {
		return createK8sClient(p)
	}, nil)
}

// CallK8sList: Send a request to a k8s server to list resources.
// If group and version are both empty, RESTMapper is used to search for mapped gvr
//
// NOTE: It's not necessary to define a function to get data of single resource,
// since all the data will be returned in the result of listing in the current usages
// @param: authProvider: IAuthProvider to provide pathname of kubeconfig
// @param: namespace: Parameter for k8s request
// @param: group: Parameter for k8s request
// @param: version: Parameter for k8s request
// @param: resource: Parameter for k8s request
// @param: extraParam: Parameters of ListOptions
// @return: Response data from k8s server
// @return: Error
func CallK8sList(authProvider auth.IAuthProvider, namespace string, group string, version string, resource string, listOpts map[string]any) (
	*json.RawMessage, error) {
	client, err := getK8sClient(authProvider)
	if err != nil {
		return nil, err
	}

	byListOpts, err := json.Marshal(listOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal listOpts: %w", err)
	}

	listOption := metav1.ListOptions{}
	if err := json.Unmarshal(byListOpts, &listOption); err != nil {
		return nil, fmt.Errorf("failed to unmarshal listOpts to listOption: %w", err)
	}

	gvr := schema.GroupVersionResource{Group: group, Version: version, Resource: resource}
	var rs dynamic.NamespaceableResourceInterface
	if len(group) == 0 && len(version) == 0 {
		found, err := client.m.ResourceFor(gvr)
		if err != nil {
			return nil, fmt.Errorf("failed to find resource \"%s\": %w", resource, err)
		}

		rs = client.c.Resource(found)
	} else {
		rs = client.c.Resource(schema.GroupVersionResource{Group: group, Version: version, Resource: resource})
	}

	var listRes *unstructured.UnstructuredList
	_rlK8sCloud.Take()
	if len(namespace) > 0 {
		listRes, err = rs.Namespace(namespace).List(context.TODO(), listOption)
	} else {
		listRes, err = rs.List(context.TODO(), listOption)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list k8s resource: %w", err)
	}

	return internal.JsonMarshal(listRes.UnstructuredContent())
}

// GetK8sVersion: Get version of a k8s server
//
// @param: authProvider: IAuthProvider to provide pathname of kubeconfig
// @return: Version of k8s server
// @return: Error
func GetK8sVersion(authProvider auth.IAuthProvider) (string, error) {
	client, err := getK8sClient(authProvider)
	if err != nil {
		return "", err
	}

	_ = internal.DisableInlining()

	return client.v, nil
}
