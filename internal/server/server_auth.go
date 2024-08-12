// Auth controller for apiserver
package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/s3studio/cloud-bench-checker/internal"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	"github.com/spf13/viper"
)

// serverAuthProvider: Implementation of IAuthProvider using files in the ".auth" subdirectory
type serverAuthProvider struct {
	// Name of profile passed from API
	profile string
}

// ProfileNotDefinedError: Error of profile not defined
type ProfileNotDefinedError struct {
	// Value of profile name
	key string
}

// Error: Output error string
// @return: Error string
func (e ProfileNotDefinedError) Error() string {
	return fmt.Sprintf("no profile defined for cloud: %s", e.key)
}

var mapViper internal.SyncMap[*viper.Viper]

// GetProfile: Implementation of IAuthProvider.GetProfile
// @param: cloudType: Type of the cloud, omitted in this implementation of IAuthProvider
// @return: Profile that can be accessed as Viper
// @return: Error
func (p *serverAuthProvider) GetProfile(_ def.CloudType) (*viper.Viper, error) {
	return mapViper.LoadOrCreate(p.profile, func() (any, error) {
		return readProfile(p.profile)
	}, nil)
}

func readProfile(profileName string) (*viper.Viper, error) {
	v := viper.New()

	if profileName == def.PROFILE_ENV {
		v.AutomaticEnv()
		return v, nil
	}

	if dir, _ := filepath.Split(profileName); dir != "" {
		// File not in subdirectory is not allowed
		return nil, errors.New("invalid profile name, should only contain filename without directory")
	}

	binPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get binary path: %w", err)
	}
	binDir, _ := filepath.Split(binPath)

	v.SetConfigFile(filepath.Join(binDir, ".auth", profileName))
	v.SetConfigType("properties")
	if err := v.ReadInConfig(); err != nil {
		// Do not use value of err to avoid leaking the file path
		return nil, errors.New("unable to read config file")
	}

	v.WatchConfig()

	return v, nil
}

var _defaultProfilePathname = map[string]string{
	"k8s": "~/.kube/config",
}

// GetProfilePathname: Implement of IAuthProvider.GetProfilePathname
// @param: cloudType: Type of the cloud
// @return: Pathname of profile
// @return: Error
func (p *serverAuthProvider) GetProfilePathname(cloudType def.CloudType) (string, error) {
	if p.profile == def.PROFILE_ENV {
		pathname, ok := _defaultProfilePathname[string(cloudType)]
		if !ok {
			return "", fmt.Errorf("no default pathname defined for cloud \"%s\" with profile of $ENV", cloudType)
		}

		if pathname[:2] == "~/" {
			// Parse home directory of user
			if homeDir, _ := os.UserHomeDir(); homeDir != "" {
				pathname = filepath.Join(homeDir, pathname[2:])
			}
		}

		return pathname, nil
	}

	if dir, _ := filepath.Split(p.profile); dir != "" {
		// File not in subdirectory is not allowed
		return "", errors.New("invalid profile name, should only contain filename without directory")
	}

	binPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", fmt.Errorf("failed to get binary path: %w", err)
	}
	binDir, _ := filepath.Split(binPath)

	return filepath.Join(binDir, ".auth", p.profile), nil
}
