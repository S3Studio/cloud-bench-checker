// Auth controller

package auth

import (
	"os"
	"path/filepath"
	"testing"

	def "github.com/s3studio/cloud-bench-checker/pkg/definition"
	"github.com/s3studio/cloud-bench-checker/test"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/spf13/viper"
)

func TestAuthFileProvider_GetProfile(t *testing.T) {
	patches := gomonkey.ApplyMethodFunc(&viper.Viper{}, "ReadInConfig",
		func() error {
			return nil
		})
	defer patches.Reset()

	type args struct {
		cloudType def.CloudType
	}
	tests := []struct {
		name string
		p    *AuthFileProvider
		args args
		//want    *viper.Viper
		wantErr bool
	}{
		{
			"Valid result with conf of env",
			NewAuthFileProvider(test.Test_conf_env),
			args{def.TENCENT_CLOUD},
			false,
		},
		{
			"Valid result with conf of file",
			NewAuthFileProvider(test.Test_conf_file),
			args{def.TENCENT_CLOUD},
			false,
		},
		{
			"Profile not defined",
			NewAuthFileProvider(test.Test_conf_file),
			args{def.ALIYUN_CLOUD},
			true,
		},
		{
			"Arg of TENCENT_CLOUD",
			NewAuthFileProvider(test.Test_conf_env),
			args{def.TENCENT_CLOUD},
			false,
		},
		{
			"Arg of TENCENT_COS",
			NewAuthFileProvider(test.Test_conf_env),
			args{def.TENCENT_COS},
			false,
		},
		{
			"Arg of ALIYUN_CLOUD",
			NewAuthFileProvider(test.Test_conf_env),
			args{def.ALIYUN_CLOUD},
			true,
		},
		{
			"Arg of ALIYUN_OSS",
			NewAuthFileProvider(test.Test_conf_env),
			args{def.ALIYUN_OSS},
			true,
		},
		{
			"Arg of K8S",
			NewAuthFileProvider(test.Test_conf_env),
			args{def.K8S},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetProfile(tt.args.cloudType)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthFileProvider.GetProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("AuthFileProvider.GetProfile() = %v, want a valid pointer", got)
			}
		})
	}
}

func Test_readProfile(t *testing.T) {
	patches := gomonkey.ApplyMethodFunc(&viper.Viper{}, "ReadInConfig",
		func() error {
			return nil
		})
	defer patches.Reset()

	type args struct {
		profileName string
	}
	tests := []struct {
		name string
		args args
		//want    *viper.Viper
		wantErr bool
	}{
		{
			"Valid result of file",
			args{"mock_filename"},
			false,
		},
		{
			"Valid result of env",
			args{"$ENV"},
			false,
		},
		{
			"invalid profile name, should contains filename only without dir",
			args{"/root/file"},
			true,
		},
		{
			"invalid profile name, should contains filename only without dir",
			args{"../file"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readProfile(tt.args.profileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("readProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("readProfile() = %v, want a valid pointer", got)
			}
		})
	}
}

func TestAuthFileProvider_GetProfilePathname(t *testing.T) {
	type args struct {
		cloudType def.CloudType
	}
	tests := []struct {
		name    string
		p       *AuthFileProvider
		args    args
		want    string
		wantErr bool
	}{
		{
			"Valid result with conf of file",
			NewAuthFileProvider(test.Test_conf_file),
			args{def.TENCENT_CLOUD},
			"file",
			false,
		},
		{
			"Valid result with conf of k8s",
			NewAuthFileProvider(test.Test_conf_k8s),
			args{def.K8S},
			"config",
			false,
		},
		{
			"no profile defined for cloud",
			NewAuthFileProvider(test.Test_conf_env),
			args{def.ALIYUN_CLOUD},
			"",
			true,
		},
		{
			"no pathname defined for cloud with profile of $ENV",
			NewAuthFileProvider(test.Test_conf_env),
			args{def.TENCENT_CLOUD},
			"",
			true,
		},
		{
			"invalid profile name, should contains filename only without dir",
			NewAuthFileProvider(def.ConfProfile{"tencent": "/root/file"}),
			args{def.TENCENT_CLOUD},
			"",
			true,
		},
		{
			"invalid profile name, should contains filename only without dir",
			NewAuthFileProvider(def.ConfProfile{"tencent": "../file"}),
			args{def.TENCENT_CLOUD},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetProfilePathname(tt.args.cloudType)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthFileProvider.GetProfilePathname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_, gotFilename := filepath.Split(got)
			if gotFilename != tt.want {
				t.Errorf("AuthFileProvider.GetProfilePathname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAllSet(t *testing.T) {
	envMap := map[string]string{
		"TENCENTCLOUD_SECRET_ID":  "mock_secretid",
		"TENCENTCLOUD_SECRET_KEY": "mock_secretkey",
		"TENCENTCLOUD_REGION":     "mock_region",
		"TEST_KEY":                "mock_value",
	}
	for k, v := range envMap {
		os.Setenv(k, v)
	}
	testViper := viper.New()
	testViper.AutomaticEnv()

	type args struct {
		v    *viper.Viper
		keys []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Key is set",
			args{testViper, []string{"TENCENTCLOUD_SECRET_ID", "TENCENTCLOUD_SECRET_KEY", "TENCENTCLOUD_REGION", "TEST_KEY"}},
			false,
		},
		{
			"Key is not set",
			args{testViper, []string{"invalidkey"}},
			true,
		},
		{
			"Invalid viper",
			args{nil, nil},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := IsAllSet(tt.args.v, tt.args.keys); (err != nil) != tt.wantErr {
				t.Errorf("IsAllSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
