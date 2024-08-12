// Integration test for apiserver
package apiserver_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/s3studio/cloud-bench-checker/internal"
	swaggerserversrv "github.com/s3studio/cloud-bench-checker/internal/server"
	"github.com/s3studio/cloud-bench-checker/internal/server/operations"
	"github.com/s3studio/cloud-bench-checker/internal/server/operations/baseline"
	"github.com/s3studio/cloud-bench-checker/internal/server/operations/listor"
	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	"github.com/s3studio/cloud-bench-checker/pkg/connector"
	"github.com/s3studio/cloud-bench-checker/pkg/server_model"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/go-openapi/loads"
)

const (
	TESTSUITE_DIR_NAME     = "testsuite"
	CONF_FILENAME          = "config.conf"
	CLOUD_LIST_RESULT_FILE = "TencentCloudListResult.json"
	TEST_RESULT_FILE       = "TestResult.json"
)

var (
	_testSuiteDir = ""
	_deferList    = []func(){}
	_testResult   = []*server_model.ValidateResult{}
)

func setupEnvironment() {
	if _testSuiteDir == "" {
		_, filename, _, _ := runtime.Caller(1)
		_testSuiteDir = filepath.Join(filepath.Dir(filename), TESTSUITE_DIR_NAME)
	}

	var patchOsOpen *gomonkey.Patches
	patchOsOpen = gomonkey.ApplyFunc(os.OpenFile,
		func(name string, flag int, perm os.FileMode) (*os.File, error) {
			var res *os.File
			var err error

			filename := filepath.Base(name)
			if filename == CONF_FILENAME && flag == os.O_RDONLY && perm == 0 {
				patchOsOpen.Origin(func() { res, err = os.Open(filepath.Join(_testSuiteDir, filename)) })
			} else {
				patchOsOpen.Origin(func() { res, err = os.OpenFile(name, flag, perm) })
			}
			return res, err
		})
	_deferList = append([]func(){patchOsOpen.Reset}, _deferList...)

	cloudListResultFile, _ := os.Open(filepath.Join(_testSuiteDir, CLOUD_LIST_RESULT_FILE))
	defer cloudListResultFile.Close()
	byFile, _ := io.ReadAll(cloudListResultFile)
	var v any
	json.Unmarshal(byFile, &v)
	rm, _ := internal.JsonMarshal(v)
	patches := gomonkey.ApplyFunc(connector.CallTencentCloud,
		func(authProvider auth.IAuthProvider, service string, version string, action string, extraParam map[string]any) (*json.RawMessage, error) {
			return rm, nil
		})
	_deferList = append([]func(){patches.Reset}, _deferList...)

	testResultFile, _ := os.Open(filepath.Join(_testSuiteDir, TEST_RESULT_FILE))
	defer testResultFile.Close()
	byFile, _ = io.ReadAll(testResultFile)
	json.Unmarshal(byFile, &_testResult)
}

