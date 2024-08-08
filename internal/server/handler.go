// Actual handlers of server to interact with /pkg
package server

import (
	"crypto"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/s3studio/cloud-bench-checker/internal"
	"github.com/s3studio/cloud-bench-checker/internal/server/operations"
	"github.com/s3studio/cloud-bench-checker/internal/server/operations/baseline"
	"github.com/s3studio/cloud-bench-checker/internal/server/operations/listor"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"
	"github.com/s3studio/cloud-bench-checker/pkg/framework"
	"github.com/s3studio/cloud-bench-checker/pkg/server_model"

	"github.com/go-openapi/runtime/middleware"
	"gopkg.in/yaml.v2"
)

type logger func(string, ...interface{})

// Logf logs message either via defined user logger or via system one if no user logger is defined.
func logf(l logger, f string, args ...interface{}) {
	if l != nil {
		l(f, args...)
	} else {
		log.Printf(f, args...)
	}
}

type generalError struct {
	Code int
	Msg  string
	Data struct{}
}

const CONF_FILE_NAME = "config.conf"

var (
	_conf      def.ConfFile
	_confValid bool = false

	_mapInstance internal.SyncMap[any]
)

func getInstance(insType string, id int, fnCreate func() (any, error)) (any, error) {
	key := fmt.Sprintf("%s_%d", insType, id)
	return _mapInstance.LoadOrCreate(key, fnCreate, nil)
}

func getBaseline(id int) (*framework.Baseline, error) {
	pIns, err := getInstance("baseline", id, func() (any, error) {
		return framework.NewBaseline(&_conf.Baseline[id], nil, nil), nil
	})

	return pIns.(*framework.Baseline), err
}

func getListor(id int) (*framework.Listor, error) {
	pIns, err := getInstance("listor", id, func() (any, error) {
		for _, l := range _conf.Listor {
			if id == l.Id {
				return framework.NewListor(&l, nil), nil
			}
		}
		return nil, fmt.Errorf("listor not found with id: %d", id)
	})

	return pIns.(*framework.Listor), err
}

func getBaselineHash(id int) (*[]byte, error) {
	pIns, err := getInstance("baseline_hash", id, func() (any, error) {
		b, err := getBaseline(id)
		if err != nil {
			return nil, err
		}

		bConf := &_conf.Baseline[id]
		listorHash := make([][]*[]byte, len(bConf.Checker))
		for i, c := range bConf.Checker {
			listorHash[i] = make([]*[]byte, len(c.Listor))
			for j, id := range c.Listor {
				eachHash, err := getListorHash(id)
				if err != nil {
					return nil, err
				}

				listorHash[i][j] = eachHash
			}
		}

		by, err := b.GetHash(crypto.SHA256, listorHash)
		if err != nil {
			return nil, err
		}

		return &by, nil
	})

	return pIns.(*[]byte), err
}

func getListorHash(id int) (*[]byte, error) {
	pIns, err := getInstance("listor_hash", id, func() (any, error) {
		l, err := getListor(id)
		if err != nil {
			return nil, err
		}

		by, err := l.GetHash(crypto.SHA256)
		if err != nil {
			return nil, err
		}

		return &by, nil
	})

	return pIns.(*[]byte), err
}

func prepareHandler(api *operations.CloudBenchCheckerAPIAPI) {
	binPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		logf(api.Logger, "failed to get binary path: %v\n", err)
		return
	}
	binDir, _ := filepath.Split(binPath)
	confFilePath := filepath.Join(binDir, CONF_FILE_NAME)

	if confFile, err := os.Open(confFilePath); err != nil {
		logf(api.Logger, "Failed to open \"%s\": %v\n", confFilePath, err)
		return
	} else {
		defer confFile.Close()

		confFileContent, err := io.ReadAll(confFile)
		if err != nil {
			logf(api.Logger, "Failed to read \"%s\": %v\n", confFilePath, err)
			return
		}

		if err = yaml.Unmarshal(confFileContent, &_conf); err != nil {
			logf(api.Logger, "Failed to load \"%s\" as yaml: %v\n", confFilePath, err)
			return
		}
	}

	_confValid = true
}

