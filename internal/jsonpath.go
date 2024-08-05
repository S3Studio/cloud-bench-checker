// Jsonpath util
package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bhmj/jsonslice"
)

func ParseJsonPath(input *json.RawMessage, path string) (*json.RawMessage, error) {
	if input == nil {
		return nil, errors.New("nil input when parsing JsonPath")
	}
	byInput, err := input.MarshalJSON()
	if err != nil {
		return nil, err
	}

	byReturn, err := jsonslice.Get(byInput, path)
	if err != nil {
		return nil, err
	}

	var rmReturn json.RawMessage
	rmReturn.UnmarshalJSON(byReturn)

	return &rmReturn, nil
}

// ParseJsonPathStr: Parse JsonPath and try to convert result to single string
func ParseJsonPathStr(input *json.RawMessage, path string) (string, error) {
	parseRes, err := ParseJsonPath(input, path)
	if err != nil {
		return "", err
	}
	if len(*parseRes) == 0 {
		// Treat nil as an empty string, otherwise bytes.NewReader will throw an error in JsonUnmarshal
		return "", nil
	}

	var jsonRes any
	if err := JsonUnmarshal(*parseRes, &jsonRes); err != nil {
		return "", err
	}

	return parseToStr(jsonRes)
}

func parseToStr(data any) (string, error) {
	switch d := data.(type) {
	case string:
		return d, nil
	case json.Number:
		return d.String(), nil
	case bool:
		return strconv.FormatBool(d), nil
	case []any:
		strList := make([]string, len(d))
		for i, item := range d {
			var err error
			if strList[i], err = parseToStr(item); err != nil {
				return "", err
			}
		}
		return fmt.Sprintf("[%s]", strings.Join(strList, ",")), nil
	case map[string]any:
		strList := make([]string, 0, len(d))
		for k, v := range d {
			if str_v, err := parseToStr(v); err != nil {
				return "", err
			} else {
				strList = append(strList, fmt.Sprintf("%s:%s", k, str_v))
			}
		}
		return fmt.Sprintf("{%s}", strings.Join(strList, ",")), nil
	case nil:
		// Treat nil as an empty string
		// Should be checked afterwards if it is not intended
		return "", nil
	default:
		return "", fmt.Errorf("failed to parse to string from data: %v", data)
	}
}

// ParseJsonPathList: Parse JsonPath and try to return the result as a list
func ParseJsonPathList(input *json.RawMessage, path string, bObjAsList bool) ([]*json.RawMessage, error) {
	parseRes, err := ParseJsonPath(input, path)
	if err != nil {
		return nil, err
	}

	if bObjAsList {
		return []*json.RawMessage{parseRes}, nil
	}

	var listRes []*json.RawMessage
	if err := JsonUnmarshal(*parseRes, &listRes); err != nil {
		return nil, err
	}

	return listRes, nil
}
