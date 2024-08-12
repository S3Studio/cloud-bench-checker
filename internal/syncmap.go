// SyncMap util
package internal

import "sync"

// SyncMap[T]: Generics type of wrapper of sync.Map
type SyncMap[T any] struct {
	sync.Map
}

// LoadOrCreate: Load an object of type T, and create a new one if it does not exist in the sync.map.
// @param: key: Key of object in the map
// @param: fnCreate: Function to create an object, normally a closure
// @param: nilVal: Nil value if fnCreate failed. Generics type does not support return nil as T
// @return: Actual object of type T loaded or created
// @return: Error
func (m *SyncMap[T]) LoadOrCreate(key string, fnCreate func() (any, error), nilVal T) (T, error) {
	val, ok := m.Load(key)
	if !ok {
		newVal, err := fnCreate()
		if err != nil {
			return nilVal, err
		}
		// May have already been created by other goroutions,
		// but it's ok to spend a little more time creating them
		val, _ = m.LoadOrStore(key, newVal)
	}

	tVal, ok := val.(T)
	if !ok {
		panic("internal error, not a valid type in sync.Map")
	}
	return tVal, nil
}