func baselineGetBaselineGetDefinitionHandler(params baseline.GetBaselineGetDefinitionParams) middleware.Responder {
	if !_confValid {
		return middleware.Error(500, generalError{Code: 500, Msg: "config.conf file not loaded"})
	}

	if params.ID <= 0 || params.ID > int64(len(_conf.Baseline)) {
		return baseline.NewGetBaselineGetDefinitionNotFound().WithPayload(
			generalError{Code: 404, Msg: "Id not found"})
	}

	b := _conf.Baseline[params.ID-1]
	b4api := server_model.Baseline4api{
		ID:       params.ID,
		Tag:      b.Tag,
		Metadata: b.Metadata,
		Checker:  make([]*server_model.Checker4api, len(b.Checker)),
	}
	for i, c := range b.Checker {
		b4api.Checker[i] = &server_model.Checker4api{
			CloudType: server_model.Cloudtype4api(c.CloudType),
			Listor:    make([]int64, len(c.Listor)),
		}
		for index, id := range c.Listor {
			b4api.Checker[i].Listor[index] = int64(id)
		}
	}

	if params.WithHash != nil && *params.WithHash {
		byHash, err := getBaselineHash(int(params.ID - 1))
		if err != nil {
			return middleware.Error(500, generalError{
				Code: 500,
				Msg:  fmt.Sprintf("Failed to get Baseline hash: %v", err),
			})
		}

		b4api.Hash = &server_model.ItemHash{Sha256: fmt.Sprintf("%x", *byHash)}
	}

	if params.WithYaml != nil && *params.WithYaml {
		b4api.YamlHidden = _conf.Option.ServerHideYaml
		if !_conf.Option.ServerHideYaml {
			byBaseline, err := yaml.Marshal(b)
			if err != nil {
				return middleware.Error(500, generalError{
					Code: 500,
					Msg:  fmt.Sprintf("Failed to marshal to yaml: %v", err),
				})
			}
			b4api.Yaml = string(byBaseline)
		} else {
			b4api.Yaml = ""
		}
	}

	return baseline.NewGetBaselineGetDefinitionOK().WithPayload(
		&baseline.GetBaselineGetDefinitionOKBody{
			Code: 200,
			Msg:  "success",
			Data: &b4api,
		})
}

func baselineGetBaselineGetIdsHandler(params baseline.GetBaselineGetIdsParams) middleware.Responder {
	if !_confValid {
		return middleware.Error(500, generalError{Code: 500, Msg: "config.conf file not loaded"})
	}

	var ids []int
	for i, c := range _conf.Baseline {
		bAdd := false
	TagFilterLoop:
		for _, tagBaseline := range c.Tag {
			if len(params.Tag) == 0 {
				// Return all Baseline if filter is empty
				bAdd = true
			} else {
				for _, tagFilter := range params.Tag {
					if tagBaseline == tagFilter {
						bAdd = true
						break TagFilterLoop
					}
				}
			}
		}
		if bAdd {
			ids = append(ids, i+1)
		}
	}

	idsRes := make([]int64, len(ids))
	for i, v := range ids {
		idsRes[i] = int64(v)
	}

	return baseline.NewGetBaselineGetIdsOK().WithPayload(
		&baseline.GetBaselineGetIdsOKBody{
			Code: 200,
			Msg:  "Success",
			Data: idsRes,
		})
}

func baselineGetBaselineGetListorIDHandler(params baseline.GetBaselineGetListorIDParams) middleware.Responder {
	if !_confValid {
		return middleware.Error(500, generalError{Code: 500, Msg: "config.conf file not loaded"})
	}

	if params.ID <= 0 || params.ID > int64(len(_conf.Baseline)) {
		return baseline.NewGetBaselineGetListorIDNotFound().WithPayload(
			generalError{Code: 404, Msg: "Id not found"})
	}

	bIns, err := getBaseline(int(params.ID - 1))
	if err != nil {
		return middleware.Error(500, generalError{
			Code: 500,
			Msg:  fmt.Sprintf("Failed to get Baseline instance: %v", err),
		})
	}

	ids := bIns.GetListorId()
	idsRes := make([]int64, len(ids))
	for i, v := range ids {
		idsRes[i] = int64(v)
	}

	return baseline.NewGetBaselineGetListorIDOK().WithPayload(
		&baseline.GetBaselineGetListorIDOKBody{
			Code: 200,
			Msg:  "Success",
			Data: idsRes,
		})
}

