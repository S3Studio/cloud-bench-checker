// Param util
package internal

import (
	"fmt"
	"strconv"

	"github.com/s3studio/cloud-bench-checker/pkg/definition"
)

// AddParamString: Insert a string into map[string]any as targetType
func AddParamString(m map[string]any, key string, value string, targetType definition.ParamType) error {
	switch targetType {
	case definition.PARAM_INT:
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("failed to convert string to int: %w", err)
		}
		m[key] = i
	case definition.PARAM_STRING:
		m[key] = value
	case definition.PARAM_STRING_LIST:
		m[key] = []string{value}
	default:
		return fmt.Errorf("invalid param type: %s", targetType)
	}

	return nil
}

// AddParamInt: Insert a int into map[string]any as targetType
func AddParamInt(m map[string]any, key string, value int, targetType definition.ParamType) error {
	switch targetType {
	case definition.PARAM_INT:
		m[key] = value
	case definition.PARAM_STRING:
		m[key] = strconv.Itoa(value)
	case definition.PARAM_STRING_LIST:
		m[key] = []string{strconv.Itoa(value)}
	default:
		return fmt.Errorf("invalid param type: %s", targetType)
	}

	return nil
}
