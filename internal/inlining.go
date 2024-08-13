// Util to disable inlining

package internal

func DisableInlining() int {
	for i := 0; i < 10; i++ {
	}
	return 0
}