func listorGetListorGetDefinitionHandler(params listor.GetListorGetDefinitionParams) middleware.Responder {
	if !_confValid {
		return middleware.Error(500, generalError{Code: 500, Msg: "config.conf file not loaded"})
	}

	var l *def.ConfListor
	for _, item := range _conf.Listor {
		if params.ID == int64(item.Id) {
			l = &item
		}
	}
	if l == nil {
		return listor.NewGetListorGetDefinitionNotFound().WithPayload(
			generalError{Code: 404, Msg: "Id not found"})
	}

	l4api := server_model.Listor4api{
		ID:        params.ID,
		CloudType: server_model.Cloudtype4api(l.CloudType),
		RsType:    l.RsType,
	}

	if params.WithHash != nil && *params.WithHash {
		byHash, err := getListorHash(int(params.ID))
		if err != nil {
			return middleware.Error(500, generalError{
				Code: 500,
				Msg:  fmt.Sprintf("Failed to get Listor hash: %v", err),
			})
		}

		l4api.Hash = &server_model.ItemHash{Sha256: fmt.Sprintf("%x", *byHash)}
	}

	if params.WithYaml != nil && *params.WithYaml {
		l4api.YamlHidden = _conf.Option.ServerHideYaml
		if !_conf.Option.ServerHideYaml {
			byListor, err := yaml.Marshal(l)
			if err != nil {
				return middleware.Error(500, generalError{
					Code: 500,
					Msg:  fmt.Sprintf("Failed to marshal to yaml: %v", err),
				})
			}
			l4api.Yaml = string(byListor)
		} else {
			l4api.Yaml = ""
		}
	}

	return listor.NewGetListorGetDefinitionOK().WithPayload(
		&listor.GetListorGetDefinitionOKBody{
			Code: 200,
			Msg:  "success",
			Data: &l4api,
		})
}

func listorGetListorGetIdsHandler(params listor.GetListorGetIdsParams) middleware.Responder {
	if !_confValid {
		return middleware.Error(500, generalError{Code: 500, Msg: "config.conf file not loaded"})
	}

	var ids []int
	for _, l := range _conf.Listor {
		if params.CloudType == nil || *params.CloudType == string(l.CloudType) {
			bAdd := true
			for _, idExist := range ids {
				if l.Id == idExist {
					bAdd = false
					break
				}
			}

			if bAdd {
				ids = append(ids, l.Id)
			}
		}
	}

	idsRes := make([]int64, len(ids))
	for i, v := range ids {
		idsRes[i] = int64(v)
	}

	return listor.NewGetListorGetIdsOK().WithPayload(
		&listor.GetListorGetIdsOKBody{
			Code: 200,
			Msg:  "success",
			Data: idsRes,
		})
}

func listorGetListorListDataHandler(params listor.GetListorListDataParams) middleware.Responder {
	if !_confValid {
		return middleware.Error(500, generalError{Code: 500, Msg: "config.conf file not loaded"})
	}

	var l *def.ConfListor
	for _, item := range _conf.Listor {
		if params.ID == int64(item.Id) {
			l = &item
		}
	}
	if l == nil {
		return listor.NewGetListorListDataNotFound().WithPayload(
			generalError{Code: 404, Msg: "Id not found"})
	}

	lIns, err := getListor(int(params.ID))
	if err != nil {
		return middleware.Error(500, generalError{
			Code: 500,
			Msg:  fmt.Sprintf("Failed to get Listor instance: %v", err),
		})
	}

	hash, err := getListorHash(int(params.ID))
	if err != nil {
		return middleware.Error(500, generalError{
			Code: 500,
			Msg:  fmt.Sprintf("Failed to get Listor hash: %v", err),
		})
	}

	data, err := lIns.ListData(framework.SetListorAuthProvider(&serverAuthProvider{params.Profile}))
	if err != nil {
		return listor.NewGetListorListDataBadRequest().WithPayload(
			generalError{Code: 400, Msg: fmt.Sprintf("Failed in ListData: %v", err)})
	}

	byData, err := json.Marshal(data)
	if err != nil {
		return middleware.Error(500, generalError{
			Code: 500,
			Msg:  fmt.Sprintf("Failed to marshal to json: %v", err),
		})
	}

	data4api := server_model.ListorData{
		ListorID:   params.ID,
		ListorHash: &server_model.ItemHash{Sha256: fmt.Sprintf("%x", *hash)},
		CloudType:  server_model.Cloudtype4api(l.CloudType),
		Data:       string(byData),
	}

	return listor.NewGetListorListDataOK().WithPayload(
		&listor.GetListorListDataOKBody{
			Code: 200,
			Msg:  "success",
			Data: &data4api,
		})
}

