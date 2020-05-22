package utils

import "errors"

func Min(values ...int) int {
	if len(values) == 0 {
		panic(errors.New("min: no arguments"))
	}

	minValue := values[0]
	for _, value := range values[1:] {
		if minValue > value {
			minValue = value
		}
	}

	return minValue
}

func Max(values ...int) int {
	if len(values) == 0 {
		panic(errors.New("max: no arguments"))
	}

	maxValue := values[0]
	for _, value := range values[1:] {
		if value > maxValue {
			maxValue = value
		}
	}

	return maxValue
}
