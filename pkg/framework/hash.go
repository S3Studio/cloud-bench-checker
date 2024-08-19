// Util to calculate hash

package framework

import (
	"crypto"
	"encoding/json"
	"fmt"
)

// CalcHash: Calculate specific hash of any object
// @param: hashType: Method of hash
// @param: obj: Object to calculate
// @return: Hash value as []byte. Convert to string with `fmt.Sprintf("%x", hash)` is recommended
// @return: Error
func CalcHash(hashType crypto.Hash, obj any) ([]byte, error) {
	byCal, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object to json: %w", err)
	}

	hashInstance := hashType.New()
	_, err = hashInstance.Write(byCal)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}

	return hashInstance.Sum(nil), nil
}
