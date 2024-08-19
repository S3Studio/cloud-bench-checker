// Json util

package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// JsonUnmarshal: Unmarshal json string and treat number as json.Number
func JsonUnmarshal(data []byte, v any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(v)
}

// JsonMarshal: Convert interface{} to *json.RawMessage
func JsonMarshal(v any) (*json.RawMessage, error) {
	byJson, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value to json: %w", err)
	}

	var rm json.RawMessage = byJson
	return &rm, nil
}
