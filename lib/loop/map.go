package loop

func Map[I any, R any](array []I, callback func(value I) R) []R {
	result := []R{}

	for _, value := range array {
		result = append(result, callback(value))
	}

	return result
}
