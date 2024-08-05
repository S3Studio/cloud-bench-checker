// Connector for k8s

package connector

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"

	"github.com/agiledragon/gomonkey/v2"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func TestCallK8sList(t *testing.T) {
	listRes := unstructured.UnstructuredList{}
	var patchList *gomonkey.Patches
	patchGetK8sClient := gomonkey.ApplyFunc(getK8sClient,
		func(p auth.IAuthProvider) (*k8sClient, error) {
			c := &dynamic.DynamicClient{}
			orig := c.Resource(schema.GroupVersionResource{})
			patchList = gomonkey.ApplyMethodFunc(orig, "List",
				func(ctx context.Context, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
					return &listRes, nil
				})
			m := meta.PriorityRESTMapper{}
			return &k8sClient{c, &m}, nil
		})
	defer patchGetK8sClient.Reset()
	patchResourceFor := gomonkey.ApplyMethodFunc(meta.PriorityRESTMapper{}, "ResourceFor",
		func(input schema.GroupVersionResource) (schema.GroupVersionResource, error) {
			if input.Resource == "invalid" {
				return schema.GroupVersionResource{}, errors.New("mock invalid resource err")
			}
			return schema.GroupVersionResource{}, nil
		})
	defer patchResourceFor.Reset()
	defer func() {
		if patchList != nil {
			patchList.Reset()
		}
	}()
	rmList, _ := internal.JsonMarshal((&listRes).UnstructuredContent())

	type args struct {
		authProvider auth.IAuthProvider
		namespace    string
		group        string
		version      string
		resource     string
		extraParam   map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    *json.RawMessage
		wantErr bool
	}{
		{
			"Valid result with empty param",
			args{nil, "", "", "", "mock_rs", make(map[string]any)},
			rmList,
			false,
		},
		{
			"Valid result with non-empty version",
			args{nil, "", "", "mock_version", "mock_rs", make(map[string]any)},
			rmList,
			false,
		},
		{
			"Valid result using non-empty namespace",
			args{nil, "mock_ns", "", "", "mock_rs", make(map[string]any)},
			rmList,
			false,
		},
		{
			"failed to unmarshal extraParam to listOption",
			args{nil, "", "", "", "mock_rs", map[string]any{"labelSelector": []int{}}},
			nil,
			true,
		},
		{
			"failed to find resource",
			args{nil, "", "", "", "invalid", make(map[string]any)},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CallK8sList(tt.args.authProvider, tt.args.namespace, tt.args.group, tt.args.version, tt.args.resource, tt.args.extraParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("CallK8sList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				sgot, _ := got.MarshalJSON()
				swant, _ := tt.want.MarshalJSON()
				t.Errorf("CallK8sList() = %v, want %v", string(sgot), string(swant))
			}
		})
	}
}
