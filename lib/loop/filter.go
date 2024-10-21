package loop

func Filter[I any](array []I, callback func(value I) bool) []I {
	result := []I{}

	for _, value := range array {
		if callback(value) {
			result = append(result, value)
		}
	}

	return result
}
