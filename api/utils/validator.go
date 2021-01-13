package utils

import (
	validator "github.com/asaskevich/govalidator"
)

// IsInteger to validate if all value is integer from string type values
func IsInteger(values ...string) bool {
	for _, value := range values {
		if !validator.IsInt(value) {
			return false
		}
	}

	return true
}

// IsPositiveInteger to validate if all value is positive integer from integer type values
func IsPositiveInteger(values ...int) bool {
	for _, value := range values {
		if !validator.IsPositive(float64(value)) {
			return false
		}
	}

	return true
}
