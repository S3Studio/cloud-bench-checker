package test

import (
	"errors"

	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	"github.com/spf13/viper"
)

var (
	Test_conf_env     = def.ConfProfile{"tencent": "$ENV"}
	Test_conf_file    = def.ConfProfile{"tencent": "file"}
	Test_conf_aliyun  = def.ConfProfile{"aliyun": "$ENV"}
	Test_conf_k8s     = def.ConfProfile{"k8s": "$ENV"}
	Test_conf_azure   = def.ConfProfile{"azure": "$ENV"}
	Test_conf_invalid = def.ConfProfile{}
)

type MockKeyNotSetAuthProvider struct{}

func (p *MockKeyNotSetAuthProvider) GetProfile(cloudType def.CloudType) (*viper.Viper, error) {
	return viper.New(), nil
}

func (p *MockKeyNotSetAuthProvider) GetProfilePathname(cloudType def.CloudType) (string, error) {
	return "", errors.New("mock function not implemented")
}
