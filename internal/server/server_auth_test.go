// Auth controller for apiserver

package server

import (
	"path/filepath"
	"testing"

	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/spf13/viper"
)

func TestServerAuthProvider_GetProfile(t *testing.T) {
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
		p    *serverAuthProvider
		args args
		//want    *viper.Viper
		wantErr bool
	}{
		{
			"Valid result",
			&serverAuthProvider{"mock"},
			args{def.TENCENT_CLOUD},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetProfile(tt.args.cloudType)
			if (err != nil) != tt.wantErr {
				t.Errorf("serverAuthProvider.GetProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("serverAuthProvider.GetProfile() = %v, want a valid pointer", got)
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

func TestServerAuthProvider_GetProfilePathname(t *testing.T) {
	type args struct {
		cloudType def.CloudType
	}
	tests := []struct {
		name    string
		p       *serverAuthProvider
		args    args
		want    string
		wantErr bool
	}{
		{
			"Valid result with conf of file",
			&serverAuthProvider{"file"},
			args{def.TENCENT_CLOUD},
			"file",
			false,
		},
		{
			"Valid result with conf of $ENV",
			&serverAuthProvider{"$ENV"},
			args{def.K8S},
			"config",
			false,
		},
		{
			"no pathname defined for cloud with profile of $ENV",
			&serverAuthProvider{"$ENV"},
			args{def.TENCENT_CLOUD},
			"",
			true,
		},
		{
			"invalid profile name, should contains filename only without dir",
			&serverAuthProvider{"/root/file"},
			args{def.TENCENT_CLOUD},
			"",
			true,
		},
		{
			"invalid profile name, should contains filename only without dir",
			&serverAuthProvider{"../file"},
			args{def.TENCENT_CLOUD},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.GetProfilePathname(tt.args.cloudType)
			if (err != nil) != tt.wantErr {
				t.Errorf("serverAuthProvider.GetProfilePathname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_, gotFilename := filepath.Split(got)
			if gotFilename != tt.want {
				t.Errorf("serverAuthProvider.GetProfilePathname() = %v, want %v", got, tt.want)
			}
		})
	}
}
