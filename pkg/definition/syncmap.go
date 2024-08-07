// Util of SyncMap
package definition

import "sync"

type SyncMap[T any] struct {
	sync.Map
}

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
