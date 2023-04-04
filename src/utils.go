package main

func Contains[T any](slice []T, block func(T) bool) bool {
	for _, element := range slice {
		if block(element) {
			return true
		}
	}
	return false
}
