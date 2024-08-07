// Baseline of management for process of checking
package framework

import (
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/s3studio/cloud-bench-checker/pkg/auth"
	def "github.com/s3studio/cloud-bench-checker/pkg/definition"
)

// Baseline: Used to manage checkers and listors.
//
// Usage of Baseline consists of 3 steps:
// (Optional) 1. GetListorId: Get the ids of the listors used in all the Checkers of the Baseline,
// which may be used to prepare the raw data in advance
// 2. GetProp: Extract properties from the raw data provided by the IDataProvider,
// which can retrieve it from the cloud connector or cache.
// Additional data would be retrieved directly via the cloud connector on demand.
// 3. Validate: Validate the property against the benchmark and return the result
type Baseline struct {
	// Definition of Baseline
	conf *def.ConfBaseline
	// Checkers correspond to the Baseline
	checker []*Checker
}

// BaselinePropList: Type alias of list of CheckerPropList for a Baseline
//
// The order of the items in the list must be the same as the order of the Checkers in the Baseline
type BaselinePropList []CheckerPropList

// NewBaseline: Constructor of Baseline
// @param: conf: Definition of Baseline
// @param: authProvider: IAuthProvider to provide profile of auth
// @param: dataProvider: IDataProvider to provide raw data
func NewBaseline(conf *def.ConfBaseline, authProvider auth.IAuthProvider, dataProvider IDataProvider) *Baseline {
	baseline := Baseline{conf, make([]*Checker, len(conf.Checker))}
	for i, confChecker := range conf.Checker {
		baseline.checker[i] = NewChecker(&confChecker, authProvider, dataProvider)
	}

	return &baseline
}

// SetAuthProvider: Set new authProvider for all checkers
// @param: authProvider: New provider
func (b *Baseline) SetAuthProvider(authProvider auth.IAuthProvider) {
	for _, c := range b.checker {
		c.SetAuthProvider(authProvider)
	}
}

// SetDataProvider: Set new dataProvider for all checkers
// @param: dataProvider: New provider
func (b *Baseline) SetDataProvider(dataProvider IDataProvider) {
	for _, c := range b.checker {
		c.SetDataProvider(dataProvider)
	}
}

// GetListorId: Get the ids of the listors used in all the Checkers of the Baseline
// @return: ids of listors
func (b *Baseline) GetListorId() []int {
	listorIds := make([]int, 0, len(b.conf.Checker))

	for _, checker := range b.conf.Checker {
		for _, idAdd := range checker.Listor {
			bAdd := true
			for _, idExist := range listorIds {
				if idAdd == idExist {
					bAdd = false
					break
				}
			}

			if bAdd {
				listorIds = append(listorIds, idAdd)
			}
		}
	}

	return listorIds
}

// GetProp: Extract properties from the raw data
//
// The length of the outer list is equal to the length of checkers
// @return: List of the result of GetProp of each checker, whose' elements are the list of props extracted from raw data
func (b *Baseline) GetProp(opts ...GetPropOption) BaselinePropList {
	var checkerPropList = make(BaselinePropList, len(b.checker))

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(b.checker))

	for i, checker := range b.checker {
		go func(target *CheckerPropList) {
			singleCheckerProp, err := checker.GetProp(opts...)
			if err != nil {
				// Print error and skip the current checker
				glog().Println(err)
			} else {
				*target = append(*target, singleCheckerProp...)
			}

			waitGroup.Done()
		}(&checkerPropList[i])
	}

	waitGroup.Wait()

	return checkerPropList
}

// Validate: Validate the property against the benchmark and return the result
//
// NOTE: The length of the list of data must be the same as the length of checkers,
// as each item in the list is sent to a checker in order
// @param: data: List of properties to be validated
// @return: List of validation results
// @return: Error
func (b *Baseline) Validate(data BaselinePropList) ([]*ValidateResult, error) {
	if len(data) != len(b.checker) {
		return nil, errors.New("size mismatch between props and checkers, please review the given data")
	}

	var validateResult []*ValidateResult
	for i, checker := range b.checker {
		singleResult, err := checker.Validate(data[i])
		if err != nil {
			// Print error and skip the current checker
			glog().Println(err)
			continue
		}

		validateResult = append(validateResult, singleResult...)
	}

	return validateResult, nil
}

// GetMetadata: Get the metadata defined in Baseline.conf
// @return: metadata
func (b *Baseline) GetMetadata() *map[string]string {
	return &b.conf.Metadata
}

func (b *Baseline) GetHash(hashType crypto.Hash, listorHashList [][]*[]byte) ([]byte, error) {
	if len(listorHashList) != len(b.conf.Checker) {
		return nil, errors.New("size mismatch between Checker and given hash list")
	}

	for i, c := range b.conf.Checker {
		if len(listorHashList[i]) != len(c.Listor) {
			return nil, fmt.Errorf("size mismatch between Checker and given hash list of #%d", i)
		}
	}

	// Copy conf to a var of any
	byBaseline, err := json.Marshal(*b.conf)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal conf to json: %w", err)
	}

	var objForHash any
	if err := json.Unmarshal(byBaseline, &objForHash); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conf from json: %w", err)
	}

	// Replace Listor of each Checker to hash of Listor
	var objChecker []any
	objChecker, ok := objForHash.(map[string]any)["Checker"].([]any)
	if !ok {
		return nil, errors.New("failed to unmarshal Checker from json")
	}

	for i, item := range objChecker {
		objItem, ok := item.(map[string]any)
		if !ok {
			return nil, errors.New("failed to unmarshal Checker item from json")
		}

		delete(objItem, "Listor")
		objItem["ListorHash"] = listorHashList[i]
	}

	// Calculate hash
	return CalcHash(hashType, objForHash)
}
