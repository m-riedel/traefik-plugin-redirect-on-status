package util

// ContainsStr returns true if the given value is in the given collection
func ContainsStr(collection []string, value string) bool {
	for _, field := range collection {
		if field == value {
			return true
		}
	}
	return false
}

// ContainsFunc returns true if for any given value of collection, the function f returns true
func ContainsFunc[C ~[]E, E any](collection C, f func(v1 E) bool) bool {
	for _, field := range collection {
		if f(field) {
			return true
		}
	}
	return false
}