func baselinePostBaselineGetPropHandler(params baseline.PostBaselineGetPropParams) middleware.Responder {
	if !_confValid {
		return middleware.Error(500, generalError{Code: 500, Msg: "config.conf file not loaded"})
	}

	if params.ID <= 0 || params.ID > int64(len(_conf.Baseline)) {
		return baseline.NewPostBaselineGetPropNotFound().WithPayload(
			generalError{Code: 404, Msg: "Id not found"})
	}

	b := _conf.Baseline[params.ID-1]
	bIns, err := getBaseline(int(params.ID - 1))
	if err != nil {
		return middleware.Error(500, generalError{
			Code: 500,
			Msg:  fmt.Sprintf("Failed to get Baseline instance: %v", err),
		})
	}

	dataProvider := listDataProvider{params.ListorData}
	for _, c := range b.Checker {
		for _, id := range c.Listor {
			if provideCloudType, err := dataProvider.GetCloudTypeByListorId(id); err != nil {
				return baseline.NewPostBaselineGetPropBadRequest().WithPayload(
					generalError{Code: 400, Msg: fmt.Sprintf("%v", err)})
			} else if provideCloudType != string(c.CloudType) {
				return baseline.NewPostBaselineGetPropBadRequest().WithPayload(
					generalError{Code: 400, Msg: fmt.Sprintf("Cloud type mismatch between %s and %s",
						provideCloudType, c.CloudType)})
			}

			providedListorHash, err := dataProvider.GetListorHashByListorId(id)
			if err != nil {
				return baseline.NewPostBaselineGetPropBadRequest().WithPayload(
					generalError{Code: 400, Msg: fmt.Sprintf("%v", err)})
			}
			byHash, err := getListorHash(id)
			if err != nil {
				return middleware.Error(500, generalError{
					Code: 500,
					Msg:  fmt.Sprintf("Failed to get Listor hash: %v", err),
				})
			}

			if providedListorHash.Sha256 != fmt.Sprintf("%x", *byHash) {
				return baseline.NewPostBaselineGetPropBadRequest().WithPayload(
					generalError{Code: 400, Msg: "Listor hash mismatch"})
			}
		}
	}

	listProp := bIns.GetProp(
		framework.SetAuthProviderOpt(&serverAuthProvider{params.Profile}),
		framework.SetDataProviderOpt(&dataProvider),
	)
	if len(listProp) != len(b.Checker) {
		return middleware.Error(500, generalError{Code: 500, Msg: "size mismatch between CheckerProp and Checker"})
	}

	listPropRes := make([]string, len(listProp))
	for i, prop := range listProp {
		byProp, err := json.Marshal(prop)
		if err != nil {
			return middleware.Error(500, generalError{
				Code: 500,
				Msg:  fmt.Sprintf("Failed to marshal to json: %v", err),
			})
		}

		listPropRes[i] = string(byProp)
	}

	byBaselineHash, err := getBaselineHash(int(params.ID - 1))
	if err != nil {
		return middleware.Error(500, generalError{
			Code: 500,
			Msg:  fmt.Sprintf("Failed to get Baseline hash: %v", err),
		})
	}

	data4api := server_model.BaselineData{
		ID:           params.ID,
		BaselineHash: &server_model.ItemHash{Sha256: fmt.Sprintf("%x", *byBaselineHash)},
		CheckerProp:  listPropRes,
	}

	return baseline.NewPostBaselineGetPropOK().WithPayload(
		&baseline.PostBaselineGetPropOKBody{
			Code: 200,
			Msg:  "success",
			Data: &data4api,
		})
}