func setupServer() http.Handler {
	swaggerSpec, err := loads.Embedded(swaggerserversrv.SwaggerJSON, swaggerserversrv.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewCloudBenchCheckerAPIAPI(swaggerSpec)
	server := swaggerserversrv.NewServer(api)
	defer server.Shutdown()

	// parser := flags.NewParser(server, flags.Default)
	// parser.ShortDescription = "Cloud Bench Checker API"
	// parser.LongDescription = "API for https://github.com/S3Studio/cloud-bench-checker described with Swagger 2.0\n"
	// server.ConfigureFlags()
	// for _, optsGroup := range api.CommandLineOptionsGroups {
	// 	_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }

	// if _, err := parser.Parse(); err != nil {
	// 	code := 1
	// 	if fe, ok := err.(*flags.Error); ok {
	// 		if fe.Type == flags.ErrHelp {
	// 			code = 0
	// 		}
	// 	}
	// 	os.Exit(code)
	// }

	server.ConfigureAPI()

	return server.GetHandler()
}

func TestApiserverIntegration(t *testing.T) {
	setupEnvironment()
	defer func() {
		for _, eachFn := range _deferList {
			eachFn()
		}
	}()

	handler := setupServer()

	baselineId := 0
	t.Run("GET /baseline/getIds", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/baseline/getIds", nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("%s = %d, want %d", t.Name(), resp.Code, http.StatusOK)
			return
		}

		body := baseline.GetBaselineGetIdsOKBody{}
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Errorf("%s json.Unmarshal error: %v", t.Name(), err)
			return
		}

		if len(body.Data) != 3 {
			t.Errorf("%s len(ids) = %d, want 3", t.Name(), len(body.Data))
			return
		}

		foundId := false
		idValue := 1
		for _, id := range body.Data {
			if id == int64(idValue) {
				foundId = true
				break
			}
		}

		if !foundId {
			t.Errorf("%s id of %d not returned", t.Name(), idValue)
			return
		}

		baselineId = idValue
	})

	listorId := 0
	t.Run("GET /baseline/getListorId", func(t *testing.T) {
		if baselineId == 0 {
			t.Fatal("Preconditions not met")
		}

		req := httptest.NewRequest(
			"GET",
			fmt.Sprintf("/api/baseline/getListorId?id=%d", baselineId),
			nil,
		)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("%s = %d, want %d", t.Name(), resp.Code, http.StatusOK)
			return
		}

		body := baseline.GetBaselineGetListorIDOKBody{}
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Errorf("%s json.Unmarshal error: %v", t.Name(), err)
			return
		}

		if len(body.Data) != 2 {
			t.Errorf("%s len(ids) = %d, want 2", t.Name(), len(body.Data))
			return
		}

		foundId := false
		idValue := 1
		for _, id := range body.Data {
			if id == int64(idValue) {
				foundId = true
				break
			}
		}

		if !foundId {
			t.Errorf("%s id of %d not returned", t.Name(), idValue)
			return
		}

		listorId = idValue
	})

	t.Run("GET /listor/getDefinition", func(t *testing.T) {
		if listorId == 0 {
			t.Fatal("Preconditions not met")
		}

		req := httptest.NewRequest(
			"GET",
			fmt.Sprintf("/api/listor/getDefinition?id=%d", listorId),
			nil,
		)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("%s = %d, want %d", t.Name(), resp.Code, http.StatusOK)
			return
		}

		body := listor.GetListorGetDefinitionOKBody{}
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Errorf("%s json.Unmarshal error: %v", t.Name(), err)
			return
		}

		if body.Data.CloudType != "tencent_cloud" {
			t.Errorf("%s Data.CloudType = \"%s\", want \"tencent_cloud\"", t.Name(), body.Data.CloudType)
			return
		}
	})

	var resListData *server_model.ListorData
	t.Run("GET /listor/listData", func(t *testing.T) {
		if listorId == 0 {
			t.Fatal("Preconditions not met")
		}

		req := httptest.NewRequest(
			"GET",
			fmt.Sprintf("/api/listor/listData?id=%d", listorId),
			nil,
		)
		req.Header.Set("profile", "mock")
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("%s = %d, want %d", t.Name(), resp.Code, http.StatusOK)
			return
		}

		body := listor.GetListorListDataOKBody{}
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Errorf("%s json.Unmarshal error: %v", t.Name(), err)
			return
		}

		if body.Data.CloudType != "tencent_cloud" {
			t.Errorf("%s Data.CloudType = \"%s\", want \"tencent_cloud\"", t.Name(), body.Data.CloudType)
			return
		}

		if body.Data.ListorID != int64(listorId) {
			t.Errorf("%s Data.ListorID = %d, want %d", t.Name(), body.Data.ListorID, listorId)
			return
		}

		resListData = body.Data
	})

	var resGetProp *server_model.BaselineData
	t.Run("POST /baseline/getProp", func(t *testing.T) {
		if resListData == nil {
			t.Fatal("Preconditions not met")
		}

		// remember to merge previous result to list
		jsonListData, _ := json.Marshal([]server_model.ListorData{*resListData})
		req := httptest.NewRequest(
			"POST",
			fmt.Sprintf("/api/baseline/getProp?id=%d", baselineId),
			bytes.NewReader(jsonListData),
		)
		req.Header.Set("profile", "mock")
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("%s = %d, want %d", t.Name(), resp.Code, http.StatusOK)
			return
		}

		body := baseline.PostBaselineGetPropOKBody{}
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Errorf("%s json.Unmarshal error: %v", t.Name(), err)
			return
		}

		if body.Data.ID != int64(baselineId) {
			t.Errorf("%s Data.ID = %d, want %d", t.Name(), body.Data.ID, baselineId)
			return
		}

		resGetProp = body.Data
	})

	t.Run("POST /baseline/validate", func(t *testing.T) {
		if resGetProp == nil {
			t.Fatal("Preconditions not met")
		}

		jsonListData, _ := json.Marshal(resGetProp)
		req := httptest.NewRequest(
			"POST",
			fmt.Sprintf("/api/baseline/validate?id=%d", baselineId),
			bytes.NewReader(jsonListData),
		)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("%s = %d, want %d", t.Name(), resp.Code, http.StatusOK)
			return
		}

		body := baseline.PostBaselineValidateOKBody{}
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Errorf("%s json.Unmarshal error: %v", t.Name(), err)
			return
		}

		if !reflect.DeepEqual(body.Data, _testResult) {
			t.Errorf("%s = %v, want %v", t.Name(), body.Data, _testResult)
			return
		}
	})

	// The following API is not necessary for client-side
	t.Run("GET /listor/getIds", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/listor/getIds", nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("%s = %d, want %d", t.Name(), resp.Code, http.StatusOK)
		}

		body := listor.GetListorGetIdsOKBody{}
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Errorf("%s json.Unmarshal error: %v", t.Name(), err)
			return
		}

		if len(body.Data) != 4 {
			t.Errorf("%s len(ids) = %d, want 4", t.Name(), len(body.Data))
			return
		}
	})

	t.Run("GET /baseline/getDefinition", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/baseline/getDefinition?id=1&with_hash=true&with_yaml=true", nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("%s = %d, want %d", t.Name(), resp.Code, http.StatusOK)
		}

		body := baseline.GetBaselineGetDefinitionOKBody{}
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Errorf("%s json.Unmarshal error: %v", t.Name(), err)
			return
		}
	})
}
