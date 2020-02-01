package utils

// ReuseString is delete value in slice
func ReuseString(origin []string, v string) []string {
	target := origin[:0]
	for _, item := range origin {
		if item != v {
			target = append(target, item)
		}
	}
	return target
}

// ReuseInt64 is delete value in slice
func ReuseInt64(origin []int64, v int64) []int64 {
	target := origin[:0]
	for _, item := range origin {
		if item != v {
			target = append(target, item)
		}
	}
	return target
}

// ChickInString is check value in slice
func ChickInString(origin []string, v string) bool {
	for _, item := range origin {
		if item == v {
			return true
		}
	}
	return false
}

// ChickInInt64 is check value in slice
func ChickInInt64(origin []int64, v int64) bool {
	for _, item := range origin {
		if item == v {
			return true
		}
	}
	return false
}