func baselinePostBaselineValidateHandler(params baseline.PostBaselineValidateParams) middleware.Responder {
	if !_confValid {
		return middleware.Error(500, generalError{Code: 500, Msg: "config.conf file not loaded"})
	}

	if params.ID <= 0 || params.ID > int64(len(_conf.Baseline)) {
		return baseline.NewPostBaselineValidateNotFound().WithPayload(
			generalError{Code: 404, Msg: "Id not found"})
	} else if params.ID != params.Data.ID {
		return baseline.NewPostBaselineValidateBadRequest().WithPayload(
			generalError{Code: 400, Msg: fmt.Sprintf("Baseline id mismatch between %d and %d",
				params.ID, params.Data.ID)})
	}

	b := _conf.Baseline[params.ID-1]
	bIns, err := getBaseline(int(params.ID - 1))
	if err != nil {
		return middleware.Error(500, generalError{
			Code: 500,
			Msg:  fmt.Sprintf("Failed to get Baseline instance: %v", err),
		})
	}

	byHash, err := getBaselineHash(int(params.ID - 1))
	if err != nil {
		return middleware.Error(500, generalError{
			Code: 500,
			Msg:  fmt.Sprintf("Failed to get Baseline hash: %v", err),
		})
	}
	if params.Data.BaselineHash.Sha256 != fmt.Sprintf("%x", *byHash) {
		return baseline.NewPostBaselineValidateBadRequest().WithPayload(
			generalError{Code: 400, Msg: "Baseline hash mismatch"})
	}
	if len(params.Data.CheckerProp) != len(b.Checker) {
		return baseline.NewPostBaselineValidateBadRequest().WithPayload(
			generalError{Code: 400, Msg: fmt.Sprintf("size mismatch between given CheckerProp %d and Checker %d",
				len(params.Data.CheckerProp), len(b.Checker))})
	}

	data := make(framework.BaselinePropList, len(params.Data.CheckerProp))
	for i, prop := range params.Data.CheckerProp {
		var provideData framework.CheckerPropList
		if err := internal.JsonUnmarshal([]byte(prop), &provideData); err != nil {
			return baseline.NewPostBaselineValidateBadRequest().WithPayload(
				generalError{Code: 400, Msg: fmt.Sprintf("failed to unmarshal from json: %s", err)})
		}

		data[i] = provideData
	}

	resBaseline, err := bIns.Validate(data)
	if err != nil {
		return baseline.NewPostBaselineValidateBadRequest().WithPayload(
			generalError{Code: 400, Msg: fmt.Sprintf("Failed in Validate: %v", err)})
	}

	data4api := make([]*server_model.ValidateResult, 0)
	for _, res := range resBaseline {
		if res.InRisk || params.RiskOnly == nil || !*params.RiskOnly {
			singleOutputData := server_model.ValidateResult{
				CloudType:      string(res.CloudType),
				ResourceID:     res.Id,
				ResourceName:   res.Name,
				ActualValue:    res.Value,
				ResourceInRisk: res.InRisk,
				Metadata:       make(map[string]string),
			}

			for _, key := range params.Metadata {
				value := (*bIns.GetMetadata())[key]
				singleOutputData.Metadata[key] = value
			}

			data4api = append(data4api, &singleOutputData)
		}
	}

	return baseline.NewPostBaselineValidateOK().WithPayload(
		&baseline.PostBaselineValidateOKBody{
			Code: 200,
			Msg:  "success",
			Data: data4api,
		})
}

func setupHandler(api *operations.CloudBenchCheckerAPIAPI) {
	prepareHandler(api)

	api.BaselineGetBaselineGetDefinitionHandler = baseline.GetBaselineGetDefinitionHandlerFunc(
		baselineGetBaselineGetDefinitionHandler)
	api.BaselineGetBaselineGetIdsHandler = baseline.GetBaselineGetIdsHandlerFunc(
		baselineGetBaselineGetIdsHandler)
	api.BaselineGetBaselineGetListorIDHandler = baseline.GetBaselineGetListorIDHandlerFunc(
		baselineGetBaselineGetListorIDHandler)
	api.ListorGetListorGetDefinitionHandler = listor.GetListorGetDefinitionHandlerFunc(
		listorGetListorGetDefinitionHandler)
	api.ListorGetListorGetIdsHandler = listor.GetListorGetIdsHandlerFunc(
		listorGetListorGetIdsHandler)
	api.ListorGetListorListDataHandler = listor.GetListorListDataHandlerFunc(
		listorGetListorListDataHandler)
	api.BaselinePostBaselineGetPropHandler = baseline.PostBaselineGetPropHandlerFunc(
		baselinePostBaselineGetPropHandler)
	api.BaselinePostBaselineValidateHandler = baseline.PostBaselineValidateHandlerFunc(
		baselinePostBaselineValidateHandler)
}
