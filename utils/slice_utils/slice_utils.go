package slice_utils

func FindIndexInSlice(n int, comparer func(index int) bool) int {
	for index := 0; index < n; index++ {
		if comparer(index) {
			return index
		}
	}
	return -1
}
