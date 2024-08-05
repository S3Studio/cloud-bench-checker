// Package auth:
// Auth controller
package auth

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	def "github.com/s3studio/cloud-bench-checker/pkg/definition"

	"github.com/spf13/viper"
)

// IAuthProvider: Interface that provides different management of the profile of Auth
type IAuthProvider interface {
	// GetProfile: Get profile for the cloud to connect to
	// @param: cloudType: Type of the cloud
	// @return: Profile that can be accessed as Viper
	// @return: Error
	GetProfile(cloudType def.CloudType) (*viper.Viper, error)
	// GetProfilePathname:
	// Get pathname of profile for the cloud connector that needs to read it directly
	// @param: cloudType: Type of the cloud
	// @return: Pathname of profile
	// @return: Error
	GetProfilePathname(cloudType def.CloudType) (string, error)
}

// AuthFileProvider: Implementation of IAuthProvider using files in the ".auth" subdirectory
type AuthFileProvider struct {
	// Definition of profile
	profile def.ConfProfile
	// sync.Map which stores the cache of vipers of profile
	mapViper sync.Map
}

// NewAuthFileProvider: Constructor of AuthFileProvider
// @param: profile: Definition of profile
func NewAuthFileProvider(profile def.ConfProfile) *AuthFileProvider {
	return &AuthFileProvider{profile: profile}
}

// _mapCloudTypeToName: Stores mapping from CloudType to profile key.
// Multiple cloud types can use the same profile key for simplicity,
// e.g.: TencentCloud and TencentCos can both use "tencent"
var _mapCloudTypeToName = map[def.CloudType]string{
	def.TENCENT_CLOUD: "tencent",
	def.TENCENT_COS:   "tencent",
	def.ALIYUN_CLOUD:  "aliyun",
	def.ALIYUN_OSS:    "aliyun",
	def.K8S:           "k8s",
	def.AZURE:         "azure",
}

// ProfileNotDefinedError: Error of profile not defined
// It may be acceptable to use one conf file in different projects
// with different cloud environments
type ProfileNotDefinedError struct {
	// Value of profile name
	key string
}

// Error: Output error string
// @return: Error string
func (e ProfileNotDefinedError) Error() string {
	return fmt.Sprintf("no profile defined for cloud: %s", e.key)
}

// GetProfile: Implementation of IAuthProvider.GetProfile
// @param: cloudType: Type of the cloud
// @return: Profile that can be accessed as Viper
// @return: Error
func (p *AuthFileProvider) GetProfile(cloudType def.CloudType) (*viper.Viper, error) {
	key, ok := _mapCloudTypeToName[cloudType]
	if !ok {
		panic(fmt.Sprintf("internal error, key name of cloudType \"%s\" not assigned", cloudType))
	}

	profileName, ok := p.profile[key]
	if !ok {
		return nil, ProfileNotDefinedError{key}
	}

	v, ok := p.mapViper.Load(profileName)
	if !ok {
		newViper, err := readProfile(profileName)
		if err != nil {
			return nil, err
		}
		// May have already been created by other goroutions,
		// but it's ok to spend a little more time creating them
		v, _ = p.mapViper.LoadOrStore(profileName, newViper)
	}

	viperValue, ok := v.(*viper.Viper)
	if !ok {
		panic("internal error, not a valid viper in sync.Map")
	}
	return viperValue, nil
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
		return nil, err
	}

	return v, nil
}

var _defaultProfilePathname = map[string]string{
	"k8s": "~/.kube/config",
}

// GetProfilePathname: Implement of IAuthProvider.GetProfilePathname
// @param: cloudType: Type of the cloud
// @return: Pathname of profile
// @return: Error
func (p *AuthFileProvider) GetProfilePathname(cloudType def.CloudType) (string, error) {
	key, ok := _mapCloudTypeToName[cloudType]
	if !ok {
		panic(fmt.Sprintf("internal error, key name of cloudType \"%s\" not assigned", cloudType))
	}

	profileName, ok := p.profile[key]
	if !ok {
		return "", ProfileNotDefinedError{key}
	}

	if profileName == def.PROFILE_ENV {
		pathname, ok := _defaultProfilePathname[key]
		if !ok {
			return "", fmt.Errorf("no default pathname defined for cloud \"%s\" with profile of $ENV", key)
		}

		if pathname[:2] == "~/" {
			// Parse home directory of user
			if homeDir, _ := os.UserHomeDir(); homeDir != "" {
				pathname = filepath.Join(homeDir, pathname[2:])
			}
		}

		return pathname, nil
	}

	if dir, _ := filepath.Split(profileName); dir != "" {
		// File not in subdirectory is not allowed
		return "", errors.New("invalid profile name, should only contain filename without directory")
	}

	binPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", fmt.Errorf("failed to get binary path: %w", err)
	}
	binDir, _ := filepath.Split(binPath)

	return filepath.Join(binDir, ".auth", profileName), nil
}

// IsAllSet: Check for all required keys,
// otherwise Viper.GetString will return an empty value other than error
// @param: v: Instance of viper
// @param: keys: List of keys to be checked
// @return: Error
func IsAllSet(v *viper.Viper, keys []string) error {
	if v == nil {
		return errors.New("invalid viper instance of nil")
	}

	for _, k := range keys {
		if !v.IsSet(k) {
			return fmt.Errorf("failed to read key from profile: %s", k)
		}
	}

	return nil
}
